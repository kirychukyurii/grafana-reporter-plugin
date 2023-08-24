package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/service"
)

type Handler interface {
	Ping(w http.ResponseWriter, req *http.Request)
	Echo(w http.ResponseWriter, req *http.Request)

	Report(w http.ResponseWriter, req *http.Request)
	Reports(w http.ResponseWriter, req *http.Request)
	NewReport(w http.ResponseWriter, req *http.Request)
	UpdateReport(w http.ResponseWriter, req *http.Request)
	DeleteReport(w http.ResponseWriter, req *http.Request)

	ReportSchedule(w http.ResponseWriter, req *http.Request)
	NewReportSchedule(w http.ResponseWriter, req *http.Request)
	UpdateReportSchedule(w http.ResponseWriter, req *http.Request)
	DeleteReportSchedule(w http.ResponseWriter, req *http.Request)
}

type handler struct {
	service service.Service
}

func New(service service.Service) Handler {
	return &handler{
		service: service,
	}
}

// Ping is an example HTTP GET resource that returns a {"message": "ok"} JSON response.
func (h *handler) Ping(w http.ResponseWriter, req *http.Request) {
	writeJsonResponse(w, "ok", nil)
}

// Echo is an example HTTP POST resource that accepts a JSON with a "message" key and
// returns to the client whatever it is sent.
func (h *handler) Echo(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var body struct {
		Message string `json:"message"`
	}

	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJsonResponse(w, body, nil)
}

func writeJsonResponse(w http.ResponseWriter, rsp interface{}, err error) {
	w.Header().Add("Content-Type", "application/json")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"message": "%s"}`, err.Error())))
	} else {
		_ = json.NewEncoder(w).Encode(rsp)
	}
}
