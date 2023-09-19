package service

type ReportScheduleService interface{}

type ReportSchedule struct{}

func NewReportScheduleService() *ReportSchedule {
	return &ReportSchedule{}
}
