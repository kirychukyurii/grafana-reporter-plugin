//go:build wireinject
// +build wireinject

package plugin

import (
	"github.com/google/wire"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/adapter/grafana"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/adapter/store"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/config"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/domain/service"
	httphandler "github.com/kirychukyurii/grafana-reporter-plugin/pkg/handler/http"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/cdp"
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
)

func Initialize(*config.ReporterAppConfig, boltdb.DatabaseManager, *log.Logger, grafana.GrafanaHTTPAdapter, cdp.BrowserPoolManager) (*App, error) {
	wire.Build(wireBasicSet, newApp)
	return &App{}, nil
}
