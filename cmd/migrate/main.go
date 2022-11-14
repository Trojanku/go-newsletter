package main

import (
	"Goo/storage"
	"Goo/utils"
	"fmt"
	"go.uber.org/zap"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	os.Exit(start())
}

// value := os.Getenv(key)
func start() int {
	_ = utils.Load()

	logEnv := getStringOrDefault("LOG_ENV", "development")
	log, err := createLogger(logEnv)
	if err != nil {
		fmt.Println("Error setting up the logger:", err)
		return 1
	}

	if len(os.Args) < 2 {
		log.Warn("Usage: migrate up|down|to")
		return 1
	}

	if os.Args[1] == "to" && len(os.Args) < 3 {
		log.Info("Usage: migrate to <version>")
		return 1
	}

	db := storage.NewDatabase(storage.NewDatabaseOptions{
		Host:     getStringOrDefault("DB_HOST", "localhost"),
		Port:     getIntOrDefault("DB_PORT", 5432),
		User:     getStringOrDefault("DB_USER", ""),
		Password: getStringOrDefault("DB_PASSWORD", ""),
		Name:     getStringOrDefault("DB_NAME", ""),
	})

	if err = db.Connect(); err != nil {
		log.Error("Error connection to database", zap.Error(err))
		return 1
	}

	url := fmt.Sprintf("postgresql://%v:%v@%v:%v/%v?sslmode=disable",
		getStringOrDefault("DB_USER", ""),
		getStringOrDefault("DB_PASSWORD", ""),
		getStringOrDefault("DB_HOST", "localhost"),
		getIntOrDefault("DB_PORT", 5432),
		getStringOrDefault("DB_NAME", ""),
	)

	m, err := migrate.New("file://storage/migrations/", url)
	if err != nil {
		log.Error("Error", zap.Error(err))
		return 1
	}

	switch os.Args[1] {
	case "up":
		err = m.Up()
	case "down":
		err = m.Down()
	case "to":
		version, _ := strconv.Atoi(os.Args[2])
		err = m.Migrate(uint(version))
	default:
		log.Error("Unknown command", zap.String("name", os.Args[1]))
		return 1
	}
	if err != nil {
		log.Error("Error migrating", zap.Error(err))
		return 1
	}
	return 0
}

func getStringOrDefault(name, defaultV string) string {
	v, ok := os.LookupEnv(name)
	if !ok {
		return defaultV
	}
	return v
}

func getIntOrDefault(name string, defaultV int) int {
	v, ok := os.LookupEnv(name)
	if !ok {
		return defaultV
	}
	vAsInt, err := strconv.Atoi(v)
	if err != nil {
		return defaultV
	}
	return vAsInt
}

func createLogger(env string) (*zap.Logger, error) {
	switch env {
	case "production":
		return zap.NewProduction()
	case "development":
		return zap.NewDevelopment()
	default:
		return zap.NewNop(), nil
	}
}
