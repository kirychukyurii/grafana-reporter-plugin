package handler

import (
	"encoding/json"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/model"
	"net/http"
)

func (h *handler) Report(w http.ResponseWriter, req *http.Request) {
	_, err := h.service.Report(req.Context(), req.Context().Value("id").(int))
	if err != nil {
		return
	}
}

func (h *handler) Reports(w http.ResponseWriter, req *http.Request) {}

func (h *handler) NewReport(w http.ResponseWriter, req *http.Request) {
	var report model.Report

	if err := json.NewDecoder(req.Body).Decode(&report); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.NewReport(req.Context(), report); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	return
}

func (h *handler) UpdateReport(w http.ResponseWriter, req *http.Request) {}

func (h *handler) DeleteReport(w http.ResponseWriter, req *http.Request) {}
