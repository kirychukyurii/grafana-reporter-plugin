package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
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

	pagePool := rod.NewPagePool(workers)

	defer func(pagePool rod.PagePool) {
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

func (s *service) exportDashboardPNG(dashboard *gutils.Dashboard, tmpDir string, pagePool rod.PagePool, b *rod.Browser) error {
	page, err := pagePool.Get(b)
	if err != nil {
		return err
	}
	defer pagePool.Put(page)

	_, err = page.SetExtraHeaders([]string{"Authorization", s.settings.BasicAuth.String()})
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/d/%s/db?kiosk&theme=light", s.settings.GrafanaBaseURL, dashboard.Model.Uid)
	if err = page.Navigate(url); err != nil {
		return err
	}

	if err = page.WaitLoad(); err != nil {
		return err
	}

	w := page.MustWaitRequestIdle()
	w()

	scrollHeightObj, err := page.Eval("() => document.documentElement.scrollHeight")
	if err != nil {
		return fmt.Errorf("scrollHeightObj: %v", err)
	}

	clientHeightObj, err := page.Eval("() => document.documentElement.clientHeight")
	if err != nil {
		return fmt.Errorf("clientHeightObj: %v", err)
	}

	/*
		scrollHeightObj, err := page.Evaluate(&rod.EvalOptions{
			JS: "document.documentElement.scrollHeight",
		})

		clientHeightObj, err := page.Evaluate(&rod.EvalOptions{
			JS: "document.documentElement.clientHeight",
		})
	*/

	backend.Logger.Debug("heights", "scrollHeight", scrollHeightObj, "clientHeight", clientHeightObj)
	scrollHeight := scrollHeightObj.Value.Num()
	clientHeight := clientHeightObj.Value.Num()
	backend.Logger.Debug("heights", "scrollHeight", scrollHeight, "clientHeight", clientHeight)

	if scrollHeight < clientHeight {

	}

	scrolls := int(scrollHeight / clientHeight)
	for i := 1; i < scrolls; i++ {
		if err = page.Mouse.Scroll(0, clientHeight, 0); err != nil {
			return fmt.Errorf("scroll: %v", err)
		}

		time.Sleep(500 * time.Millisecond)
	}

	if err = page.Mouse.Scroll(0, 0, 0); err != nil {
		return fmt.Errorf("scroll to 0,0: %v", err)
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

	page.MustScreenshotFullPage(filepath.Join(tmpDir, fmt.Sprintf("%s.png", dashboard.Model.Uid)))

	return nil
}

func (s *service) exportPanelPNG(dashboard *gutils.Dashboard, panel gutils.Panel, pagePool rod.PagePool, b *rod.Browser, tmpDir string) error {
	page, err := pagePool.Get(b)
	if err != nil {
		return err
	}
	defer pagePool.Put(page)

	if err = page.SetViewport(&proto.EmulationSetDeviceMetricsOverride{Width: panel.Width(), Height: panel.Height()}); err != nil {
		return err
	}

	_, err = page.SetExtraHeaders([]string{"Authorization", s.settings.BasicAuth.String()})
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/d-solo/%s/db?panelId=%d&width=%d&height=%d&render=1", s.settings.GrafanaBaseURL, dashboard.Model.Uid, panel.Id, panel.Width(), panel.Height())
	if err = page.Navigate(url); err != nil {
		return err
	}

	if err = page.WaitLoad(); err != nil {
		return err
	}

	w := page.MustWaitRequestIdle()
	w()
	/*
		waitReqIdle := page.WaitRequestIdle(300*time.Millisecond, nil, nil, nil)
		waitReqIdle()
	*/

	page.MustScreenshot(filepath.Join(tmpDir, fmt.Sprintf("panel-%s-%d.png", panel.Type, panel.Id)))
	// page.MustPDF(tmpDir + "test.pdf")
	/*
			screenshot, err := page.Screenshot(false, nil)
			if err != nil {
				return err
			}

		if err = utils.OutputFile(tmpDir+"test.png", screenshot); err != nil {
			return err
		}
	*/

	return nil
}

func (s *service) exportCSV(dashboard *gutils.Dashboard, panel gutils.Panel, pagePool rod.PagePool, b *rod.Browser, tmpDir string) error {
	page, err := pagePool.Get(b)
	if err != nil {
		return err
	}
	defer pagePool.Put(page)

	/*// https://cloud.webitel.ua/grafana/d/cl1CQ2Gnk/operatory?orgId=1&inspect=4&inspectTab=data

	page, err := pagePool.Get(b, "https://webhook.site/9e65782a-732f-4809-aa23-417eb8e830a1")
	if err != nil {
		return err
	}
	defer page.Close()
	*/

	url := fmt.Sprintf("%s/d/%s/db?orgId=1&inspect=%d&inspectTab=data", s.settings.GrafanaBaseURL, dashboard.Model.Uid, panel.Id)
	if err = page.Navigate(url); err != nil {
		return err
	}

	if err = page.WaitLoad(); err != nil {
		return err
	}

	w := page.MustWaitRequestIdle()
	w()

	waitDownload := b.WaitDownload(tmpDir)
	e, err := page.ElementR("span", "Download CSV")
	if err != nil {
		return err
	}

	if err := e.Click(proto.InputMouseButtonLeft, 1); err != nil {
		return err
	}

	_ = utils.OutputFile(filepath.Join(tmpDir, fmt.Sprintf("panel-%s-%d.csv", panel.Type, panel.Id)), waitDownload())

	return nil
}
