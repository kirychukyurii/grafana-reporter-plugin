package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/browser"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/model"
	"golang.org/x/sync/errgroup"
	"path/filepath"
	"time"

	gutils "github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/grafana"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/utils"
)

type Reporter interface {
	Report(ctx context.Context, id int) (*model.Report, error)
	NewReport(ctx context.Context, report model.Report) error
}

func (s *service) Report(ctx context.Context, id int) (*model.Report, error) {
	_, err := s.database.Report(ctx, id)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *service) NewReport(ctx context.Context, report model.Report) error {
	g, gctx := errgroup.WithContext(ctx)
	for _, r := range report.Dashboards {
		g.Go(func() error {
			opts := model.DashboardOpts{
				DashboardID: r.Dashboard.UID,
				//Variables:   r.Variables,
			}

			dashboard, err := s.grafanaHTTPClient.Dashboard(ctx, opts)
			if err != nil {
				return err
			}

			if err := export(gctx, s.settings, s.browserPool, dashboard); err != nil {
				return err
			}

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

	return nil
}

func export(ctx context.Context, settings model.ReporterAppSetting, browserPool browser.Browser, dashboard *gutils.Dashboard) error {
	tmpDir := filepath.Join(settings.TemporaryDirectory, utils.NewUUID().String())
	panelsCnt := len(dashboard.Model.Panels)

	// fetch images in parallel form Grafana sever.
	// limit concurrency using a worker pool to avoid overwhelming grafana
	// for dashboards with many panels.
	workers := settings.WorkersCount
	if workers > panelsCnt {
		workers = panelsCnt
	}

	panels := make(chan gutils.Panel, panelsCnt)
	for _, p := range dashboard.Model.Panels {
		panels <- p
	}
	close(panels)

	b, err := browserPool.Get(settings)
	if err != nil {
		return err
	}

	pagePool := browser.NewPagePool(workers)
	g, _ := errgroup.WithContext(ctx)
	for i := 0; i < workers; i++ {
		g.Go(func() error {
			for panel := range panels {
				if !panel.IsTable() {
					if err := exportPNG(pagePool, b, tmpDir); err != nil {
						return err
					}
				} else {
					if err := exportCSV(pagePool, b, tmpDir); err != nil {
						return err
					}
				}
			}

			return nil
		})
	}

	return g.Wait()
}

func exportPNG(pagePool browser.PagePool, b *rod.Browser, tmpDir string) error {
	return nil
}

func exportCSV(pagePool browser.PagePool, b *rod.Browser, tmpDir string) error {
	page, err := pagePool.Get(b, "https://webhook.site/9e65782a-732f-4809-aa23-417eb8e830a1")
	if err != nil {
		return err
	}
	defer page.Close()

	headers, err := page.SetExtraHeaders([]string{"Authorization", "Basic YWRtaW46ME1QMVFjSm5rQTlC"})
	if err != nil {
		return err
	}

	headers()

	waitReqIdle := page.WaitRequestIdle(300*time.Millisecond, nil, nil, nil)
	waitReqIdle()

	waitDownload := b.WaitDownload(tmpDir)
	e, err := page.ElementR("span", "Download CSV")
	if err != nil {
		return err
	}

	if err := e.Click(proto.InputMouseButtonLeft, 1); err != nil {
		return err
	}

	_ = utils.OutputFile(filepath.Join(tmpDir, "t.csv"), waitDownload())

	return nil
}
