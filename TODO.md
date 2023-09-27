plugin/app.go > handler/cron/cron.go > handler/http/report_schedule (all schedules) 
                                     > domain/service/report_schedule (schedule job[s])

handler/http/report_schedule (new) > domain/service/report_schedule > store/boltdb/report_schedule
                                   > domain/service/report_schedule (schedule job)
