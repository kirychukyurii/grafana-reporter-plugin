package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/domain/entity"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/domain/service"
)

type ReportSchedule struct {
	service service.ReportScheduleService
}

func NewReportScheduleHandler(service service.ReportScheduleService) *ReportSchedule {
	return &ReportSchedule{service: service}
}

func (r *ReportSchedule) ReportSchedule(w http.ResponseWriter, req *http.Request) {
	schedule, err := r.service.ReportSchedule(req.Context(), req.Context().Value("id").(int))
	if err != nil {
		Response{Error: err}.JSON(w)
		return
	}

	Response{Code: http.StatusOK, Data: schedule}.JSON(w)
}

func (r *ReportSchedule) ReportSchedules(w http.ResponseWriter, req *http.Request) {}

func (r *ReportSchedule) NewReportSchedule(w http.ResponseWriter, req *http.Request) {
	var schedule entity.ReportSchedule

	if err := json.NewDecoder(req.Body).Decode(&schedule); err != nil {
		Response{Code: http.StatusBadRequest, Error: fmt.Errorf("decode json: %v", err)}.JSON(w)
		return
	}

	if err := r.service.NewReportSchedule(req.Context(), schedule); err != nil {
		Response{Error: err}.JSON(w)
		return
	}

	Response{Code: http.StatusOK}.JSON(w)
}

func (r *ReportSchedule) UpdateReportSchedule(w http.ResponseWriter, req *http.Request) {}

func (r *ReportSchedule) DeleteReportSchedule(w http.ResponseWriter, req *http.Request) {}
