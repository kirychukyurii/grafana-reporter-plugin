package grafana

import (
	"context"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/grafana"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/models"
)

type GrafanaHTTPAdapter interface {
	Dashboard(ctx context.Context, opts models.DashboardOpts) (*grafana.Dashboard, error)
}

type grafanaHTTPAdapter struct {
	client grafana.Client
}

func New(client grafana.Client) GrafanaHTTPAdapter {
	return &grafanaHTTPAdapter{
		client: client,
	}
}
