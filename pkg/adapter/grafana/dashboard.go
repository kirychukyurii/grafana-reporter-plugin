package grafana

import (
	"context"

	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/domain/entity"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/grafana"
)

type DashboardAdapter interface {
	Dashboard(ctx context.Context, opts entity.DashboardOpts) (*grafana.Dashboard, error)
}

type DashboardClient struct {
	connection *grafana.Client
}

func NewDashboardClient(client *grafana.Client) *DashboardClient {
	return &DashboardClient{
		connection: client,
	}
}

func (c *DashboardClient) Dashboard(ctx context.Context, opts entity.DashboardOpts) (*grafana.Dashboard, error) {
	dashboard, err := c.connection.Dashboard(ctx, opts)
	if err != nil {
		return nil, err
	}

	return dashboard, nil
}
