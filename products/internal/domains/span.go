package domains

import (
	"context"
	"fmt"
	"strconv"

	"github.com/openzipkin/zipkin-go"
)

type (
	Span struct {
		Name         string                 `validate:"required"`
		Cid          string                 `validate:"required"`
		Resource     string                 `validate:"required"`
		Version      string                 `validate:"required"`
		OrgId        string                 `validate:"required"`
		Line         int                    `validate:"required"`
		FuncName     string                 `validate:"required"`
		FileName     string                 `validate:"required"`
		Custom       map[string]interface{} `validate:"required"`
		InternalSpan zipkin.Span
	}
	SpanConfig interface {
		Apply(ctx context.Context, s *Span)
	}
	nameSpanConfigStruct struct {
		value string
	}
	orgIdSpanConfigStruct struct {
		value string
	}
	customSpanConfigStruct struct {
		key   string
		value interface{}
	}
)

func (s *Span) Tag(k string, v interface{}) *Span {
	s.Custom[k] = v
	return s
}

func (s *Span) Finish() {
	if s.InternalSpan != nil {
		s.InternalSpan.Tag("cid", s.Cid)
		s.InternalSpan.Tag("line", strconv.Itoa(s.Line))
		s.InternalSpan.Tag("func-name", s.FuncName)
		s.InternalSpan.Tag("file-name", s.FileName)
		for k, v := range s.Custom {
			s.InternalSpan.Tag(k, fmt.Sprintf("%v", v))
		}
		s.InternalSpan.Finish()
	}
	s.Custom = nil
}

func (s *Span) Error(err error) *Span {
	if err != nil {
		s.Tag("error", true)
		s.Tag("error.message", err)
	}
	return s
}

func WithName(v string) SpanConfig {
	return &nameSpanConfigStruct{value: v}
}

func (n *nameSpanConfigStruct) Apply(_ context.Context, s *Span) {
	s.Name = n.value
}

func WithOrgId(v string) SpanConfig {
	return &orgIdSpanConfigStruct{value: v}
}

func (n *orgIdSpanConfigStruct) Apply(_ context.Context, s *Span) {
	s.OrgId = n.value
}

func WithCustom(k string, v interface{}) SpanConfig {
	return &customSpanConfigStruct{key: k, value: v}
}

func (n *customSpanConfigStruct) Apply(_ context.Context, s *Span) {
	s.Custom[n.key] = n.value
}
