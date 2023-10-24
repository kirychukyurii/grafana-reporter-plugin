package main

import (
	"github.com/grafana/grafana-plugin-sdk-go/backend"

	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/log"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/plugin"
)

const pluginID = "kirychukyurii-reporter-app"

func main() {
	backend.SetupPluginEnvironment(pluginID)
	logger := log.New()
	app, err := plugin.NewApp(logger)
	if err != nil {
		logger.Fatal(err.Error())
	}

	opts := backend.ServeOpts{
		CheckHealthHandler:  app,
		CallResourceHandler: app,
	}

	logger.Info("serve app")

	if err := backend.Serve(opts); err != nil {
		logger.Fatal(err.Error())
	}
}
