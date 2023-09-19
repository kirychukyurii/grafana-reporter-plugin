package main

import (
	"context"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/adapter/store"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/config"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/domain/entity"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/domain/service"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/cdp"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/db"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/grafana"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/migration"
	"log"
)

func main() {
	s := config.ReporterAppSetting{}
	s.Browser.Url = "127.0.0.1"
	s.GrafanaBaseURL = "http://grafana:3000"
	s.BasicAuth.Username = "admin"
	s.BasicAuth.Password = "admin"
	s.WorkersCount = 10
	s.TemporaryDirectory = "/tmp/reporter/tmp"
	//s.Browser.Url = "chrome"

	database, err := db.New()
	if err != nil {
		log.Fatalln("msg", "db.New", "err", err)
	}

	if err = migration.Migrate(database); err != nil {
		log.Fatalln("msg", "migration.Migrate", "err", err)
	}

	browserPool := cdp.NewBrowserPool(2)
	gclient, _ := grafana.New(s)
	stor := store.New()
	svc := service.New(s, stor, gclient, browserPool)

	ctx := context.Background()
	report := entity.Report{
		Dashboard: entity.ReportDashboard{UID: "efee5237-d0c5-4fe4-897e-261e1aca6a1c"},
	}

	err = svc.NewReport(ctx, report)
	if err != nil {
		log.Fatalln("msg", "svc.NewReport", "err", err)
	}

	return
}
