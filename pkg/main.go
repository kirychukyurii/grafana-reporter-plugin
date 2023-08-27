package main

import (
	"os"

	"github.com/grafana/grafana-plugin-sdk-go/backend/app"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/plugin"
)

const pluginID = "kirychukyurii-reporter-app"

func main() {
	log.DefaultLogger.Info("serving the app over gPRC with automatic instance management")
	if err := app.Manage(pluginID, plugin.New, app.ManageOpts{}); err != nil {
		log.DefaultLogger.Error(err.Error())
		os.Exit(1)
	}
}
