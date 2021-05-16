package environment

import (
	"testing"

	"github.com/dalmarcogd/mobstore/users/internal/services"
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
	if sm.Environment().UserDatabaseDsn() == "" {
		t.Error("UserDatabaseDsn env should be default value")
	}
	if sm.Environment().UserReplicaDatabaseDsn() == "" {
		t.Error("UserDatabaseDsn env should be default value")
	}
	if sm.Environment().AwsRegion() == "" {
		t.Error("AwsRegion env should be default value")
	}
	if !sm.Environment().DebugPprof() {
		t.Error("DebugPprof env should be default value")
	}
	if sm.Environment().SpanUrl() == "" {
		t.Error("SpanUrl env should be default value")
	}
	if err := serviceImpl.ServiceManager().Close(); err != nil {
		t.Error(err)
	}
}
