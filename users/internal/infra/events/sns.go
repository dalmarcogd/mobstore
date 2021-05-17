package events

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"

	"github.com/dalmarcogd/mobstore/users/internal/domains"
	"github.com/dalmarcogd/mobstore/users/internal/infra/errors"
	"github.com/dalmarcogd/mobstore/users/internal/infra/httptransport"
	"github.com/dalmarcogd/mobstore/users/internal/infra/iaws"
	"github.com/dalmarcogd/mobstore/users/internal/infra/ptrs"
	"github.com/dalmarcogd/mobstore/users/internal/services"
)

const (
	productTopic Topic = iota
)

type (
	Topic int

	eventsService struct {
		serviceManager services.ServiceManager
		ctx            context.Context
		client         iaws.SNSClient
		topic          Topic
		address        string
	}

	snsEndpointResolve struct {
		endpointResolverFunc aws.EndpointResolverFunc
	}
)

func (s snsEndpointResolve) ResolveEndpoint(region string, _ sns.EndpointResolverOptions) (aws.Endpoint, error) {
	return s.endpointResolverFunc("", region)
}

func New() *eventsService {
	return &eventsService{}
}

func (s *eventsService) WithProductTopic() services.Events {
	s.topic = productTopic
	return s
}

func (s *eventsService) Init(ctx context.Context) error {
	s.ctx = ctx

	var topicAddress string
	switch s.topic {
	case productTopic:
		topicAddress = s.ServiceManager().Environment().UserTopicAddress()
	default:
		return errors.EventAddressIsRequiredError
	}

	s.address = topicAddress

	if s.client == nil {
		client := &http.Client{}
		transport, err := httptransport.NewTransport(s.ServiceManager().Spans().Tracer())
		if err != nil {
			return err
		}
		client.Transport = transport

		customResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           s.ServiceManager().Environment().AwsEndpoint(),
				SigningRegion: s.ServiceManager().Environment().AwsRegion(),
			}, nil
		})

		staticCredentialsProvider := credentials.NewStaticCredentialsProvider(
			s.ServiceManager().Environment().AwsAccessKey(),
			s.ServiceManager().Environment().AwsSecretKey(),
			"",
		)

		cfg, err := config.LoadDefaultConfig(
			s.ctx,
			config.WithRegion(s.ServiceManager().Environment().AwsRegion()),
			config.WithHTTPClient(client),
			config.WithEndpointResolver(customResolver),
			config.WithCredentialsProvider(staticCredentialsProvider),
		)
		if err != nil {
			return err
		}

		s.client = sns.NewFromConfig(cfg, sns.WithEndpointResolver(snsEndpointResolve{customResolver}))
	}
	return nil
}

func (s *eventsService) Close() error {
	s.client = nil
	return nil
}

func (s *eventsService) Readiness(ctx context.Context) error {
	_, err := s.client.GetTopicAttributes(ctx, &sns.GetTopicAttributesInput{
		TopicArn: aws.String(s.address),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *eventsService) Liveness(ctx context.Context) error {
	return s.Readiness(ctx)
}

func (s *eventsService) WithServiceManager(c services.ServiceManager) services.Events {
	s.serviceManager = c
	return s
}

func (s *eventsService) ServiceManager() services.ServiceManager {
	return s.serviceManager
}

func (s *eventsService) Send(ctx context.Context, req domains.EventRequest) error {
	ctx, span := s.ServiceManager().Spans().New(ctx)
	defer span.Finish()

	if err := s.ServiceManager().Validator().Validate(ctx, req); err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during validating event with err=%v", err))
		span.Error(err)
		return err
	}

	message, err := json.Marshal(req)
	if err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Was unable to marshall message with err=%v", err))
		span.Error(err)
		return err
	}

	in := &sns.PublishInput{
		TopicArn: aws.String(s.address),
		Message:  aws.String(string(message)),
		MessageAttributes: map[string]types.MessageAttributeValue{
			"type": {
				DataType:    ptrs.String("String"),
				StringValue: ptrs.String(req.EventType()),
			},
			"operation": {
				DataType:    ptrs.String("String"),
				StringValue: ptrs.String(req.Operation()),
			},
		},
		MessageGroupId: ptrs.String(req.EventId()),
	}

	res, err := s.client.Publish(ctx, in)
	if err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during Publish on sns err=%v", err))
		span.Error(err)
		return err
	}

	s.ServiceManager().Logger().Info(ctx, fmt.Sprintf("Event sent eventId=%v, messageId=%v", req.EventId(), ptrs.StringValue(res.MessageId)))

	return nil
}

func (s *eventsService) SendBulk(ctx context.Context, req []domains.EventRequest) error {
	ctx, span := s.ServiceManager().Spans().New(ctx)
	defer span.Finish()

	if err := s.ServiceManager().Validator().ValidateSlice(ctx, req); err != nil {
		span.Error(err)
		return err
	}

	chunks := s.chunks(req, 500)

	var wg sync.WaitGroup
	wg.Add(len(chunks))
	for _, items := range chunks {
		go s.sendBulk(ctx, &wg, items)
	}
	wg.Wait()
	chunks = nil

	return nil
}

func (s *eventsService) chunks(eventRequests []domains.EventRequest, lim int) [][]domains.EventRequest {
	var chunk []domains.EventRequest
	chunks := make([][]domains.EventRequest, 0, len(eventRequests)/lim+1)
	for len(eventRequests) >= lim {
		chunk, eventRequests = eventRequests[:lim], eventRequests[lim:]
		chunks = append(chunks, chunk)
	}
	if len(eventRequests) > 0 {
		chunks = append(chunks, eventRequests[:])
	}
	return chunks
}

func (s *eventsService) sendBulk(ctx context.Context, wg *sync.WaitGroup, items []domains.EventRequest) {
	ctx, span := s.ServiceManager().Spans().New(ctx)
	defer span.Finish()
	defer wg.Done()
	for _, req := range items {
		message, err := json.Marshal(req)
		if err != nil {
			s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Was unable to marshall message with err=%v", err))
			span.Error(err)
			continue
		}

		in := &sns.PublishInput{
			TopicArn: aws.String(s.address),
			Message:  aws.String(string(message)),
			MessageAttributes: map[string]types.MessageAttributeValue{
				"type": {
					DataType:    ptrs.String("String"),
					StringValue: ptrs.String(req.EventType()),
				},
				"operation": {
					DataType:    ptrs.String("String"),
					StringValue: ptrs.String(req.Operation()),
				},
			},
			MessageGroupId: ptrs.String(req.EventId()),
		}

		res, err := s.client.Publish(ctx, in)
		if err != nil {
			s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during Publish on sns err=%v", err))
			span.Error(err)
			continue
		}

		s.ServiceManager().Logger().Info(ctx, fmt.Sprintf("Event sent eventId=%v, messageId=%v", req.EventId(), ptrs.StringValue(res.MessageId)))
	}
}
