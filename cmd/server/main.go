// Package main is the entry point to the server. It reads configuration, sets up logging and error handling,
// handles signals from the OS, and starts and stops the server.
package main

import (
	"Goo/jobs"
	"Goo/messaging"
	"Goo/server"
	"Goo/storage"
	"Goo/utils"
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/smithy-go/logging"
	"go.uber.org/zap"

	"golang.org/x/sync/errgroup"
)

// release is set through the linker at build time, generally from a git sha.
// Used for logging and error reporting.
var release string

func main() {
	os.Exit(start())
}

func start() int {
	_ = utils.Load()

	logEnv := utils.GetStringOrDefault("LOG_ENV", "development")
	log, err := createLogger(logEnv)

	if err != nil {
		fmt.Println("Error setting up the logger: ", err)
		return 1
	}

	log = log.With(zap.String("release", release))

	defer func() {
		// If we cannot sync, there's probably something wrong with outputting logs,
		// so we probably cannot write using fmt.Println either. So just ignore the error.
		_ = log.Sync()
	}()

	host := utils.GetStringOrDefault("HOST", "localhost")
	port := utils.GetIntOrDefault("PORT", 8080)

	awsConfig, err := config.LoadDefaultConfig(context.Background(),
		config.WithLogger(createAWSLogAdapter(log)),
		config.WithEndpointResolverWithOptions(createAWSEndpointResolver()),
	)
	if err != nil {
		log.Info("Error creating AWS config", zap.Error(err))
	}

	registry := prometheus.NewRegistry()
	registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	registry.MustRegister(collectors.NewGoCollector())

	queue := createQueue(log, awsConfig)
	db := createDatabase(log, registry)
	if err = db.Connect(); err != nil {
		log.Info("Error connecting to database", zap.Error(err))
		return 1
	}

	s := server.New(server.Options{
		AdminPassword:   utils.GetStringOrDefault("ADMIN_PASSWORD", "eyDawVH9LLZtaG2q"),
		Database:        db,
		Host:            host,
		Log:             log,
		MetricsPassword: utils.GetStringOrDefault("METRICS_PASSWORD", "12345678"),
		Metrics:         registry,
		Port:            port,
		Queue:           queue,
	})

	r := jobs.NewRunner(jobs.NewRunnerOptions{
		Emailer: createEmailer(log, host, port),
		Log:     log,
		Metrics: registry,
		Queue:   queue,
	})

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		if err := s.Start(); err != nil {
			log.Info("Error starting server", zap.Error(err))
			return err
		}
		return nil
	})

	eg.Go(func() error {
		r.Start(ctx)
		return nil
	})

	<-ctx.Done()

	eg.Go(func() error {
		if err := s.Stop(); err != nil {
			log.Info("Error stopping server", zap.Error(err))
			return err
		}
		return nil
	})
	if err := eg.Wait(); err != nil {
		return 1
	}
	return 0
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

func createAWSLogAdapter(log *zap.Logger) logging.LoggerFunc {
	return func(classification logging.Classification, format string, v ...interface{}) {
		switch classification {
		case logging.Debug:
			log.Sugar().Debugf(format, v...)
		case logging.Warn:
			log.Sugar().Warnf(format, v...)
		}
	}
}

// createAWSEndpointResolver used for local development endpoints.
// See https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/endpoints/
func createAWSEndpointResolver() aws.EndpointResolverWithOptionsFunc {
	sqsEndpointURL := utils.GetStringOrDefault("SQS_ENDPOINT_URL", "")

	return func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if sqsEndpointURL != "" && service == sqs.ServiceID {
			return aws.Endpoint{
				URL: sqsEndpointURL,
			}, nil
		}
		// Fallback to default endpoint
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	}
}

func createDatabase(log *zap.Logger, registry *prometheus.Registry) *storage.Database {
	return storage.NewDatabase(storage.NewDatabaseOptions{
		Host:                  utils.GetStringOrDefault("DB_HOST", "localhost"),
		Port:                  utils.GetIntOrDefault("DB_PORT", 5432),
		User:                  utils.GetStringOrDefault("DB_USER", ""),
		Password:              utils.GetStringOrDefault("DB_PASSWORD", ""),
		Name:                  utils.GetStringOrDefault("DB_NAME", ""),
		MaxOpenConnections:    utils.GetIntOrDefault("DB_MAX_OPEN_CONNECTIONS", 10),
		MaxIdleConnections:    utils.GetIntOrDefault("DB_MAX_IDLE_CONNECTIONS", 10),
		ConnectionMaxLifetime: utils.GetDurationOrDefault("DB_CONNECTION_MAX_LIFETIME", time.Hour),
		Log:                   log,
		Metrics:               registry,
	})
}

func createQueue(log *zap.Logger, awsConfig aws.Config) *messaging.Queue {
	return messaging.NewQueue(messaging.NewQueueOptions{
		Config:   awsConfig,
		Log:      log,
		Name:     utils.GetStringOrDefault("QUEUE_NAME", "jobs"),
		WaitTime: utils.GetDurationOrDefault("QUEUE_WAIT_TIME", 20*time.Second),
	})
}

func createEmailer(log *zap.Logger, host string, port int) *messaging.Emailer {
	return messaging.NewEmailer(messaging.NewEmailerOptions{
		BaseURL:                   utils.GetStringOrDefault("BASE_URL", fmt.Sprintf("http://%v:%v", host, port)),
		Host:                      utils.GetStringOrDefault("EMAIL_HOST", "localhost"),
		Port:                      utils.GetIntOrDefault("EMAIL_PORT", 1025),
		MarketingUsername:         utils.GetStringOrDefault("MARKETING_USERNAME", "Goo bot"),
		MarketingPassword:         utils.GetStringOrDefault("MARKETING_EMAIL_PASSWORD", ""),
		TransactionalUsername:     utils.GetStringOrDefault("TRANSACTIONAL_USERNAME", "Goo bot"),
		TransactionalPassword:     utils.GetStringOrDefault("TRANSACTIONAL_PASSWORD", ""),
		MarketingEmailAddress:     utils.GetStringOrDefault("MARKETING_EMAIL", "goo.marketing@example.com"),
		MarketingEmailName:        utils.GetStringOrDefault("MARKETING_EMAIL_NAME", ""),
		TransactionalEmailAddress: utils.GetStringOrDefault("TRANSACTIONAL_EMAIL", "goo.transactional@example.com"),
		TransactionalEmailName:    utils.GetStringOrDefault("TRANSACTIONAL_EMAIL_NAME", ""),
		Log:                       log,
	})
}
