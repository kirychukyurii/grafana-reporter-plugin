package service

type ServiceManager interface {
	*ReportService
	*ReportScheduleService
}

type Service struct {
	report         Report
	reportSchedule ReportSchedule
}

func New(report ReportService, reportSchedule ReportScheduleService) *Service {
	return &Service{}
}
