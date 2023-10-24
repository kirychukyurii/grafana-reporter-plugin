package plugin

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/app"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource/httpadapter"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/config"
	cronhandler "github.com/kirychukyurii/grafana-reporter-plugin/pkg/handler/cron"
	handler "github.com/kirychukyurii/grafana-reporter-plugin/pkg/handler/http"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/cdp"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/cron"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/grafana"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/log"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/smtp"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/store"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/store/boltdb"
	"time"
)

var _ instancemgmt.InstanceDisposer = (*AppInstance)(nil)

type AppInstance struct {
	OrgID               int
	httpResourceHandler backend.CallResourceHandler

	settings *config.ReporterAppConfig
	logger   *log.Logger
	router   *mux.Router
	handler  handler.HandlerManager
	cron     cronhandler.ReportScheduleCronHandler
}

// New creates a new *App instance.
// func New(ctx context.Context, s backend.AppInstanceSettings) (instancemgmt.Instance, error)
func New(logger *log.Logger) app.InstanceFactoryFunc {
	return func(ctx context.Context, s backend.AppInstanceSettings) (instancemgmt.Instance, error) {
		pluginContext := httpadapter.PluginConfigFromContext(ctx)
		logger.Info("New", "pluginContext", pluginContext)
		setting, err := config.New(s)
		if err != nil {
			return nil, fmt.Errorf("initializing config: %v", err)
		}

		logger.Info("settings", "set", s)

		database, err := store.New(logger, store.Opts{
			Type: store.BoltDB,
			BoltDBOpts: &boltdb.Opts{
				DataDirectory:   setting.DataDirectory,
				EncryptionKey:   setting.DatabaseConfig.EncryptionKey,
				Timeout:         setting.DatabaseConfig.Timeout,
				InitialMmapSize: setting.DatabaseConfig.InitialMmapSize,
				MaxBatchSize:    setting.DatabaseConfig.MaxBatchSize,
				MaxBatchDelay:   setting.DatabaseConfig.MaxBatchDelay,
			},
		})
		if err != nil {
			return nil, fmt.Errorf("initializing database client: %v", err)
		}

		pool := cdp.NewBrowserPool(setting)
		grafanaClient, err := grafana.New(setting)
		if err != nil {
			return nil, fmt.Errorf("initializing grafana client: %v", err)
		}

		m, err := smtp.New(setting.MailConfig.Host, setting.MailConfig.Port, setting.MailConfig.Username, setting.MailConfig.Password)
		if err != nil {
			return nil, fmt.Errorf("initializing mail client: %v", err)
		}

		timezone, err := time.LoadLocation("Local")
		if err != nil {
			return nil, fmt.Errorf("load default timezone: %v", err)
		}

		scheduler := cron.NewScheduler(timezone)

		a, err := Initialize(setting, database, logger, grafanaClient, pool, scheduler, m)
		if err != nil {
			return nil, fmt.Errorf("initializing app: %v", err)
		}

		if err = a.cron.LoadSchedules(); err != nil {
			return nil, err
		}

		a.registerRoutes()

		return a, nil
	}
}

func newAppInstance(setting *config.ReporterAppConfig, handler handler.HandlerManager, cronHandler cronhandler.ReportScheduleCronHandler) (*AppInstance, error) {
	router := mux.NewRouter()
	httpResourceHandler := httpadapter.New(router)

	a := &AppInstance{
		httpResourceHandler: httpResourceHandler,

		settings: setting,
		router:   router,
		handler:  handler,
		cron:     cronHandler,
	}

	return a, nil
}

// Dispose here tells plugin SDK that plugin wants to clean up resources when a new instance
// created.
func (a *AppInstance) Dispose() {
	a.logger.Info("called when the settings change", "cfg", a.settings)
}
