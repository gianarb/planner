package planner

import (
	"context"

	"go.uber.org/zap"
)

// Procedure describe every single step to be executed. It is the smallest unit
// of work in a plan.
type Procedure interface {
	// Name identifies a specific procedure.
	Name() string
	// Do execute the business logic for a specific procedure.
	Do(ctx context.Context) ([]Procedure, error)
}

// Plan describe a set of procedures and the way to calculate them
type Plan interface {
	// Create returns the list of procedures that needs to be executed.
	Create(ctx context.Context) ([]Procedure, error)
	// Name identifies a specific plan
	Name() string
}

// Loggable describes an object that support its custom logger
type Loggable interface {
	// WithLogger set the logger to the object
	WithLogger(*zap.Logger)
}
