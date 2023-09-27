package plugin

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource/httpadapter"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/cron"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/store/boltdb"
	"net/http"
	"time"

	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/config"
	cronhandler "github.com/kirychukyurii/grafana-reporter-plugin/pkg/handler/cron"
	handler "github.com/kirychukyurii/grafana-reporter-plugin/pkg/handler/http"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/cdp"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/grafana"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/log"
)

// Make sure App implements required interfaces.
// This is important to do since otherwise we will only get a
// not implemented error response from plugin in runtime.
var (
	_ backend.CallResourceHandler   = (*App)(nil)
	_ backend.CheckHealthHandler    = (*App)(nil)
	_ instancemgmt.InstanceDisposer = (*App)(nil)
)

type App struct {
	settings *config.ReporterAppConfig
	router   *mux.Router
	handler  handler.HandlerManager
	cron     cronhandler.ReportScheduleCronHandler
}

// New creates a new *App instance.
func New(ctx context.Context, s backend.AppInstanceSettings) (instancemgmt.Instance, error) {
	setting, err := config.New(s)
	if err != nil {
		return nil, err
	}

	logger := log.New()

	database, err := boltdb.New(setting, logger)
	if err != nil {
		return nil, err
	}

	pool := cdp.NewBrowserPool(setting)
	grafanaClient, err := grafana.New(setting)
	if err != nil {
		return nil, err
	}

	timezone, err := time.LoadLocation("Local")
	if err != nil {
		return nil, err
	}

	scheduler := cron.NewScheduler(timezone)

	app, err := Initialize(setting, database, logger, grafanaClient, pool, scheduler)
	if err != nil {
		return nil, err
	}

	if err = app.cron.LoadSchedules(); err != nil {
		return nil, err
	}

	app.registerRoutes()

	return app, nil
}

func newApp(setting *config.ReporterAppConfig, handler handler.HandlerManager, cronHandler cronhandler.ReportScheduleCronHandler) (*App, error) {
	return &App{
		settings: setting,
		router:   mux.NewRouter(),
		handler:  handler,
		cron:     cronHandler,
	}, nil
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}

// CallResource HTTP style resource
func (a *App) CallResource(ctx context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error {
	return httpadapter.New(a.router).CallResource(ctx, req, sender)
}

// CheckHealth handles health checks sent from Grafana to the plugin.
func (a *App) CheckHealth(_ context.Context, _ *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: "ok",
	}, nil
}

// Dispose here tells plugin SDK that plugin wants to clean up resources when a new instance
// created.
func (a *App) Dispose() {
	backend.Logger.Info("called when the settings change", "cfg", a.settings)
}
