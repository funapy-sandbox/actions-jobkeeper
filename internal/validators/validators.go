package validators

import (
	"context"
	"errors"
)

var (
	ErrValidate = errors.New("validate error")
)

type JobStatus uint8

const (
	JobStatusUnknown = iota
	JobStatusPending
	JobStatusError
	JobStatusFailure
	JobStatusSuccess
)

func (js JobStatus) String() string {
	switch js {
	case JobStatusPending:
		return "pending"
	case JobStatusError:
		return "error"
	case JobStatusFailure:
		return "failure"
	case JobStatusSuccess:
		return "success"
	default:
		return "unknown"
	}
}

type Status struct {
	Context     string
	Description string
	State       JobStatus
}

type Validator interface {
	Validate(ctx context.Context) error
}
