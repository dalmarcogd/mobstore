package environment

import (
	"testing"

	"github.com/dalmarcogd/mobstore/products/internal/services"
)

func TestLoadEnvsSuccessfully(t *testing.T) {
	serviceImpl := New()
	sm := services.New().WithEnvironment(serviceImpl)

	if err := sm.Init(); err != nil {
		t.Error("unexpected address")
	}

	if sm.Environment().Env() == "" {
		t.Error("Env env should be default value")
	}
	if sm.Environment().Service() == "" {
		t.Error("envService env should be default value")
	}
	if sm.Environment().Version() == "" {
		t.Error("Version env should be default value")
	}
	if sm.Environment().ProductDatabaseDsn() == "" {
		t.Error("ProductDatabaseDsn env should be default value")
	}
	if sm.Environment().ProductReplicaDatabaseDsn() == "" {
		t.Error("ProductDatabaseDsn env should be default value")
	}
	if sm.Environment().AwsRegion() == "" {
		t.Error("AwsRegion env should be default value")
	}
	if !sm.Environment().DebugPprof() {
		t.Error("DebugPprof env should be default value")
	}
	if sm.Environment().DiscountsAddress() == "" {
		t.Error("DiscountsAddress env should be default value")
	}
	if sm.Environment().SpanUrl() == "" {
		t.Error("SpanUrl env should be default value")
	}
	if err := serviceImpl.ServiceManager().Close(); err != nil {
		t.Error(err)
	}
}
