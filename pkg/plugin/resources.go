package plugin

import "net/http"

func (a *App) registerRoutes() {
	a.router.HandleFunc("/ping", a.handler.Ping)
	a.router.HandleFunc("/echo", a.handler.Echo)

	a.router.HandleFunc("/reports/{id:[0-9]+}", a.handler.Report).Methods(http.MethodGet)
	a.router.HandleFunc("/reports", a.handler.Reports).Methods(http.MethodGet)
	a.router.HandleFunc("/reports", a.handler.NewReport).Methods(http.MethodPost)
	a.router.HandleFunc("/reports/{id:[0-9]+}", a.handler.UpdateReport).Methods(http.MethodPatch, http.MethodPut)
	a.router.HandleFunc("/reports/{id:[0-9]+}", a.handler.DeleteReport).Methods(http.MethodDelete)

	a.router.HandleFunc("/report/schedules", a.handler.ReportSchedule).Methods(http.MethodGet)
}
