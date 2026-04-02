package app_test

import (
	"context"
	"testing"

	"go.uber.org/zap"

	"github.com/neofyis/geopulse/internal/app"
	"github.com/neofyis/geopulse/internal/config"
	"github.com/neofyis/geopulse/internal/mocks"
)

func testConfig() *config.Config {
	return &config.Config{
		Port:        "8080",
		ScrapeTypes: "restaurant",
	}
}

func newApp(s *mocks.MockStore) *app.App {
	return app.New(s, nil, testConfig(), zap.NewNop())
}

func TestReady_NilPinger_ReturnsError(t *testing.T) {
	a := app.New(mocks.NewMockStore(t), nil, testConfig(), zap.NewNop())
	if err := a.Ready(context.Background()); err == nil {
		t.Error("expected error when pinger is nil")
	}
}
