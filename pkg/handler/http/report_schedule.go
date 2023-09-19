package http

import "net/http"

type ReportSchedule struct{}

func NewReportScheduleHandler() *ReportSchedule {
	return &ReportSchedule{}
}

func (r *ReportSchedule) ReportSchedule(w http.ResponseWriter, req *http.Request) {}

func (r *ReportSchedule) NewReportSchedule(w http.ResponseWriter, req *http.Request) {}

func (r *ReportSchedule) UpdateReportSchedule(w http.ResponseWriter, req *http.Request) {}

func (r *ReportSchedule) DeleteReportSchedule(w http.ResponseWriter, req *http.Request) {}
