package cli

import (
	"context"
	"errors"
	"time"

	"github.com/funapy-sandbox/actions-jobkeeper/internal/github"
	"github.com/funapy-sandbox/actions-jobkeeper/internal/validators"
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
	targetJobName       string
)

func validateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate github actions job",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				targetJobName = defaultJobName
				return nil
			}
			targetJobName = args[0]
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			statusValidator := status.CreateValidator(github.NewClient(ctx, ghToken),
				status.WithTargetJob(targetJobName),
				status.WithGitHubOwnerAndRepo(ghOwner, ghRepo),
				status.WithGitHubRef(ghRef),
			)
			return doValidateCmd(ctx, cmd, statusValidator)
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

func doValidateCmd(ctx context.Context, logger logger, vs ...validators.Validator) error {
	timeoutT := time.NewTicker(time.Duration(timeoutSecond) * time.Second)
	defer timeoutT.Stop()

	invalT := time.NewTicker(time.Duration(validateInvalSecond) * time.Second)
	defer invalT.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timeoutT.C:
			return errors.New("validation is timeout")
		case <-invalT.C:
			logger.Println("start to validate")
			var successCnt int
			for _, validator := range vs {
				err := validator.Validate(ctx)
				if err != nil {
					if !errors.Is(err, validators.ErrValidate) {
						return err
					}
					logger.PrintErrln(err)
					break
				} else {
					successCnt++
				}
			}
			logger.Println("finish to validate")
			if successCnt == len(vs) {
				return nil
			}
		}
	}
}
