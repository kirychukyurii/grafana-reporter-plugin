package dto

import "time"

type ReportSchedule struct {
	ID                int
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         time.Time
	Name              string
	Active            bool
	Report            Report
	StartDate         time.Time
	EndDate           time.Time
	Timezone          string
	Frequency         string
	IntervalFrequency string
	IntervalAmount    string
	WorkDays          bool
	DayOfMonth        int
}
