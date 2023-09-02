package model

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

type ReporterAppSetting struct {
	GrafanaBaseURL string
	BasicAuth      BasicAuth

	TemporaryDirectory string
	WorkersCount       int
	Browser            BrowserSettings

	InsecureSkipVerify bool
}

type BrowserSettings struct {
	BinPath string
	Url     string
}

type BasicAuth struct {
	Username string
	Password string
}

func (s *ReporterAppSetting) Load(config backend.AppInstanceSettings) error {
	if config.JSONData != nil && len(config.JSONData) > 1 {
		if err := json.Unmarshal(config.JSONData, s); err != nil {
			return fmt.Errorf("could not unmarshal AppInstanceSettings json: %w", err)
		}
	}

	return nil
}

func (s *ReporterAppSetting) Validate() error {
	return nil
}

func (a *BasicAuth) String() string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", a.Username, a.Password)))
}
