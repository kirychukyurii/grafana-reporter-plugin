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
	db "github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/store/sqlite"
)

var wireBasicSet = wire.NewSet(
	store.ProviderSet,
	wire.Bind(new(store.ReportStoreManager), new(*store.ReportStore)),
	service.ProviderSet,
	wire.Bind(new(service.ReportService), new(*service.Report)),
	httphandler.ProviderSet,
	wire.Bind(new(httphandler.ReportHandler), new(*httphandler.Report)),
)

func Initialize(*config.ReporterAppSetting, *db.DB, *log.Logger, grafana.GrafanaHTTPAdapter, cdp.BrowserPoolManager) (*App, error) {
	wire.Build(wireBasicSet, newApp)
	return &App{}, nil
}
