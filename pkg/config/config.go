package config

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

type ReporterAppConfig struct {
	WorkersCount  int    `json:"workers_count,omitempty" env:"GF_PLUGIN_WORKERS_COUNT"`
	DataDirectory string `json:"data_directory,omitempty" env:"GF_PLUGIN_DATA_DIRECTORY"`

	MailConfig     MailConfig
	GrafanaConfig  GrafanaConfig
	DatabaseConfig DatabaseConfig
	BrowserConfig  BrowserConfig
}

type MailConfig struct {
	Host     string `json:"mail_host" env:"GF_PLUGIN_MAIL_HOST"`
	Port     int    `json:"mail_port" env:"GF_PLUGIN_MAIL_PORT"`
	Username string `json:"mail_username" env:"GF_PLUGIN_MAIL_USERNAME"`
	Password string `json:"mail_password" env:"GF_PLUGIN_MAIL_PASSWORD"`
}

type GrafanaConfig struct {
	URL                string `json:"grafana_url,omitempty" env:"GF_PLUGIN_GRAFANA_URL"`
	InsecureSkipVerify bool   `json:"grafana_insecure_skip_verify,omitempty" env:"GF_PLUGIN_GRAFANA_INSECURE_SKIP_VERIFY"`

	APIToken string `json:"grafana_api_token,omitempty" env:"GF_PLUGIN_GRAFANA_API_TOKEN"`
	Username string `json:"grafana_username,omitempty" env:"GF_PLUGIN_GRAFANA_USERNAME"`
	Password string `json:"grafana_password,omitempty" env:"GF_PLUGIN_GRAFANA_PASSWORD"`
}

type DatabaseConfig struct {
	MaxBatchSize    int           `json:"database_max_batch_size,omitempty" env:"GF_PLUGIN_DATABASE_MAX_BATCH_SIZE"`
	MaxBatchDelay   time.Duration `json:"database_max_batch_delay,omitempty" env:"GF_PLUGIN_DATABASE_MAX_BATCH_DELAY"`
	InitialMmapSize int           `json:"database_initial_mmap_size,omitempty" env:"GF_PLUGIN_DATABASE_INITIAL_MMAP_SIZE"`
	EncryptionKey   []byte        `json:"database_encryption_key,omitempty" env:"GF_PLUGIN_DATABASE_ENCRYPTION_KEY"`
}

type BrowserConfig struct {
	Type    string `json:"browser_type,omitempty" env:"GF_PLUGIN_BROWSER_TYPE"`
	BinPath string `json:"browser_bin_path,omitempty" env:"GF_PLUGIN_BROWSER_BIN_PATH"`
	URL     string `json:"browser_url,omitempty" env:"GF_PLUGIN_BROWSER_URL"`
}

func New(settings backend.AppInstanceSettings) (*ReporterAppConfig, error) {
	var config ReporterAppConfig

	if settings.JSONData != nil && len(settings.JSONData) > 1 {
		if err := json.Unmarshal(settings.JSONData, &config); err != nil {
			return nil, fmt.Errorf("could not unmarshal AppInstanceSettings json: %w", err)
		}
	}

	if err := cleanenv.ReadEnv(&config); err != nil {
		return nil, fmt.Errorf("read env: %v", err)
	}

	return &config, nil
}

func (a *GrafanaConfig) BasicAuth() string {
	return fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", a.Username, a.Password))))
}
