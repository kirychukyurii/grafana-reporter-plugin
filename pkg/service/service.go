package service

import (
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/adapter/grafana"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/adapter/store"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/browser"
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
	browserPool       browser.Browser
}

func New(settings model.ReporterAppSetting, database store.DatabaseAdapter, client grafana.GrafanaHTTPAdapter) Service {
	return &service{
		settings: settings,

		database:          database,
		grafanaHTTPClient: client,
	}
}
