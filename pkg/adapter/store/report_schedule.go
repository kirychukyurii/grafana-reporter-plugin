package store

import (
	"context"

	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/domain/entity"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/log"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/store/boltdb"
)

const reportScheduleBucketName = "report_schedule"

type ReportScheduleStoreManager interface {
	ReportSchedule(ctx context.Context, id int) (*entity.ReportSchedule, error)
	ReportSchedules(ctx context.Context, query string) (*entity.Result, error)
	NewReportSchedule(ctx context.Context, schedule entity.ReportSchedule) (*entity.ReportSchedule, error)
	UpdateReportSchedule(ctx context.Context, id int, schedule entity.ReportSchedule) (*entity.ReportSchedule, error)
	DeleteReportSchedule(ctx context.Context, id int) (*entity.ReportSchedule, error)
}

type ReportScheduleStore struct {
	db     boltdb.DatabaseManager
	logger *log.Logger
}

func NewReportScheduleStore(db boltdb.DatabaseManager, logger *log.Logger) *ReportScheduleStore {
	if err := db.SetServiceName(reportScheduleBucketName); err != nil {
		return nil
	}

	return &ReportScheduleStore{
		db:     db,
		logger: logger,
	}
}

func (s *ReportScheduleStore) ReportSchedule(ctx context.Context, id int) (*entity.ReportSchedule, error) {
	var schedule *entity.ReportSchedule

	if err := s.db.GetObject(reportScheduleBucketName, s.db.ConvertToKey(id), schedule); err != nil {
		return nil, err
	}

	return schedule, nil
}

func (s *ReportScheduleStore) ReportSchedules(ctx context.Context, query string) (*entity.Result, error) {
	var (
		schedule  entity.ReportSchedule
		schedules = make([]entity.ReportSchedule, 0)
	)

	if err := s.db.GetAll(reportScheduleBucketName, &schedule, boltdb.AppendFn(&schedules)); err != nil {
		return nil, err
	}

	return &entity.Result{Data: schedules}, nil
}

func (s *ReportScheduleStore) NewReportSchedule(ctx context.Context, schedule entity.ReportSchedule) (*entity.ReportSchedule, error) {
	objFn := func(id uint64) (int, interface{}) {
		return int(id), schedule
	}

	if err := s.db.CreateObject(reportScheduleBucketName, objFn); err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *ReportScheduleStore) UpdateReportSchedule(ctx context.Context, id int, schedule entity.ReportSchedule) (*entity.ReportSchedule, error) {
	if err := s.db.UpdateObject(reportScheduleBucketName, s.db.ConvertToKey(id), schedule); err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *ReportScheduleStore) DeleteReportSchedule(ctx context.Context, id int) (*entity.ReportSchedule, error) {
	if err := s.db.DeleteObject(reportScheduleBucketName, s.db.ConvertToKey(id)); err != nil {
		return nil, err
	}

	return nil, nil
}
