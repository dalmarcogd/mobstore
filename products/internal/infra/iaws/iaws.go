package iaws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type (
	SNSClient interface {
		GetTopicAttributes(ctx context.Context, params *sns.GetTopicAttributesInput, optFns ...func(*sns.Options)) (*sns.GetTopicAttributesOutput, error)
		Publish(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error)
	}
)
