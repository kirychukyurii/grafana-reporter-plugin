package service

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
	"path/filepath"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"

	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/adapter/grafana"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/adapter/store"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/config"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/domain/entity"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/cdp"
	gutils "github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/grafana"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/util"
)

type ReportService interface {
	Report(ctx context.Context, id int) (*entity.Report, error)
	NewReport(ctx context.Context, report entity.Report) error
}

type Report struct {
	settings config.ReporterAppSetting

	store             store.ReportStoreManager
	grafanaHTTPClient grafana.GrafanaHTTPAdapter

	browserPool cdp.BrowserPoolManager
	browser     cdp.BrowserManager
	pagePool    cdp.PagePoolManager
	page        cdp.PageManager
}

func NewReportService(settings config.ReporterAppSetting,
	reportStore store.ReportStoreManager,
	grafanaHTTPClient grafana.GrafanaHTTPAdapter,
	browserPool cdp.BrowserPoolManager) *Report {
	return &Report{
		settings:          settings,
		store:             reportStore,
		grafanaHTTPClient: grafanaHTTPClient,
		browserPool:       browserPool,
	}
}

func (r *Report) Report(ctx context.Context, id int) (*entity.Report, error) {
	_, err := r.store.Report(ctx, id)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *Report) NewReport(ctx context.Context, report entity.Report) error {
	defer r.browserPool.Cleanup()

	dashboard, err := r.grafanaHTTPClient.Dashboard(ctx, entity.DashboardOpts{
		DashboardID: report.Dashboard.UID,
		// Variables:   r.Variables,
	})
	if err != nil {
		return fmt.Errorf("grafanaHTTPClient.Dashboard: %v", err)
	}

	tmpDir, err := util.TemporaryDir(r.settings.TemporaryDirectory)
	if err != nil {
		return fmt.Errorf("ensure dir RW: %v", err)
	}

	r.browser, err = r.browserPool.Get(r.settings)
	if err != nil {
		return err
	}
	defer r.browserPool.Put(r.browser)

	workers := util.Workers(r.settings.WorkersCount, len(dashboard.Model.Panels))
	r.pagePool = cdp.NewPagePool(workers)
	defer func(pagePool cdp.PagePoolManager) {
		if tmpErr := pagePool.Cleanup(); tmpErr != nil {
			err = tmpErr
		}
	}(r.pagePool)

	if err = r.exportDashboardPNG(ctx, dashboard, tmpDir); err != nil {
		return err
	}

	eg, gctx := errgroup.WithContext(ctx)
	eg.SetLimit(workers)
	for _, panel := range dashboard.Model.Panels {
		panel := panel // https://golang.org/doc/faq#closures_and_goroutines

		eg.Go(func() error {
			backend.Logger.Debug("processing panel", "panel.id", panel.Id, "panel.title", panel.Title, "panel.type", panel.Type, "output", tmpDir)

			if err = r.exportPanelPNG(gctx, dashboard, panel, tmpDir); err != nil {
				return fmt.Errorf("exportPanelPNG: %v", err)
			}

			//time.Sleep(5 * time.Second)

			return nil
		})
	}

	if err = eg.Wait(); err != nil {
		return err
	}

	return nil

}

