package service

import "github.com/kirychukyurii/grafana-reporter-plugin/pkg/service/report"

type Service struct {
	Reporter report.Service
}

func New() Service {
	reportService := report.New()

	return Service{
		Reporter: reportService,
	}
}
