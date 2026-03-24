package app_test

import (
	"context"
	"testing"

	"go.uber.org/zap"

	"github.com/neofyis/geopulse/bff/internal/app"
	"github.com/neofyis/geopulse/bff/internal/mocks"
)

func newApp(s *mocks.MockStore) *app.App {
	return app.New(s, nil, zap.NewNop())
}

func TestReady_NilPinger_ReturnsError(t *testing.T) {
	a := app.New(mocks.NewMockStore(t), nil, zap.NewNop())
	if err := a.Ready(context.Background()); err == nil {
		t.Error("expected error when pinger is nil")
	}
}
