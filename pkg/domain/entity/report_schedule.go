package entity

import "time"

type ReportSchedule struct {
	ID        int64     `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
	Name      string    `json:"name,omitempty"`
	Active    bool      `json:"active,omitempty"`
	Report    Report    `json:"report"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Timezone  string    `json:"timezone,omitempty"`
	Interval  string    `json:"interval,omitempty"`
	WorkDays  bool      `json:"work_days,omitempty"`
}
