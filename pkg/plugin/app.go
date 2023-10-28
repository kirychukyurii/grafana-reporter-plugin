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
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/cdp"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/cron"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/grafana"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/log"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/smtp"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/store"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/store/boltdb"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/setting"
	"net/url"
	"time"
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

	as, err := setting.NewAppSetting()
	if err != nil {
		return nil, err
	}

	db, err := store.New(logger, &store.Options{
		Type: store.BoltDB,
		BoltDBOpts: &boltdb.Options{
			DataDirectory:   as.DataDirectory,
			EncryptionKey:   as.DatabaseConfig.EncryptionKey,
			Timeout:         as.DatabaseConfig.Timeout,
			InitialMmapSize: as.DatabaseConfig.InitialMmapSize,
			MaxBatchSize:    as.DatabaseConfig.MaxBatchSize,
			MaxBatchDelay:   as.DatabaseConfig.MaxBatchDelay,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("initializing database client: %v", err)
	}

	is, err := setting.InstanceSettingFromStore(db)
	if err != nil {
		return nil, err
	}

	s := &setting.Setting{
		AppSetting:       as,
		InstancesSetting: is,
	}

	var schedulers cron.Schedulers
	for _, os := range is {
		timezone, err := time.LoadLocation(os.Timezone)
		if err != nil {
			return nil, fmt.Errorf("load default timezone: %v", err)
		}

		schedulers[os.OrgID] = cron.NewScheduler(timezone)
	}

	b := cdp.NewBrowserPool(as.WorkersCount)
	gURL, err := url.Parse(as.GrafanaSetting.URL)
	if err != nil {
		return nil, err
	}

	gCli, err := grafana.New(&grafana.Options{
		URL:                  gURL,
		InsecureSkipVerify:   as.GrafanaSetting.InsecureSkipVerify,
		RetryNum:             as.GrafanaSetting.RetryNum,
		RetryTimeout:         as.GrafanaSetting.RetryTimeout,
		RetryStatusCodes:     as.GrafanaSetting.RetryStatusCodesArr(),
		APIToken:             as.GrafanaSetting.APIToken,
		BasicAuthCredentials: as.GrafanaSetting.BasicAuth(),
	})
	if err != nil {
		return nil, fmt.Errorf("initializing grafana client: %v", err)
	}

	m, err := smtp.New(&smtp.Options{
		Host:     as.MailConfig.Host,
		Port:     as.MailConfig.Port,
		Username: as.MailConfig.Username,
		Password: as.MailConfig.Password,
		From:     "",
	})
	if err != nil {
		return nil, fmt.Errorf("initializing mail client: %v", err)
	}

	cronHandler, err := InitializeCronHandler(s, db, logger, gCli, b, schedulers, m)
	if err != nil {
		return nil, err
	}

	if err := cronHandler.LoadSchedules(); err != nil {
		return nil, err
	}

	im := app.NewInstanceManager(New(logger, db, s, schedulers, m, b, gCli))

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
