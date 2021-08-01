package cli

import (
	"context"
	"errors"
	"time"

	"github.com/funapy-sandbox/actions-jobkeeper/internal/github"
	"github.com/funapy-sandbox/actions-jobkeeper/internal/validators/status"
	"github.com/spf13/cobra"
)

const defaultJobName = "jobkeeper"

// Tease variables will be set by command line flags.
var (
	ghOwner             string
	ghRepo              string
	ghRef               string
	timeoutSecond       uint
	validateInvalSecond uint
)

func validateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate github actions job",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("target job context is not set")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			jobName := defaultJobName
			if len(args) > 0 {
				jobName = args[0]
			}
			return doValidateCmd(cmd.Context(), jobName)
		},
	}

	cmd.PersistentFlags().StringVarP(&ghRepo, "repo", "r", "", "set github repository")
	cmd.MarkPersistentFlagRequired("repo")
	cmd.PersistentFlags().StringVarP(&ghOwner, "owner", "o", "", "set owner of github repository")
	cmd.MarkPersistentFlagRequired("owpner")
	cmd.PersistentFlags().StringVar(&ghRef, "ref", "", "set ref of github repository. the ref can be a SHA, a branch name, or tag name")
	cmd.MarkPersistentFlagRequired("ref")
	cmd.PersistentFlags().UintVar(&timeoutSecond, "timeout", 600, "set validate timeout second")
	cmd.MarkPersistentFlagRequired("timeout")
	cmd.PersistentFlags().UintVar(&validateInvalSecond, "interval", 120, "set validate interval second")
	cmd.MarkPersistentFlagRequired("timeout")

	return cmd
}

func doValidateCmd(ctx context.Context, targetJobName string) error {
	timeoutT := time.NewTicker(time.Duration(timeoutSecond) * time.Second)
	defer timeoutT.Stop()

	invalT := time.NewTicker(time.Duration(validateInvalSecond) * time.Second)
	defer invalT.Stop()

	statusValidator := status.CreateValidator(github.NewClient(ctx, ghToken),
		status.WithTargetJob(targetJobName),
		status.WithGitHubOwnerAndRepo(ghOwner, ghRepo),
		status.WithGitHubRef(ghRef),
	)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-timeoutT.C:
			return errors.New("timeout validate")
		case <-invalT.C:
			if err := statusValidator.Validate(ctx); err != nil {
				return err
			}
		}
	}
}
