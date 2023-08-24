package handler

import (
	"net/http"
)

func (h *handler) Report(w http.ResponseWriter, req *http.Request) {
	_, err := h.service.Report(req.Context(), req.Context().Value("id").(int))
	if err != nil {
		return
	}
}

func (h *handler) Reports(w http.ResponseWriter, req *http.Request) {}

func (h *handler) NewReport(w http.ResponseWriter, req *http.Request) {}

func (h *handler) UpdateReport(w http.ResponseWriter, req *http.Request) {}

func (h *handler) DeleteReport(w http.ResponseWriter, req *http.Request) {}
