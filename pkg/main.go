package main

import (
	"os"

	"github.com/grafana/grafana-plugin-sdk-go/backend/app"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/plugin"
)

func main() {
	if err := app.Manage("kirychukyurii-reporter-app", plugin.New, app.ManageOpts{}); err != nil {
		log.DefaultLogger.Error(err.Error())
		os.Exit(1)
	}
}
