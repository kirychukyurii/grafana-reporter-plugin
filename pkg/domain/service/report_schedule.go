package service

import (
	"context"
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/util"
	"strconv"

	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/adapter/store"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/config"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/domain/entity"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/cron"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/log"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/smtp"
)

type ReportScheduleService interface {
	ReportSchedule(ctx context.Context, id int) (*entity.ReportSchedule, error)
	ReportSchedules(ctx context.Context, query string) ([]entity.ReportSchedule, error)
	NewReportSchedule(ctx context.Context, report entity.ReportSchedule) error
	NewReportScheduleJob(ctx context.Context, schedule entity.ReportSchedule) error
	UpdateReportSchedule(ctx context.Context, id int, report entity.ReportSchedule) error
	DeleteReportSchedule(ctx context.Context, id int) error
}

type ReportSchedule struct {
	settings *config.ReporterAppConfig
	logger   *log.Logger
	report   ReportService
	store    store.ReportScheduleStoreManager
	schedule cron.ScheduleManager
	sender   smtp.Sender
}

func NewReportScheduleService(settings *config.ReporterAppConfig, logger *log.Logger, report ReportService, store store.ReportScheduleStoreManager, schedule cron.ScheduleManager, sender smtp.Sender) *ReportSchedule {
	subLogger := &log.Logger{
		Logger: logger.With("component.type", "service", "component", "reportSchedule"),
	}

	return &ReportSchedule{
		settings: settings,
		logger:   subLogger,
		report:   report,
		store:    store,
		schedule: schedule,
		sender:   sender,
	}
}

func (r *ReportSchedule) ReportSchedule(ctx context.Context, id int) (*entity.ReportSchedule, error) {
	schedule, err := r.store.ReportSchedule(ctx, id)
	if err != nil {
		return nil, err
	}

	return schedule, nil
}

func (r *ReportSchedule) ReportSchedules(ctx context.Context, query string) ([]entity.ReportSchedule, error) {
	schedule, err := r.store.ReportSchedules(ctx, query)
	if err != nil {
		return nil, err
	}

	return schedule, nil
}

func (r *ReportSchedule) NewReportSchedule(ctx context.Context, report entity.ReportSchedule) error {
	_, err := r.store.NewReportSchedule(ctx, report)
	if err != nil {
		return err
	}

	if err = r.NewReportScheduleJob(ctx, report); err != nil {
		return err
	}

	return nil
}

func (r *ReportSchedule) NewReportScheduleJob(ctx context.Context, schedule entity.ReportSchedule) error {
	var err error

	fn := func(job gocron.Job) error {
		defer func() {
			subLogger := r.logger.With("schedule.id", schedule.ID, "job.next_run", job.NextRun(), "job.count", job.RunCount())
			if err != nil {
				subLogger = subLogger.With("error", err.Error())
			}

			subLogger.Debug("completed scheduled job")
		}()

		tmpDir, err := r.report.NewReport(ctx, schedule.Report)
		if err != nil {
			return err
		}

		attachments, err := util.ReadDir(tmpDir)
		if err != nil {
			return err
		}

		if err := r.sender.Send(schedule.Report.Recipients,
			[]byte(fmt.Sprintf("%s: %s", schedule.Name, schedule.Report.Name)),
			[]byte(schedule.Report.Message), attachments); err != nil {
			return err
		}

		return nil
	}

	job, err := r.schedule.ScheduleJob(schedule.Interval, strconv.FormatInt(schedule.ID, 10), fn)
	if err != nil {
		return err
	}

	r.logger.Debug("schedule job", "schedule.id", schedule.ID, "job.next_run", job.NextRun())

	return nil
}

func (r *ReportSchedule) UpdateReportSchedule(ctx context.Context, id int, report entity.ReportSchedule) error {
	_, err := r.store.UpdateReportSchedule(ctx, id, report)
	if err != nil {
		return err
	}

	return nil
}

func (r *ReportSchedule) DeleteReportSchedule(ctx context.Context, id int) error {
	_, err := r.store.DeleteReportSchedule(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *ReportSchedule) RunReportSchedule(ctx context.Context, id int) error {
	return nil
}