func (r *Report) export(ctx context.Context, dashboard *gutils.Dashboard) (err error) {

	tmpDir := filepath.Join(r.settings.TemporaryDirectory, util.NewUUID().String())
	if err = util.EnsureDirRW(tmpDir); err != nil {
		return fmt.Errorf("ensure dir RW: %v", err)
	}

	r.browser, err = r.browserPool.Get(r.settings)
	if err != nil {
		return err
	}
	defer r.browserPool.Put(r.browser)

	panelsCnt := len(dashboard.Model.Panels)
	workers := r.settings.WorkersCount
	if workers > panelsCnt {
		workers = panelsCnt
	}

	r.pagePool = cdp.NewPagePool(workers)
	defer func(pagePool cdp.PagePoolManager) {
		if tmpErr := pagePool.Cleanup(); tmpErr != nil {
			err = tmpErr
		}
	}(r.pagePool)

	if err = r.exportDashboardPNG(ctx, dashboard, tmpDir); err != nil {
		return err
	}

	// sem is a weighted semaphore allowing up to 10 concurrent operations.
	var sem = semaphore.NewWeighted(int64(workers))
	errs := make(chan error, panelsCnt)
	for _, panel := range dashboard.Model.Panels {
		if err = sem.Acquire(ctx, 1); err != nil {
			errs <- fmt.Errorf("acquire semaphore: %v", err)
		}

		go func(dashboard *gutils.Dashboard, panel gutils.Panel, errs chan<- error) {
			defer sem.Release(1)
			backend.Logger.Debug("processing panel", "panel.id", panel.Id, "panel.title", panel.Title, "panel.type", panel.Type, "output", tmpDir)

			if err = r.exportPanelPNG(ctx, dashboard, panel, tmpDir); err != nil {
				errs <- fmt.Errorf("exportPanelPNG: %v", err)
			}

		}(dashboard, panel, errs)

	}

	go func(ctx context.Context) {
		select {
		case err = <-errs:
			if err != nil {
				backend.Logger.Error("error", "error", err)
				ctx.Done()
			}
		}
	}(ctx)

	if err = sem.Acquire(ctx, int64(workers)); err != nil {
		return fmt.Errorf("acquire semaphore: %v", err)
	}

	close(errs)
	for err = range errs {
		if err != nil {
			return err
		}
	}

	/*

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

		s.browser, err = s.browserPool.Get(s.settings)
		if err != nil {
			return err
		}
		defer s.browserPool.Put(s.browser)

		s.pagePool = cdp.NewPagePool(workers)
		defer func(pagePool cdp.PagePoolManager) {
			if tmpErr := pagePool.Cleanup(); tmpErr != nil {
				err = tmpErr
			}
		}(s.pagePool)

		if err = s.exportDashboardPNG(dashboard, tmpDir); err != nil {
			return err
		}

		g, _ := errgroup.WithContext(ctx)
		for i := 0; i < workers; i++ {
			g.Go(func() error {
				for panel := range panels {
					backend.Logger.Debug("processing panel", "panel.id", panel.Id, "panel.title", panel.Title, "panel.type", panel.Type, "output", tmpDir)

					if err = s.exportPanelPNG(dashboard, panel, tmpDir); err != nil {
						return err
					}


						if panel.IsTable() {
							if err = s.exportCSV(dashboard, panel, tmpDir); err != nil {
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
	*/
	return nil
}

func (r *Report) exportDashboardPNG(ctx context.Context, dashboard *gutils.Dashboard, tmpDir string) error {
	since := time.Now()
	defer func() { backend.Logger.Debug(util.TimeTrack(since)) }()

	page, err := r.pagePool.Get(r.browser)
	if err != nil {
		return err
	}
	defer r.pagePool.Put(page)

	url := fmt.Sprintf("%s/d/%s/db?kiosk&theme=light", r.settings.GrafanaBaseURL, dashboard.Model.Uid)
	headers := []string{"Authorization", r.settings.BasicAuth.String()}
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

func (r *Report) exportPanelPNG(ctx context.Context, dashboard *gutils.Dashboard, panel gutils.Panel, tmpDir string) error {
	//return fmt.Errorf("hello")

	since := time.Now()
	defer func() { backend.Logger.Debug(util.TimeTrack(since), "panel.id", panel.Id, "panel.title", panel.Title) }()

	page, err := r.pagePool.Get(r.browser)
	if err != nil {
		return err
	}
	defer r.pagePool.Put(page)

	url := fmt.Sprintf("%s/d-solo/%s/db?panelId=%d&width=%d&height=%d&render=1", r.settings.GrafanaBaseURL, dashboard.Model.Uid, panel.Id, panel.Width(), panel.Height())
	headers := []string{"Authorization", r.settings.BasicAuth.String()}
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

func (r *Report) exportCSV(ctx context.Context, dashboard *gutils.Dashboard, panel gutils.Panel, tmpDir string) error {
	since := time.Now()
	defer func() { backend.Logger.Debug(util.TimeTrack(since)) }()

	page, err := r.pagePool.Get(r.browser)
	if err != nil {
		return err
	}
	defer r.pagePool.Put(page)

	url := fmt.Sprintf("%s/d/%s/db?orgId=1&inspect=%d&inspectTab=data", r.settings.GrafanaBaseURL, dashboard.Model.Uid, panel.Id)
	headers := []string{"Authorization", r.settings.BasicAuth.String()}
	if err = page.Prepare(url, headers, nil); err != nil {
		return err
	}

	waitDownload := r.browser.WaitDownload(tmpDir)
	e, err := page.Element("span", "Download CSV")
	if err != nil {
		return err
	}

	if err = e.Click(cdp.InputMouseButtonLeft, cdp.InputMouseButtonSingleClick); err != nil {
		return err
	}

	if err = cdp.OutputFile(filepath.Join(tmpDir, fmt.Sprintf("panel-%s-%d.csv", panel.Type, panel.Id)), waitDownload()); err != nil {
		return err
	}

	return nil
}
