package http

import (
	"encoding/json"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/domain/entity"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/domain/service"
	"net/http"
)

type Report struct {
	service service.ReportService
}

func NewReportHandler(service service.ReportService) *Report {
	return &Report{service: service}
}

func (r *Report) Report(w http.ResponseWriter, req *http.Request) {
	_, err := r.service.Report(req.Context(), req.Context().Value("id").(int))
	if err != nil {
		return
	}
}

func (r *Report) Reports(w http.ResponseWriter, req *http.Request) {}

func (r *Report) NewReport(w http.ResponseWriter, req *http.Request) {
	var report entity.Report

	if err := json.NewDecoder(req.Body).Decode(&report); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := r.service.NewReport(req.Context(), report); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	return
}

func (r *Report) UpdateReport(w http.ResponseWriter, req *http.Request) {}

func (r *Report) DeleteReport(w http.ResponseWriter, req *http.Request) {}
