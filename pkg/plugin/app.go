package plugin

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource/httpadapter"

	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/config"
	handler "github.com/kirychukyurii/grafana-reporter-plugin/pkg/handler/http"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/cdp"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/cron"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/db"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/db/migration"
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
	settings *config.ReporterAppSetting
	router   *mux.Router
	handler  handler.ReportHandler
}

// New creates a new *App instance.
func New(s backend.AppInstanceSettings) (instancemgmt.Instance, error) {
	setting, err := config.New(s)
	if err != nil {
		return nil, err
	}

	logger := log.New()
	database, err := db.New()
	if err != nil {
		return nil, err
	}

	if err = migration.Migrate(database); err != nil {
		return nil, fmt.Errorf("migrate: %v", err)
	}

	pool := cdp.NewBrowserPool(setting)
	grafanaClient, err := grafana.New(setting)
	if err != nil {
		return nil, err
	}

	app, err := Initialize(setting, database, logger, grafanaClient, pool)
	if err != nil {
		return nil, err
	}

	app.registerRoutes()
	return app, nil
}

func newApp(setting *config.ReporterAppSetting, handler handler.ReportHandler) (*App, error) {
	return &App{
		settings: setting,
		router:   mux.NewRouter(),
		handler:  handler,
	}, nil
}

func runScheduler() error {
	timezone, err := time.LoadLocation("Local")
	if err != nil {
		return err
	}

	scheduler := cron.NewScheduler(timezone)
	scheduler.Cron.StartAsync()

	return nil
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
