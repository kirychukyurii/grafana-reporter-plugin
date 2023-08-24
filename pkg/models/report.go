package models

import "time"

// ReportSchedule represents the schedule from a Grafana report.
type ReportSchedule struct {
	StartDate         *time.Time `json:"start_date,omitempty"`
	EndDate           *time.Time `json:"end_date,omitempty"`
	Frequency         string     `json:"frequency"`
	IntervalFrequency string     `json:"interval_frequency"`
	IntervalAmount    int64      `json:"interval_amount"`
	WorkdaysOnly      bool       `json:"workdays_only"`
	TimeZone          string     `json:"timezone"`
	DayOfMonth        string     `json:"day_of_month,omitempty"`
}

// ReportOptions represents the options for a Grafana report.
type ReportOptions struct {
	Orientation string `json:"orientation"`
	Layout      string `json:"layout"`
}

// ReportDashboardTimeRange represents the time range from a dashboard on a Grafana report.
type ReportDashboardTimeRange struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// ReportDashboardIdentifier represents the identifier for a dashboard on a Grafana report.
type ReportDashboardIdentifier struct {
	ID   int64  `json:"id,omitempty"`
	UID  string `json:"uid,omitempty"`
	Name string `json:"name,omitempty"`
}

// ReportDashboard represents a dashboard on a Grafana report.
type ReportDashboard struct {
	Dashboard ReportDashboardIdentifier `json:"dashboard"`
	TimeRange ReportDashboardTimeRange  `json:"time_range"`
	Variables map[string]string         `json:"report_variables"`
}

// Report represents a Grafana report.
type Report struct {
	// ReadOnly
	ID     int64  `json:"id,omitempty"`
	UserID int64  `json:"user_id,omitempty"`
	OrgID  int64  `json:"org_id,omitempty"`
	State  string `json:"state,omitempty"`

	Dashboards []ReportDashboard `json:"dashboards"`

	Name               string          `json:"name"`
	Recipients         string          `json:"recipients"`
	ReplyTo            string          `json:"reply_to"`
	Message            string          `json:"message"`
	Schedule           *ReportSchedule `json:"schedule"`
	Options            ReportOptions   `json:"options"`
	EnableDashboardURL bool            `json:"enable_dashboard_url"`
	EnableCSV          bool            `json:"enable_csv"`
	Formats            []string        `json:"formats"`
	ScaleFactor        int64           `json:"scale_factor"`
}
