//go:build wireinject
// +build wireinject

package plugin

import (
	"github.com/google/wire"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/setting"

	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/adapter/grafana"
	storeadapter "github.com/kirychukyurii/grafana-reporter-plugin/pkg/adapter/store"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/domain/service"
	cronhandler "github.com/kirychukyurii/grafana-reporter-plugin/pkg/handler/cron"
	httphandler "github.com/kirychukyurii/grafana-reporter-plugin/pkg/handler/http"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/cdp"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/cron"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/log"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/smtp"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/store"
)

var wireBasicSet = wire.NewSet(
	storeadapter.ProviderSet,
	wire.Bind(new(storeadapter.ReportScheduleStoreManager), new(*storeadapter.ReportScheduleStore)),
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

func Initialize(*setting.Setting, store.DatabaseManager, *log.Logger, grafana.DashboardAdapter, cdp.BrowserPoolManager, *cron.Schedulers, smtp.Sender) (*AppInstance, error) {
	wire.Build(wireBasicSet, newAppInstance)

	return &AppInstance{}, nil
}

func InitializeCronHandler(*setting.Setting, store.DatabaseManager, *log.Logger, grafana.DashboardAdapter, cdp.BrowserPoolManager, *cron.Schedulers, smtp.Sender) (*cronhandler.ReportScheduleCron, error) {
	wire.Build(wireBasicSet)

	return &cronhandler.ReportScheduleCron{}, nil
}
