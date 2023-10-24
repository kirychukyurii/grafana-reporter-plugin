package sqlite

import (
	"embed"
	"fmt"
	"github.com/pressly/goose/v3"
)

//go:embed schema/*.sql
var embedMigrations embed.FS

func Migrate(database *DB) error {
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("sqlite"); err != nil {
		return fmt.Errorf("set dialect: %v", err)
	}

	if err := goose.Up(database.SQL, "schema"); err != nil {
		return fmt.Errorf("up migrations: %v", err)
	}

	return nil
}
