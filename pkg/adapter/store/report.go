package store

import (
	"context"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/db"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/log"

	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/dto"
)

type ReportStoreManager interface {
	Report(ctx context.Context, id int) (*dto.Result, error)
	Reports(ctx context.Context, query string) (*dto.Result, error)
	NewReport(ctx context.Context, report dto.Report) (*dto.Result, error)
	UpdateReport(ctx context.Context, id int, report dto.Report) (*dto.Result, error)
	DeleteReport(ctx context.Context, id int) (*dto.Result, error)
}

type ReportStore struct {
	db     db.DB
	logger log.Logger
}

func NewReportStore(db db.DB, logger log.Logger) ReportStore {
	return ReportStore{
		db:     db,
		logger: logger,
	}
}

func (s *ReportStore) Report(ctx context.Context, id int) (*dto.Result, error) {
	return nil, nil
}

func (s *ReportStore) Reports(ctx context.Context, query string) (*dto.Result, error) {
	return nil, nil
}

func (s *ReportStore) NewReport(ctx context.Context, report dto.Report) (*dto.Result, error) {
	return nil, nil
}

func (s *ReportStore) UpdateReport(ctx context.Context, id int, report dto.Report) (*dto.Result, error) {
	return nil, nil
}

func (s *ReportStore) DeleteReport(ctx context.Context, id int) (*dto.Result, error) {
	return nil, nil
}
