package cron

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
)

type ScheduleManager interface {
	ScheduleJob(interval string, jobTag string, jobFun interface{}, params ...interface{}) (*gocron.Job, error)
	RemoveJob(jobTag string) error
}

type Scheduler struct {
	Cron *gocron.Scheduler
}

func NewScheduler(location *time.Location) *Scheduler {
	scheduler := gocron.NewScheduler(location)
	scheduler.StartAsync()

	return &Scheduler{
		Cron: scheduler,
	}
}

func (s *Scheduler) ScheduleJob(interval string, jobTag string, jobFun interface{}, params ...interface{}) (*gocron.Job, error) {
	job, err := s.Cron.Cron(interval).Tag(jobTag).DoWithJobDetails(jobFun, params...)
	if err != nil {
		return nil, fmt.Errorf("scheduling job: %v", err)
	}

	return job, nil
}

func (s *Scheduler) RemoveJob(jobTag string) error {
	if err := s.Cron.RemoveByTag(jobTag); err != nil {
		return fmt.Errorf("removing job by tag: %v", err)
	}

	return nil
}
