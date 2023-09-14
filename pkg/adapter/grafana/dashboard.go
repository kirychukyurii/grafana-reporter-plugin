package grafana

import (
	"context"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/model"

	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/grafana"
)

func (g *grafanaHTTPAdapter) Dashboard(ctx context.Context, opts model.DashboardOpts) (*grafana.Dashboard, error) {
	dashboard, err := g.client.Dashboard(ctx, opts)
	if err != nil {
		return nil, err
	}

	return dashboard, nil
}
