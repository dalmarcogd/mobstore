package validator

import (
	"context"
	goerrors "errors"
	"reflect"

	"github.com/go-playground/validator/v10"

	"github.com/dalmarcogd/mobstore/users/internal/infra/errors"
	"github.com/dalmarcogd/mobstore/users/internal/services"
)

type (
	validatorService struct {
		services.NoopHealth
		serviceManager services.ServiceManager
		ctx            context.Context
		validate       *validator.Validate
	}
)

func New() *validatorService {
	return &validatorService{}
}

func (s *validatorService) Init(ctx context.Context) error {
	s.ctx = ctx
	if s.validate == nil {
		s.validate = validator.New()
	}

	return nil
}

func (s *validatorService) Close() error {
	return nil
}

func (s *validatorService) WithServiceManager(c services.ServiceManager) services.Validator {
	s.serviceManager = c
	return s
}

func (s *validatorService) ServiceManager() services.ServiceManager {
	return s.serviceManager
}

func (s *validatorService) Validate(ctx context.Context, obj interface{}) error {
	err := s.validate.StructCtx(ctx, obj)
	if err != nil {
		return err
	}
	return nil
}

func (s *validatorService) ValidateSlice(ctx context.Context, objs interface{}) error {
	sv := reflect.ValueOf(objs)
	if sv.Kind() != reflect.Slice || sv.IsNil() {
		return errors.ObjsIsNotSliceValidatorError
	}

	ret := make([]interface{}, sv.Len())
	for i := 0; i < sv.Len(); i++ {
		ret[i] = sv.Index(i).Interface()
	}

	var validationErrors validator.ValidationErrors
	for _, obj := range ret {
		err := s.Validate(ctx, obj)
		if err != nil {
			if goerrors.Is(err, validator.ValidationErrors{}) {
				validationErrors = append(validationErrors, err.(validator.ValidationErrors)...)
			} else {
				return err
			}
		}
	}
	ret = nil
	if len(validationErrors) > 0 {
		return validationErrors
	}
	return nil
}
