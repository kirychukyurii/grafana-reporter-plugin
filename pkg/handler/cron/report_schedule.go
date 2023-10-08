package cron

import (
	"context"

	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/domain/service"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/log"
)

type ReportScheduleCronHandler interface {
	LoadSchedules() error
}

type ReportScheduleCron struct {
	logger  *log.Logger
	service service.ReportScheduleService
}

func NewReportScheduleCronHandler(logger *log.Logger, service service.ReportScheduleService) *ReportScheduleCron {
	return &ReportScheduleCron{logger: logger, service: service}
}

func (r *ReportScheduleCron) LoadSchedules() error {
	ctx := context.Background()

	schedules, err := r.service.ReportSchedules(ctx, "")
	if err != nil {
		return err
	}

	for _, schedule := range schedules {
		r.logger.Debug("schedule job", "job", schedule)
		if err = r.service.NewReportScheduleJob(ctx, schedule); err != nil {
			return err
		}
	}

	return nil
}
