//go:build wireinject
// +build wireinject

package plugin

import (
	"github.com/google/wire"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/smtp"

	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/adapter/grafana"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/adapter/store"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/config"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/domain/service"
	cronhandler "github.com/kirychukyurii/grafana-reporter-plugin/pkg/handler/cron"
	httphandler "github.com/kirychukyurii/grafana-reporter-plugin/pkg/handler/http"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/cdp"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/cron"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/log"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/store/boltdb"
)

var wireBasicSet = wire.NewSet(
	store.ProviderSet,
	wire.Bind(new(store.ReportScheduleStoreManager), new(*store.ReportScheduleStore)),
	service.ProviderSet,
	wire.Bind(new(service.ReportService), new(*service.Report)),
	wire.Bind(new(service.ReportScheduleService), new(*service.ReportSchedule)),
	httphandler.ProviderSet,
	wire.Bind(new(httphandler.ReportHandler), new(*httphandler.Report)),
	wire.Bind(new(httphandler.ReportScheduleHandler), new(*httphandler.ReportSchedule)),
	wire.Bind(new(httphandler.HandlerManager), new(*httphandler.Handler)),
	cronhandler.ProviderSet,
	wire.Bind(new(cronhandler.ReportScheduleCronHandler), new(*cronhandler.ReportScheduleCron)),
)

func Initialize(*config.ReporterAppConfig, boltdb.DatabaseManager, *log.Logger, grafana.GrafanaHTTPAdapter, cdp.BrowserPoolManager, cron.ScheduleManager, smtp.Sender) (*App, error) {
	wire.Build(wireBasicSet, newApp)
	return &App{}, nil
}
