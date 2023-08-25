package store

type DatabaseAdapter interface {
	Reporter
	ReportScheduler
}

type databaseAdapter struct{}

func New() DatabaseAdapter {
	return &databaseAdapter{}
}
