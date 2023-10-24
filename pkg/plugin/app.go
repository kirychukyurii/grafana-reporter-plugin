package plugin

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/app"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	cronhandler "github.com/kirychukyurii/grafana-reporter-plugin/pkg/handler/cron"
	handler "github.com/kirychukyurii/grafana-reporter-plugin/pkg/handler/http"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/log"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/setting"
)

// Make sure App implements required interfaces.
// This is important to do since otherwise we will only get a
// not implemented error response from plugin in runtime.
var (
	_ backend.CallResourceHandler = (*App)(nil)
	_ backend.CheckHealthHandler  = (*App)(nil)
)

type App struct {
	im instancemgmt.InstanceManager

	logger     *log.Logger
	appSetting *setting.AppSetting

	router  *mux.Router
	handler handler.HandlerManager
	cron    cronhandler.ReportScheduleCronHandler
}

func NewApp(logger *log.Logger) (*App, error) {
	logger.Info("NewApp")
	im := app.NewInstanceManager(New(logger))
	as, err := setting.NewAppSetting()
	if err != nil {
		return nil, err
	}

	/*database, err := boltdb.New(setting, logger)
	if err != nil {
		return nil, fmt.Errorf("initializing database client: %v", err)
	}

	pool := cdp.NewBrowserPool(setting)
	grafanaClient, err := grafana.New(setting)
	if err != nil {
		return nil, fmt.Errorf("initializing grafana client: %v", err)
	}

	mail, err := smtp.New(as.MailConfig.Host, as.MailConfig.Port, as.MailConfig.Username, as.MailConfig.Password)
	if err != nil {
		return nil, fmt.Errorf("initializing mail client: %v", err)
	}

	timezone, err := time.LoadLocation("Local")
	if err != nil {
		return nil, fmt.Errorf("load default timezone: %v", err)
	}

	scheduler := cron.NewScheduler(timezone)*/

	return &App{
		im: im,

		logger:     logger,
		appSetting: as,
	}, nil
}

// CallResource HTTP style resource
func (a *App) CallResource(ctx context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error {
	a.logger.Info("CallResource", "pluginContext", req.PluginContext)

	appInstance, err := a.appInstance(ctx, req.PluginContext)
	if err != nil {
		a.logger.Error(err.Error())
	}

	a.logger.Info("CallResource", "OrgID", appInstance.OrgID)

	appInstance.logger = a.logger

	return appInstance.httpResourceHandler.CallResource(ctx, req, sender)
}

// CheckHealth handles health checks sent from Grafana to the plugin.
func (a *App) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	_, err := a.appInstance(ctx, req.PluginContext)
	if err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: err.Error(),
		}, err
	}

	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: "ok",
	}, nil
}

// appInstance Returns cached datasource or creates new one
func (a *App) appInstance(ctx context.Context, pluginContext backend.PluginContext) (*AppInstance, error) {
	instance, err := a.im.Get(ctx, pluginContext)
	if err != nil {
		return nil, err
	}

	appInstance, ok := instance.(*AppInstance)
	if !ok {
		return nil, fmt.Errorf("cannot use instancemgmt.Instance as the type *plugin.AppInstance")
	}

	return appInstance, nil
}
