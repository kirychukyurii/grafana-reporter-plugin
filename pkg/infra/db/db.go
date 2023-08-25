package db

import (
	"database/sql"
	"fmt"
	"path/filepath"

	_ "modernc.org/sqlite"

	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/utils"
)

const (
	databaseFolder   = "/opt/reporter"
	databaseFileName = "database.db"
)

type DB struct {
	SQL *sql.DB
}

func New() (*DB, error) {
	databaseURL := filepath.Join(databaseFolder, databaseFileName)
	if err := utils.EnsureDirRW(databaseFolder); err != nil {
		return nil, fmt.Errorf("ensure directory rw: %v", err)
	}

	f, err := utils.Create(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("create database file: %v", err)
	}

	defer f.Close()

	db, err := sql.Open("sqlite", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("open database: %v", err)
	}

	return &DB{
		SQL: db,
	}, nil
}
