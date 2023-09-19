package http

import (
	"encoding/json"
	"fmt"
	"github.com/google/wire"
	"net/http"
)

// ProviderSet is handler providers.
var ProviderSet = wire.NewSet(NewReportHandler, NewReportScheduleHandler)

type HandlerManager interface {
	Ping(w http.ResponseWriter, req *http.Request)
	Echo(w http.ResponseWriter, req *http.Request)

	ReportHandler
	ReportScheduleHandler
}

type ReportHandler interface {
	Report(w http.ResponseWriter, req *http.Request)
	Reports(w http.ResponseWriter, req *http.Request)
	NewReport(w http.ResponseWriter, req *http.Request)
	UpdateReport(w http.ResponseWriter, req *http.Request)
	DeleteReport(w http.ResponseWriter, req *http.Request)
}

type ReportScheduleHandler interface {
	ReportSchedule(w http.ResponseWriter, req *http.Request)
	NewReportSchedule(w http.ResponseWriter, req *http.Request)
	UpdateReportSchedule(w http.ResponseWriter, req *http.Request)
	DeleteReportSchedule(w http.ResponseWriter, req *http.Request)
}

type Handler struct {
	*Report
	*ReportSchedule
}

func New(reportHandler ReportHandler, reportScheduleHandler ReportScheduleHandler) *Handler {
	return &Handler{
		Report:         reportHandler.(*Report),
		ReportSchedule: reportScheduleHandler.(*ReportSchedule),
	}
}

// Ping is an example HTTP GET resource that returns a {"message": "ok"} JSON response.
func (h *Handler) Ping(w http.ResponseWriter, req *http.Request) {
	writeJsonResponse(w, "ok", nil)
}

// Echo is an example HTTP POST resource that accepts a JSON with a "message" key and
// returns to the client whatever it is sent.
func (h *Handler) Echo(w http.ResponseWriter, req *http.Request) {
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
