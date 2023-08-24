package service

import (
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/adapter/grafana"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/models"
)

type Service interface {
	Reporter
	ReportScheduler
}

type service struct {
	settings          models.ReporterAppSetting
	grafanaHTTPClient grafana.GrafanaHTTPAdapter
}

func New(settings models.ReporterAppSetting, client grafana.GrafanaHTTPAdapter) Service {
	return &service{
		settings:          settings,
		grafanaHTTPClient: client,
	}
}
