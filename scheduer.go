package planner

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Scheduler takes a plan and it executes it.
type Scheduler struct {
	// stepCounter keep track of the number of steps exectued by the scheduler.
	// It is used for debug and logged out at the end of every execution.
	stepCounter int
	// logger is an instance of the zap.Logger
	logger *zap.Logger
}

// NewScheduler creates a new scheduler
func NewScheduler() *Scheduler {
	return &Scheduler{
		stepCounter: 0,
		logger:      zap.NewNop(),
	}
}

// WithLogger allows you to pass a logger from the outside.
func (s *Scheduler) WithLogger(logger *zap.Logger) {
	s.logger = logger
}

// Execute takes a plan it executes it
func (s *Scheduler) Execute(ctx context.Context, p Plan) error {
	uuidGenerator := uuid.New()
	logger := s.logger.With(zap.String("execution_id", uuidGenerator.String()))
	start := time.Now()
	if loggableP, ok := p.(Loggable); ok {
		loggableP.WithLogger(logger)
	}
	logger.Info("Started execution plan " + p.Name())
	s.stepCounter = 0
	for {
		steps, err := p.Create(ctx)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
		if len(steps) == 0 {
			break
		}
		err = s.react(ctx, steps, logger)
		if err != nil {
			logger.Error(err.Error(), zap.String("execution_time", time.Since(start).String()), zap.Int("step_executed", s.stepCounter))
			return err
		}
	}
	logger.Info("Plan executed without errors.", zap.String("execution_time", time.Since(start).String()), zap.Int("step_executed", s.stepCounter))
	return nil
}

// react is a recursive function that goes over all the steps and the one
// returned by previous steps until the plan does not return anymore steps
func (s *Scheduler) react(ctx context.Context, steps []Procedure, logger *zap.Logger) error {
	for _, step := range steps {
		s.stepCounter = s.stepCounter + 1
		if loggableS, ok := step.(Loggable); ok {
			loggableS.WithLogger(logger)
		}
		innerSteps, err := step.Do(ctx)
		if err != nil {
			logger.Error("Step failed.", zap.String("step", step.Name()), zap.Error(err))
			return err
		}
		if len(innerSteps) > 0 {
			if err := s.react(ctx, innerSteps, logger); err != nil {
				return err
			}
		}
	}
	return nil
}
