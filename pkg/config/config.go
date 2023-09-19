package config

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

type ReporterAppSetting struct {
	GrafanaBaseURL     string
	BasicAuth          BasicAuth
	InsecureSkipVerify bool

	TemporaryDirectory string
	WorkersCount       int
	Browser            BrowserSettings
}

type BrowserSettings struct {
	BinPath string
	Url     string
}

type BasicAuth struct {
	Username string
	Password string
}

func New(config backend.AppInstanceSettings) (*ReporterAppSetting, error) {
	var setting ReporterAppSetting

	if config.JSONData != nil && len(config.JSONData) > 1 {
		if err := json.Unmarshal(config.JSONData, &setting); err != nil {
			return nil, fmt.Errorf("could not unmarshal AppInstanceSettings json: %w", err)
		}
	}

	setting.GrafanaBaseURL = "http://localhost:3000"
	setting.BasicAuth.Username = "admin"
	setting.BasicAuth.Password = "admin"
	setting.WorkersCount = 10
	setting.TemporaryDirectory = "/opt/reporter/tmp"
	setting.Browser.Url = "chrome"

	return &setting, nil
}

func (a *BasicAuth) String() string {
	return fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", a.Username, a.Password))))
}
