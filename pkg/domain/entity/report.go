package entity

import "time"

// ReportDashboardTimeRange represents the time range from a dashboard on a Grafana report.
type ReportDashboardTimeRange struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// ReportDashboard represents a dashboard on a Grafana report.
type ReportDashboard struct {
	ID        int64                    `json:"id,omitempty"`
	UID       string                   `json:"uid,omitempty"`
	Name      string                   `json:"name,omitempty"`
	TimeRange ReportDashboardTimeRange `json:"time_range"`
	Variables map[string]string        `json:"variables"`
}

// ReportOptions represents the options for a Grafana report.
type ReportOptions struct {
	Orientation        string `json:"orientation,omitempty"`
	Layout             string `json:"layout,omitempty"`
	EnableDashboardURL bool   `json:"enable_dashboard_url,omitempty"`
	EnableCSV          bool   `json:"enable_csv,omitempty"`
	ScaleFactor        int    `json:"scale_factor,omitempty"`
}

type Report struct {
	ID         int64           `json:"id,omitempty"`
	Name       string          `json:"name"`
	OrgID      int64           `json:"org_id,omitempty"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
	DeletedAt  time.Time       `json:"deleted_at"`
	State      string          `json:"state,omitempty"`
	Dashboard  ReportDashboard `json:"dashboard,omitempty"`
	Recipients []string        `json:"recipients,omitempty"`
	ReplyTo    string          `json:"reply_to,omitempty"`
	Message    string          `json:"message,omitempty"`
	Options    ReportOptions   `json:"options"`
}
