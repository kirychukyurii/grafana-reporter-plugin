// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package plugin

import (
	"github.com/google/wire"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/adapter/grafana"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/adapter/store"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/config"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/domain/service"
	cron2 "github.com/kirychukyurii/grafana-reporter-plugin/pkg/handler/cron"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/handler/http"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/cdp"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/cron"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/log"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/smtp"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/store/boltdb"
)

// Injectors from wire.go:

func Initialize(reporterAppConfig *config.ReporterAppConfig, databaseManager boltdb.DatabaseManager, logger *log.Logger, grafanaHTTPAdapter grafana.GrafanaHTTPAdapter, browserPoolManager cdp.BrowserPoolManager, scheduleManager cron.ScheduleManager, sender smtp.Sender) (*App, error) {
	report := service.NewReportService(reporterAppConfig, logger, grafanaHTTPAdapter, browserPoolManager)
	httpReport := http.NewReportHandler(report)
	reportScheduleStore := store.NewReportScheduleStore(databaseManager, logger)
	reportSchedule := service.NewReportScheduleService(reporterAppConfig, logger, report, reportScheduleStore, scheduleManager, sender)
	httpReportSchedule := http.NewReportScheduleHandler(reportSchedule)
	handler := http.New(httpReport, httpReportSchedule)
	reportScheduleCron := cron2.NewReportScheduleCronHandler(logger, reportSchedule)
	app, err := newApp(reporterAppConfig, handler, reportScheduleCron)
	if err != nil {
		return nil, err
	}
	return app, nil
}

// wire.go:

var wireBasicSet = wire.NewSet(store.ProviderSet, wire.Bind(new(store.ReportScheduleStoreManager), new(*store.ReportScheduleStore)), service.ProviderSet, wire.Bind(new(service.ReportService), new(*service.Report)), wire.Bind(new(service.ReportScheduleService), new(*service.ReportSchedule)), http.ProviderSet, wire.Bind(new(http.ReportHandler), new(*http.Report)), wire.Bind(new(http.ReportScheduleHandler), new(*http.ReportSchedule)), wire.Bind(new(http.HandlerManager), new(*http.Handler)), cron2.ProviderSet, wire.Bind(new(cron2.ReportScheduleCronHandler), new(*cron2.ReportScheduleCron)))
