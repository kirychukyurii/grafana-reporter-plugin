package service

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"golang.org/x/sync/errgroup"

	gutils "github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/grafana"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/models"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/utils"
)

type Reporter interface {
	Report(ctx context.Context, reportID int) (*models.Report, error)
	NewReport(ctx context.Context, report models.Report) error
}

func (s *service) Report(ctx context.Context, reportID int) (*models.Report, error) {
	return nil, nil
}

func (s *service) NewReport(ctx context.Context, report models.Report) error {
	for _, r := range report.Dashboards {
		opts := models.DashboardOpts{
			DashboardID: r.Dashboard.UID,
			//Variables:   r.Variables,
		}

		dashboard, err := s.grafanaHTTPClient.Dashboard(ctx, opts)
		if err != nil {
			return err
		}

		_ = filepath.Join(s.settings.TemporaryDirectory, utils.NewUUID().String())

		panelsCnt := len(dashboard.Model.Panels)
		panels := make(chan gutils.Panel, panelsCnt)
		for _, p := range dashboard.Model.Panels {
			panels <- p
		}
		close(panels)

		// fetch images in parallel form Grafana sever.
		// limit concurrency using a worker pool to avoid overwhelming grafana
		// for dashboards with many panels.
		workers := s.settings.WorkersCount
		if workers > panelsCnt {
			workers = panelsCnt
		}

		g, _ := errgroup.WithContext(ctx)
		for i := 0; i < workers; i++ {
			g.Go(func() error {
				return nil
			})
		}

		// wait for all errgroup goroutines
		if err := g.Wait(); err != nil {
			if errors.Is(err, context.Canceled) {
				return fmt.Errorf("context was canceled")
			} else {
				return err
			}
		}
	}

	return nil
}
