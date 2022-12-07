// Package jobs has a Runner that can run registered jobs in parallel.
package jobs

import (
	"Goo/messaging"
	"Goo/model"
	"context"
	"go.uber.org/zap"
)

type Runner struct {
	emailer *messaging.Emailer
	jobs    map[string]Func
	log     *zap.Logger
	queue   *messaging.Queue
}

type NewRunnerOptions struct {
	Emailer *messaging.Emailer
	Log     *zap.Logger
	Queue   *messaging.Queue
}

func NewRunner(opts NewRunnerOptions) *Runner {
	if opts.Log == nil {
		opts.Log = zap.NewNop()
	}

	return &Runner{
		emailer: opts.Emailer,
		jobs:    map[string]Func{},
		log:     opts.Log,
		queue:   opts.Queue,
	}
}

// Func is the actual work to do in a job.
// The given context is the root context of the runner, which may be cancelled.
type Func = func(ctx context.Context, message model.Message) error
