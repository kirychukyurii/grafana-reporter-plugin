package plugin

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource/httpadapter"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/handler"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/models"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/service"
	"net/http"
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
	var app App

	if err := app.settings.Load(s); err != nil {
		return nil, err
	}

	if err := app.settings.Validate(); err != nil {
		return nil, err
	}

	svc := service.New()
	app.handler = handler.New(svc)

	app.router = mux.NewRouter()
	app.registerRoutes()

	return &app, nil
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}

// CallResource HTTP style resource
func (a *App) CallResource(ctx context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error {
	return httpadapter.New(a).CallResource(ctx, req, sender)
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
	// cleanup
}
