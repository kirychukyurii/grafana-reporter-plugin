package grafana

type GridPos struct {
	H int `json:"h"`
	W int `json:"w"`
	X int `json:"x"`
	Y int `json:"y"`
}

type PanelTargets struct {
	Datasource Datasource `json:"datasource"`
	QueryType  string     `json:"queryType"`
	RefId      string     `json:"refId"`
}

type Panel struct {
	Id         int            `json:"id"`
	Title      string         `json:"title"`
	Type       string         `json:"type"`
	GridPos    GridPos        `json:"gridPos"`
	Datasource Datasource     `json:"datasource"`
	Targets    []PanelTargets `json:"targets"`
}

// RowPanels represents a container for Panels
type RowPanels struct {
	Id        int     `json:"id"`
	Title     string  `json:"title"`
	Type      string  `json:"type"`
	Collapsed bool    `json:"collapsed"`
	GridPos   GridPos `json:"gridPos"`
	Panels    []Panel `json:"panels"`
}
