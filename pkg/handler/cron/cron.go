package cron

import (
	"context"
	"fmt"
	"github.com/google/wire"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/log"

	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/domain/entity"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/domain/service"
)

var ProviderSet = wire.NewSet(NewReportScheduleCronHandler)

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

	s, ok := schedules.Data.([]entity.ReportSchedule)
	if !ok {
		return fmt.Errorf("cant cast to provided type")
	}

	for _, schedule := range s {
		r.logger.Debug("schedule job", "job", schedule)
		if err = r.service.NewReportScheduleJob(ctx, schedule); err != nil {
			return err
		}
	}

	return nil
}
