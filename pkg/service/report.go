package service

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"golang.org/x/sync/errgroup"

	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/cdp"
	gutils "github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/grafana"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/model"
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
			//defer s.browserPool.Cleanup()
			opts := model.DashboardOpts{
				DashboardID: r.Dashboard.UID,
				// Variables:   r.Variables,
			}

			dashboard, err := s.grafanaHTTPClient.Dashboard(ctx, opts)
			if err != nil {
				return fmt.Errorf("grafanaHTTPClient.Dashboard: %v", err)
			}

			if err = s.export(gctx, dashboard); err != nil {
				return fmt.Errorf("export: %v", err)
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		if errors.Is(err, context.Canceled) {
			return fmt.Errorf("context was canceled")
		} else {
			return err
		}
	}

	return nil
}

func (s *service) export(ctx context.Context, dashboard *gutils.Dashboard) (err error) {
	tmpDir := filepath.Join(s.settings.TemporaryDirectory, utils.NewUUID().String())
	if err := utils.EnsureDirRW(tmpDir); err != nil {
		return err
	}

	// fetch images in parallel form Grafana sever.
	// limit concurrency using a worker pool to avoid overwhelming grafana
	// for dashboards with many panels.
	panelsCnt := len(dashboard.Model.Panels)
	workers := s.settings.WorkersCount
	if workers > panelsCnt {
		workers = panelsCnt
	}

	panels := make(chan gutils.Panel, panelsCnt)
	for _, p := range dashboard.Model.Panels {
		panels <- p
	}
	close(panels)

	b, err := s.browserPool.Get(s.settings)
	if err != nil {
		return err
	}
	defer s.browserPool.Put(b)

	pagePool := cdp.NewPagePool(workers)

	defer func(pagePool cdp.PagePoolManager) {
		if tmpErr := pagePool.Cleanup(); tmpErr != nil {
			err = tmpErr
		}
	}(pagePool)

	if err = s.exportDashboardPNG(dashboard, tmpDir, pagePool, b); err != nil {
		return err
	}

	g, _ := errgroup.WithContext(ctx)
	for i := 0; i < workers; i++ {
		g.Go(func() error {
			for panel := range panels {
				/*
					backend.Logger.Debug("processing panel", "panel.id", panel.Id, "panel.title", panel.Title, "panel.type", panel.Type, "output", tmpDir)
					if err = s.exportPanelPNG(dashboard, panel, pagePool, b, tmpDir); err != nil {
						return err
					}
				*/

				if panel.IsTable() {
					if err = s.exportCSV(dashboard, panel, pagePool, b, tmpDir); err != nil {
						return err
					}
				}
			}

			return nil
		})
	}

	if err = g.Wait(); err != nil {
		return err
	}

	return nil
}

func (s *service) exportDashboardPNG(dashboard *gutils.Dashboard, tmpDir string, pagePool cdp.PagePoolManager, b cdp.BrowserManager) error {
	page, err := pagePool.Get(b)
	if err != nil {
		return err
	}
	defer pagePool.Put(page)

	url := fmt.Sprintf("%s/d/%s/db?kiosk&theme=light", s.settings.GrafanaBaseURL, dashboard.Model.Uid)
	headers := []string{"Authorization", s.settings.BasicAuth.String()}
	if err = page.Prepare(url, headers, nil); err != nil {
		return err
	}

	if err = page.ScrollDown(500); err != nil {
		return err
	}

	panels, err := page.Elements("[data-panelId]")
	if err != nil {
		return fmt.Errorf("data-panelId elements: %v", err)
	}

	panelCount := len(panels)
	panelsRendered, err := page.Elements("[class$='panel-content']")
	if err != nil {
		return fmt.Errorf("panel-content elements: %v", err)
	}

	backend.Logger.Debug("panels", "panelCount", panelCount, "panelsRendered", panelsRendered)
	// panelRenderedCount := 0
	for i, p := range panelsRendered {
		// https://github.com/grafana/grafana-image-renderer/blob/master/src/browser/browser.ts#L344
		backend.Logger.Debug("panelRenderedCount", "i", i, "panelRendered", p)
	}

	if err = page.ScreenshotFullPage(filepath.Join(tmpDir, fmt.Sprintf("%s.png", dashboard.Model.Uid))); err != nil {
		return err
	}

	return nil
}

func (s *service) exportPanelPNG(dashboard *gutils.Dashboard, panel gutils.Panel, pagePool cdp.PagePoolManager, b cdp.BrowserManager, tmpDir string) error {
	page, err := pagePool.Get(b)
	if err != nil {
		return err
	}
	defer pagePool.Put(page)

	url := fmt.Sprintf("%s/d-solo/%s/db?panelId=%d&width=%d&height=%d&render=1", s.settings.GrafanaBaseURL, dashboard.Model.Uid, panel.Id, panel.Width(), panel.Height())
	headers := []string{"Authorization", s.settings.BasicAuth.String()}
	viewport := &cdp.PageViewportOpts{
		Width:  panel.Width(),
		Height: panel.Height(),
	}

	if err = page.Prepare(url, headers, viewport); err != nil {
		return err
	}

	if err = page.Screenshot(filepath.Join(tmpDir, fmt.Sprintf("panel-%s-%d.png", panel.Type, panel.Id)), false); err != nil {
		return err
	}

	return nil
}

func (s *service) exportCSV(dashboard *gutils.Dashboard, panel gutils.Panel, pagePool cdp.PagePoolManager, b cdp.BrowserManager, tmpDir string) error {
	page, err := pagePool.Get(b)
	if err != nil {
		return err
	}
	defer pagePool.Put(page)

	url := fmt.Sprintf("%s/d/%s/db?orgId=1&inspect=%d&inspectTab=data", s.settings.GrafanaBaseURL, dashboard.Model.Uid, panel.Id)
	headers := []string{"Authorization", s.settings.BasicAuth.String()}
	if err = page.Prepare(url, headers, nil); err != nil {
		return err
	}

	waitDownload := b.WaitDownload(tmpDir)
	e, err := page.Element("span", "Download CSV")
	if err != nil {
		return err
	}

	if err = e.Click(cdp.InputMouseButtonLeft, 1); err != nil {
		return err
	}

	if err = cdp.OutputFile(filepath.Join(tmpDir, fmt.Sprintf("panel-%s-%d.csv", panel.Type, panel.Id)), waitDownload()); err != nil {
		return err
	}

	return nil
}
