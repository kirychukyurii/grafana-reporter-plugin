package db

import (
	"database/sql"
	"fmt"
	"path/filepath"

	_ "modernc.org/sqlite"

	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/util"
)

const (
	databaseFolder   = "/tmp/reporter"
	databaseFileName = "database.db"
)

type DB struct {
	SQL *sql.DB
}

func New() (*DB, error) {
	databaseURL := filepath.Join(databaseFolder, databaseFileName)
	if err := util.EnsureDirRW(databaseFolder); err != nil {
		return nil, fmt.Errorf("ensure directory rw: %v", err)
	}

	f, err := util.Create(databaseURL)
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
