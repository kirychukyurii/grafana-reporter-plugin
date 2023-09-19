package grafana

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/domain/entity"
)

type DashboardMeta struct {
	IsStarred bool   `json:"isStarred"`
	Slug      string `json:"slug"`
	Folder    int64  `json:"folderId"`
	FolderUID string `json:"folderUid"`
	URL       string `json:"url"`
}

type DashboardModel struct {
	Uid      string  `json:"uid"`
	Title    string  `json:"title"`
	Time     Time    `json:"time"`
	Timezone string  `json:"timezone"`
	Panels   []Panel `json:"panels"`
}

// Dashboard represents a Grafana dashboard.
type Dashboard struct {
	FolderID int64          `json:"folderId"`
	Meta     DashboardMeta  `json:"meta"` // read-only
	Model    DashboardModel `json:"dashboard"`
}

// Dashboard get and create Dashboard struct from Grafana internal JSON dashboard definition
func (c *Client) Dashboard(ctx context.Context, opts entity.DashboardOpts) (*Dashboard, error) {
	var dashboard Dashboard

	dashboardUrl := fmt.Sprintf("%s/api/dashboards/uid/%s", c.Setting.GrafanaBaseURL, opts.DashboardID)
	if len(opts.Variables) > 0 {
		dashboardUrl = fmt.Sprintf("%s?%s", dashboardUrl, opts.EncodeVariables())
	}

	var u url.Values
	u.Encode()

	body, err := c.Request(ctx, "GET", dashboardUrl, nil)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(body, &dashboard); err != nil {
		return nil, fmt.Errorf("parse recieved dashboard json: %v", err)
	}

	return &dashboard, nil
}
