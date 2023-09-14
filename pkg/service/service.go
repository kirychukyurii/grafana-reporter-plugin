package service

import (
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/adapter/grafana"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/adapter/store"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/cdp"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/model"
)

type Service interface {
	Reporter
	ReportScheduler
}

type service struct {
	settings model.ReporterAppSetting

	database          store.DatabaseAdapter
	grafanaHTTPClient grafana.GrafanaHTTPAdapter
	browserPool       cdp.BrowserPoolManager
}

func New(settings model.ReporterAppSetting, database store.DatabaseAdapter, client grafana.GrafanaHTTPAdapter, browserPool cdp.BrowserPoolManager) Service {
	return &service{
		settings: settings,

		database:          database,
		grafanaHTTPClient: client,
		browserPool:       browserPool,
	}
}
