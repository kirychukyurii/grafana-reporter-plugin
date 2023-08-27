package store

import (
	"context"

	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/models/dto"
)

type Reporter interface {
	Report(ctx context.Context, id int) (*dto.Result, error)
	Reports(ctx context.Context, query string) (*dto.Result, error)
	NewReport(ctx context.Context, report dto.Report) (*dto.Result, error)
	UpdateReport(ctx context.Context, id int, report dto.Report) (*dto.Result, error)
	DeleteReport(ctx context.Context, id int) (*dto.Result, error)
}

func (a *databaseAdapter) Report(ctx context.Context, id int) (*dto.Result, error) {
	return nil, nil
}

func (a *databaseAdapter) Reports(ctx context.Context, query string) (*dto.Result, error) {
	return nil, nil
}

func (a *databaseAdapter) NewReport(ctx context.Context, report dto.Report) (*dto.Result, error) {
	return nil, nil
}

func (a *databaseAdapter) UpdateReport(ctx context.Context, id int, report dto.Report) (*dto.Result, error) {
	return nil, nil
}

func (a *databaseAdapter) DeleteReport(ctx context.Context, id int) (*dto.Result, error) {
	return nil, nil
}
