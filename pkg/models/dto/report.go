package dto

import "time"

type Report struct {
	ID           int
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    time.Time
	State        string
	Recipients   string
	ReplyTo      string
	Message      string
	Orientation  string
	Layout       string
	DashboardURL bool
	CSV          bool
	ScaleFactor  int
}
