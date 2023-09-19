package grafana

import (
	"context"

	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/domain/entity"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/grafana"
)

func (g *grafanaHTTPAdapter) Dashboard(ctx context.Context, opts entity.DashboardOpts) (*grafana.Dashboard, error) {
	dashboard, err := g.client.Dashboard(ctx, opts)
	if err != nil {
		return nil, err
	}

	return dashboard, nil
}
