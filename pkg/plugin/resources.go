package plugin

import "net/http"

func (a *AppInstance) registerRoutes() {
	//a.router.HandleFunc("/ping", a.handler.Ping)
	//a.router.HandleFunc("/echo", a.handler.Echo)

	a.router.HandleFunc("/reports/{id:[0-9]+}", a.handler.Report).Methods(http.MethodGet)
	a.router.HandleFunc("/reports", a.handler.Reports).Methods(http.MethodGet)
	a.router.HandleFunc("/reports", a.handler.NewReport).Methods(http.MethodPost)
	a.router.HandleFunc("/reports/{id:[0-9]+}", a.handler.UpdateReport).Methods(http.MethodPatch, http.MethodPut)
	a.router.HandleFunc("/reports/{id:[0-9]+}", a.handler.DeleteReport).Methods(http.MethodDelete)

	a.router.HandleFunc("/report/schedules/{id:[0-9]+}", a.handler.ReportSchedule).Methods(http.MethodGet)
	a.router.HandleFunc("/report/schedules", a.handler.ReportSchedules).Methods(http.MethodGet)
	a.router.HandleFunc("/report/schedules", a.handler.NewReportSchedule).Methods(http.MethodPost)
	a.router.HandleFunc("/report/schedules/{id:[0-9]+}", a.handler.UpdateReportSchedule).Methods(http.MethodPatch, http.MethodPut)
	a.router.HandleFunc("/report/schedules/{id:[0-9]+}", a.handler.DeleteReportSchedule).Methods(http.MethodDelete)
}
