package service

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"golang.org/x/sync/errgroup"

	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/adapter/grafana"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/config"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/domain/entity"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/cdp"
	gutils "github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/grafana"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/log"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/util"
)

type ReportService interface {
	Report(ctx context.Context, id int) (*entity.Report, error)
	NewReport(ctx context.Context, report entity.Report) (string, error)
}

type Report struct {
	settings *config.ReporterAppConfig
	logger   *log.Logger

	grafanaHTTPClient grafana.DashboardAdapter

	browserPool cdp.BrowserPoolManager
	browser     cdp.BrowserManager
	pagePool    cdp.PagePoolManager
	page        cdp.PageManager
}

func NewReportService(settings *config.ReporterAppConfig, logger *log.Logger, grafanaHTTPClient grafana.DashboardAdapter, browserPool cdp.BrowserPoolManager) *Report {
	subLogger := &log.Logger{
		Logger: logger.With("component.type", "service", "component", "report"),
	}

	return &Report{
		settings:          settings,
		logger:            subLogger,
		grafanaHTTPClient: grafanaHTTPClient,
		browserPool:       browserPool,
	}
}

func (r *Report) Report(ctx context.Context, id int) (*entity.Report, error) {

	return nil, nil
}

func (r *Report) NewReport(ctx context.Context, report entity.Report) (string, error) {
	dashboard, err := r.grafanaHTTPClient.Dashboard(ctx, entity.DashboardOpts{
		DashboardID: report.Dashboard.UID,
		// Variables:   r.Variables,
	})
	if err != nil {
		return "", fmt.Errorf("grafanaHTTPClient.Dashboard: %v", err)
	}

	tmpDir, err := util.TemporaryDir(filepath.Join(r.settings.DataDirectory, "file"))
	if err != nil {
		return "", fmt.Errorf("ensure dir RW: %v", err)
	}

	r.browser, err = r.browserPool.Get(r.settings)
	if err != nil {
		return "", err
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
		return "", err
	}
	eg, gctx := errgroup.WithContext(ctx)
	eg.SetLimit(workers)
	for _, panel := range dashboard.Model.Panels {
		panel := panel // https://golang.org/doc/faq#closures_and_goroutines

		eg.Go(func() error {
			r.logger.Debug("processing panel", "panel.id", panel.Id, "panel.title", panel.Title, "panel.type", panel.Type, "output", tmpDir)
			if err = r.exportPanelPNG(gctx, dashboard, panel, tmpDir); err != nil {
				return fmt.Errorf("exportPanelPNG: %v", err)
			}

			return nil
		})
	}

	if err = eg.Wait(); err != nil {
		return "", err
	}

	return tmpDir, nil

}

func (r *Report) exportDashboardPNG(ctx context.Context, dashboard *gutils.Dashboard, tmpDir string) error {
	since := time.Now()
	defer func() { backend.Logger.Debug(util.TimeTrack(since)) }()

	page, err := r.pagePool.Get(r.browser)
	if err != nil {
		return err
	}
	defer r.pagePool.Put(page)

	url := fmt.Sprintf("%s/d/%s/db?kiosk&theme=light", r.settings.GrafanaConfig.URL, dashboard.Model.Uid)
	headers := []string{"Authorization", r.settings.GrafanaConfig.BasicAuth()}
	if err = page.Prepare(url, headers, nil); err != nil {
		return err
	}

	if err = page.ScrollDown(500); err != nil {
		return err
	}

	panels, err := page.Elements(`[data-panelId]`)
	if err != nil {
		return fmt.Errorf("data-panelId elements: %v", err)
	}

	panelCount := len(panels)
	panelsRendered, err := page.Elements(`[class$='panel-content']`)
	if err != nil {
		return fmt.Errorf("panel-content elements: %v", err)
	}

	backend.Logger.Debug("panels", "panelCount", panelCount, "panelsRendered", panelsRendered)
	panelsRenderedCount := 0
	for i, p := range panelsRendered {
		// https://github.com/grafana/grafana-image-renderer/blob/master/src/browser/browser.ts#L344
		r.logger.Debug("panelRenderedCount", "i", i, "panelRendered", p)
		panelsRenderedCount++
	}

	dRow, err := page.Elements(`'.dashboard-row'`)
	if err != nil {
		return err
	}

	totalPanelsRendered := panelsRenderedCount + len(dRow)
	if totalPanelsRendered >= panelCount {
		if err = page.ScreenshotFullPage(filepath.Join(tmpDir, fmt.Sprintf("%s.png", dashboard.Model.Uid))); err != nil {
			return err
		}
	}

	return nil
}

func (r *Report) exportPanelPNG(ctx context.Context, dashboard *gutils.Dashboard, panel gutils.Panel, tmpDir string) error {
	since := time.Now()
	defer func() { backend.Logger.Debug(util.TimeTrack(since), "panel.id", panel.Id, "panel.title", panel.Title) }()

	page, err := r.pagePool.Get(r.browser)
	if err != nil {
		return err
	}
	defer r.pagePool.Put(page)

	url := fmt.Sprintf("%s/d-solo/%s/db?panelId=%d&width=%d&height=%d&render=1", r.settings.GrafanaConfig.URL, dashboard.Model.Uid, panel.Id, panel.Width(), panel.Height())
	headers := []string{"Authorization", r.settings.GrafanaConfig.BasicAuth()}
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

	url := fmt.Sprintf("%s/d/%s/db?orgId=1&viewPanel=%d&inspect=%d&inspectTab=data", r.settings.GrafanaConfig.URL, dashboard.Model.Uid, panel.Id, panel.Id)
	headers := []string{"Authorization", r.settings.GrafanaConfig.BasicAuth()}
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
