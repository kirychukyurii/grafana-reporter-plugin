package dto

import "time"

type ReportSchedule struct {
	ID        int64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
	Name      string
	Active    bool
	Report    string
	StartDate time.Time
	EndDate   time.Time
	Timezone  string
	Interval  string
	WorkDays  bool
}
