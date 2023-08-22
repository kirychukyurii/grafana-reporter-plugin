package models

import (
	"encoding/json"
	"fmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

type ReporterAppSetting struct {
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
