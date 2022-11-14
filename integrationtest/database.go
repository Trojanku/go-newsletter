package integrationtest

import (
	"Goo/storage"
	"Goo/utils"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

var once sync.Once

// CreateDatabase for testing.
// Usage:
//
//	db, cleanup := CreateDatabase()
//	defer cleanup()
//	…
func CreateDatabase() (*storage.Database, func()) {
	utils.MustLoad("../.env-test")

	once.Do(initDatabase)

	db, cleanup := connect("postgres")
	defer cleanup()

	dropConnections(db.DB, "template1")

	name := utils.GetStringOrDefault("DB_NAME", "test")
	dropConnections(db.DB, name)
	db.DB.MustExec(`drop database if exists ` + name)
	db.DB.MustExec(`create database ` + name)

	return connect(name)
}

func initDatabase() {
	db, cleanup := connect("template1")
	defer cleanup()

	if err := db.Ping(context.Background()); err != nil {
		time.Sleep(100 * time.Millisecond)
	}

	url := fmt.Sprintf("postgresql://%v:%v@%v:%v/%v?sslmode=disable",
		utils.GetStringOrDefault("DB_USER", "test"),
		utils.GetStringOrDefault("DB_PASSWORD", ""),
		utils.GetStringOrDefault("DB_HOST", "localhost"),
		utils.GetIntOrDefault("DB_PORT", 5432),
		"template1",
	)

	migrationPath := utils.GetStringOrDefault("MIGRATIONS_PATH", "storage/migration")

	m, err := migrate.New(fmt.Sprintf("file://%s", migrationPath), url)
	if err != nil {
		panic(err)
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		panic(err)
	}
	if err = m.Down(); err != nil && err != migrate.ErrNoChange {
		panic(err)
	}
	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		panic(err)
	}

	if err = db.DB.Close(); err != nil {
		panic(err)
	}
}

func connect(name string) (*storage.Database, func()) {
	db := storage.NewDatabase(storage.NewDatabaseOptions{
		Host:               utils.GetStringOrDefault("DB_HOST", "localhost"),
		Port:               utils.GetIntOrDefault("DB_PORT", 5432),
		User:               utils.GetStringOrDefault("DB_USER", "test"),
		Password:           utils.GetStringOrDefault("DB_PASSWORD", ""),
		Name:               name,
		MaxOpenConnections: 10,
		MaxIdleConnections: 10,
		Log:                nil,
	})
	if err := db.Connect(); err != nil {
		panic(err)
	}
	return db, func() {
		if err := db.DB.Close(); err != nil {
			panic(err)
		}
	}
}

func dropConnections(db *sqlx.DB, name string) {
	db.MustExec(`
		select pg_terminate_backend(pg_stat_activity.pid)
		from pg_stat_activity
		where pg_stat_activity.datname = $1 and pid <> pg_backend_pid()`, name)
}
