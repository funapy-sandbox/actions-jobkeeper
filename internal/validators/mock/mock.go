package mock

import (
	"context"

	"github.com/funapy-sandbox/actions-jobkeeper/internal/validators"
)

type Validator struct {
	ValidateFunc func(ctx context.Context) error
}

func (v *Validator) Validate(ctx context.Context) error {
	return v.ValidateFunc(ctx)
}

var (
	_ validators.Validator = &Validator{}
)
