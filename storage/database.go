package storage

import (
	"context"
	"embed"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"time"

	"go.uber.org/zap"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

type Database struct {
	DB                    *sqlx.DB
	host                  string
	port                  int
	user                  string
	password              string
	name                  string
	maxOpenConnections    int
	maxIdleConnections    int
	connectionMaxLifetime time.Duration
	connectionMaxIdleTime time.Duration
	log                   *zap.Logger
	metrics               *prometheus.Registry
}

type NewDatabaseOptions struct {
	Host                  string
	Port                  int
	User                  string
	Password              string
	Name                  string
	MaxOpenConnections    int
	MaxIdleConnections    int
	ConnectionMaxLifetime time.Duration
	ConnectionMaxIdleTime time.Duration
	Log                   *zap.Logger
	Metrics               *prometheus.Registry
}

func NewDatabase(opts NewDatabaseOptions) *Database {
	if opts.Log == nil {
		opts.Log = zap.NewNop()
	}
	if opts.Metrics == nil {
		opts.Metrics = prometheus.NewRegistry()
	}
	return &Database{
		host:                  opts.Host,
		port:                  opts.Port,
		user:                  opts.User,
		password:              opts.Password,
		name:                  opts.Name,
		maxOpenConnections:    opts.MaxOpenConnections,
		maxIdleConnections:    opts.MaxIdleConnections,
		connectionMaxIdleTime: opts.ConnectionMaxIdleTime,
		connectionMaxLifetime: opts.ConnectionMaxLifetime,
		log:                   opts.Log,
		metrics:               opts.Metrics,
	}
}

func (d *Database) Connect() error {
	d.log.Info("Connecting to database", zap.String("url", d.createDataSourceName(false)))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	d.DB, err = sqlx.ConnectContext(ctx, "pgx", d.createDataSourceName(true))
	if err != nil {
		return err
	}

	d.log.Debug("Setting connection pool options",
		zap.Int("max open connections", d.maxOpenConnections),
		zap.Int("map idle connections", d.maxIdleConnections),
		zap.Duration("connection max lifetime", d.connectionMaxLifetime),
		zap.Duration("connection max idle time", d.connectionMaxIdleTime))
	d.DB.SetMaxOpenConns(d.maxOpenConnections)
	d.DB.SetMaxIdleConns(d.maxIdleConnections)
	d.DB.SetConnMaxLifetime(d.connectionMaxLifetime)
	d.DB.SetConnMaxIdleTime(d.connectionMaxIdleTime)

	d.metrics.MustRegister(collectors.NewDBStatsCollector(d.DB.DB, d.name))

	return nil
}

func (d *Database) createDataSourceName(withPassword bool) string {
	password := d.password
	if !withPassword {
		password = "xxx"
	}
	return fmt.Sprintf("postgresql://%v:%v@%v:%v/%v?sslmode=disable", d.user, password, d.host, d.port, d.name)
}

func (d *Database) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	if err := d.DB.PingContext(ctx); err != nil {
		return err
	}

	_, err := d.DB.ExecContext(ctx, `select 1`)
	return err
}

//go:embed migrations
var migrations embed.FS

func (d *Database) MigrateTo(ctx context.Context, version uint) error {
	m, err := d.getMigrate()
	if err != nil {
		return err
	}
	return m.Migrate(version)
}

func (d *Database) MigrateUp(ctx context.Context) error {
	m, err := d.getMigrate()
	if err != nil {
		return err
	}
	return m.Up()
}

func (d *Database) getMigrate() (*migrate.Migrate, error) {
	filesDriver, err := iofs.New(migrations, "migrations")
	if err != nil {
		return nil, err
	}
	driver, err := postgres.WithInstance(d.DB.DB, &postgres.Config{})
	if err != nil {
		return nil, err
	}
	return migrate.NewWithInstance("iofs", filesDriver, "goo", driver)
}
