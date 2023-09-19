package grafana

import (
	"context"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/dto"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/grafana"
)

type GrafanaHTTPAdapter interface {
	Dashboard(ctx context.Context, opts dto.DashboardOpts) (*grafana.Dashboard, error)
}

type grafanaHTTPAdapter struct {
	client grafana.Client
}

func New(client grafana.Client) GrafanaHTTPAdapter {
	return &grafanaHTTPAdapter{
		client: client,
	}
}
