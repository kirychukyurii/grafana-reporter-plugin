package store

type Adapter interface {
	ReportStore
	ReportScheduleStore
}

type adapter struct{}

func New() Adapter {
	return &adapter{}
}
