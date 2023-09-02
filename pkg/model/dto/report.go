package dto

import "time"

type Report struct {
	ID         int64
	OrgID      int64
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  time.Time
	State      string
	Dashboards string
	Recipients string
	ReplyTo    string
	Message    string
	Options    string
}
