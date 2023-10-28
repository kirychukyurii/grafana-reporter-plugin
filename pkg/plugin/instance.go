package plugin

import (
	"context"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/cdp"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/cron"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/grafana"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/smtp"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/store"

	"github.com/gorilla/mux"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/app"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource/httpadapter"

	cronhandler "github.com/kirychukyurii/grafana-reporter-plugin/pkg/handler/cron"
	handler "github.com/kirychukyurii/grafana-reporter-plugin/pkg/handler/http"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/log"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/setting"
)

var _ instancemgmt.InstanceDisposer = (*AppInstance)(nil)

type AppInstance struct {
	OrgID               int
	httpResourceHandler backend.CallResourceHandler

	settings *setting.Setting
	logger   *log.Logger
	router   *mux.Router
	handler  handler.HandlerManager
	cron     cronhandler.ReportScheduleCronHandler
}

// New creates a new *App instance.
func New(logger *log.Logger, db store.DatabaseManager, is *setting.Setting, schedulers *cron.Schedulers, m smtp.Sender, b cdp.BrowserPoolManager, gcli *grafana.Client) app.InstanceFactoryFunc {
	return func(ctx context.Context, s backend.AppInstanceSettings, orgID int64, grafanaCfg backend.GrafanaCfg) (instancemgmt.Instance, error) {

		_ = setting.NewInstanceSetting(s)

		logger.Info("settings", "set", s)
		logger.Info("orgID", "orgID", orgID)
		logger.Info("grafanaCfg", "grafanaCfg", grafanaCfg)

		/*

			a, err := Initialize(setting, database, logger, grafanaClient, pool, scheduler, m)
			if err != nil {
				return nil, fmt.Errorf("initializing app: %v", err)
			}

			if err = a.cron.LoadSchedules(); err != nil {
				return nil, err
			}

			a.registerRoutes()

		*/

		return nil, nil
	}
}

func newAppInstance(handler handler.HandlerManager, cronHandler cronhandler.ReportScheduleCronHandler) (*AppInstance, error) {
	router := mux.NewRouter()
	httpResourceHandler := httpadapter.New(router)

	a := &AppInstance{
		httpResourceHandler: httpResourceHandler,

		router:  router,
		handler: handler,
		cron:    cronHandler,
	}

	return a, nil
}

// Dispose here tells plugin SDK that plugin wants to clean up resources when a new instance
// created.
func (a *AppInstance) Dispose() {
	a.logger.Info("called when the settings change", "cfg", a.settings)
}
