package model

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

func (s *ReporterAppSetting) Load(config backend.AppInstanceSettings) error {
	if config.JSONData != nil && len(config.JSONData) > 1 {
		if err := json.Unmarshal(config.JSONData, s); err != nil {
			return fmt.Errorf("could not unmarshal AppInstanceSettings json: %w", err)
		}
	}

	s.GrafanaBaseURL = "https://cloud.webitel.ua/grafana"
	s.BasicAuth.Username = "srvadm"
	s.BasicAuth.Password = "whogAQgABPkt3wzQ"
	s.WorkersCount = 10
	s.TemporaryDirectory = "/opt/reporter/tmp"
	s.Browser.Url = "chrome"

	return nil
}

func (s *ReporterAppSetting) Validate() error {
	return nil
}

func (a *BasicAuth) String() string {
	return fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", a.Username, a.Password))))
}
