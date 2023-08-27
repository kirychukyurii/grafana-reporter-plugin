package cron

import (
	"context"
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
)

type Scheduler struct {
	Cron *gocron.Scheduler
}

func NewScheduler(location *time.Location) Scheduler {
	scheduler := gocron.NewScheduler(location)
	// scheduler.SetMaxConcurrentJobs()

	return Scheduler{
		Cron: scheduler,
	}
}

func (s *Scheduler) ScheduleJob(jobTag string, jobFun interface{}, ctx context.Context, interval string) (*gocron.Job, error) {
	job, err := s.Cron.Cron(interval).Tag(jobTag).DoWithJobDetails(jobFun, ctx)
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
