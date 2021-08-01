package status

import (
	"context"

	"github.com/funapy-sandbox/actions-jobkeeper/internal/github"
	"github.com/funapy-sandbox/actions-jobkeeper/internal/validators"
)

const (
	successState = "success"
	errorState   = "error"
	pendingState = "pending"
)

type statusValidator struct {
	token         string
	repo          string
	owner         string
	ref           string
	targetJobName string
	client        github.Client
}

func CreateValidator(c github.Client, opts ...Option) validators.Validator {
	sv := &statusValidator{
		client: c,
	}
	for _, opt := range opts {
		opt(sv)
	}
	return sv
}

func (sv *statusValidator) Validate(ctx context.Context) error {
	status, _, err := sv.client.ListStatuses(ctx, sv.owner, sv.repo, sv.ref, &github.ListOptions{})
	if err != nil {
		return err
	}

	// When there is no job other than the target job.
	if len(status) == 1 {
		return nil
	}

	var successJobCnt int
	for _, status := range status {
		if status.Context == nil || status.State == nil {
			continue
		}
		if *status.Context != sv.targetJobName && *status.State == successState {
			successJobCnt++
		}
	}
	if len(status)-1 != successJobCnt {
		return validators.ErrValidate
	}
	return nil
}
