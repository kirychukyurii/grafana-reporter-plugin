package plugin

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource/httpadapter"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/adapter/store"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/handler"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/cron"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/db"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/grafana"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/migration"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/models"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/service"
	"net/http"
	"time"
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
	settings models.ReporterAppSetting
	router   *mux.Router
	handler  handler.Handler
}

// New creates a new *App instance.
func New(s backend.AppInstanceSettings) (instancemgmt.Instance, error) {
	backend.Logger.Info("load conf")
	settings := models.ReporterAppSetting{}
	if err := settings.Load(s); err != nil {
		return nil, err
	}

	if err := settings.Validate(); err != nil {
		return nil, err
	}

	return newApp(settings)
}

func newApp(settings models.ReporterAppSetting) (*App, error) {
	database, err := db.New()
	if err != nil {
		return nil, fmt.Errorf("database: %v", err)
	}

	if err = migration.Migrate(database); err != nil {

		return nil, fmt.Errorf("migrate: %v", err)
	}

	gclient, _ := grafana.New(settings)
	s := store.New()
	svc := service.New(settings, s, gclient)
	router := mux.NewRouter()
	app := &App{
		settings: settings,
		router:   router,
		handler:  handler.New(svc),
	}

	app.registerRoutes()

	return app, nil
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
	backend.Logger.Info("Called when the settings change", "cfg", a.settings)
}
