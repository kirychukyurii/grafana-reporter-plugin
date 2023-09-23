package service

import (
	"context"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/adapter/store"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/config"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/domain/entity"
)

type ReportScheduleService interface {
	ReportSchedule(ctx context.Context, id int) (*entity.ReportSchedule, error)
	ReportSchedules(ctx context.Context, query string) (*entity.Result, error)
	NewReportSchedule(ctx context.Context, report entity.ReportSchedule) error
	UpdateReportSchedule(ctx context.Context, id int, report entity.ReportSchedule) error
	DeleteReportSchedule(ctx context.Context, id int) error
}

type ReportSchedule struct {
	settings *config.ReporterAppConfig
	store    store.ReportScheduleStoreManager
}

func NewReportScheduleService(settings *config.ReporterAppConfig,
	store store.ReportScheduleStoreManager,
) *ReportSchedule {
	return &ReportSchedule{
		settings: settings,
		store:    store,
	}
}

func (r *ReportSchedule) ReportSchedule(ctx context.Context, id int) (*entity.ReportSchedule, error) {
	schedule, err := r.store.ReportSchedule(ctx, id)
	if err != nil {
		return nil, err
	}

	return schedule, nil
}

func (r *ReportSchedule) ReportSchedules(ctx context.Context, query string) (*entity.Result, error) {
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
