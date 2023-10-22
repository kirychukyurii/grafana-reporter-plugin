package plugin

import (
	"context"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
)

type AppInstance struct{}

// New creates a new *App instance.
func New(ctx context.Context, s backend.AppInstanceSettings) (instancemgmt.Instance, error) {
	return &AppInstance{}, nil
}
