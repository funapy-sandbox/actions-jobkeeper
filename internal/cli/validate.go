package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/funapy-sandbox/actions-jobkeeper/internal/github"
	"github.com/funapy-sandbox/actions-jobkeeper/internal/validators"
	"github.com/funapy-sandbox/actions-jobkeeper/internal/validators/status"
	"github.com/spf13/cobra"
)

const defaultJobName = "check-other-job-status"

// Tease variables will be set by command line flags.
var (
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
		PreRunE: func(cmd *cobra.Command, args []string) error {
			repo := os.Getenv("GITHUB_REPOSITORY")
			if len(repo) != 0 {
				ghRepo = repo
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			fmt.Println("env result")
			fmt.Println("****************************")
			fmt.Println(ghRef)
			fmt.Println(timeoutSecond)
			fmt.Println(validateInvalSecond)
			fmt.Println(targetJobName)
			fmt.Println(ghRepo)
			fmt.Println(os.Getenv("GITHUB_REPOSITORY"))
			fmt.Println(os.Getenv("GITHUB_REPOSITORY_OWNER"))
			fmt.Println("****************************")

			owner, repo := ownerAndRepository(ghRepo)
			if len(owner) == 0 || len(repo) == 0 {
				return fmt.Errorf("github owner or repository is empty. owner: %s, repository: %s", owner, repo)
			}

			fmt.Println(owner)
			fmt.Println(repo)
			fmt.Println("****************************")

			statusValidator := status.CreateValidator(github.NewClient(ctx, ghToken),
				status.WithTargetJob(targetJobName),
				status.WithGitHubOwnerAndRepo(owner, repo),
				status.WithGitHubRef(ghRef),
			)
			return doValidateCmd(ctx, cmd, statusValidator)
		},
	}

	cmd.PersistentFlags().StringVarP(&targetJobName, "job", "j", defaultJobName, "set target job name")

	cmd.PersistentFlags().StringVarP(&ghRepo, "repo", "r", "", "set github repository")

	cmd.PersistentFlags().StringVar(&ghRef, "ref", "", "set ref of github repository. the ref can be a SHA, a branch name, or tag name")

	cmd.PersistentFlags().UintVar(&timeoutSecond, "timeout", 600, "set validate timeout second")

	cmd.PersistentFlags().UintVar(&validateInvalSecond, "interval", 10, "set validate interval second")

	return cmd
}

func ownerAndRepository(str string) (owner string, repo string) {
	sp := strings.Split(str, "/")
	switch len(sp) {
	case 0:
		return "", ""
	case 1:
		return sp[0], ""
	case 2:
		return sp[0], sp[1]
	default:
		return sp[0], strings.Join(sp[1:], "/")
	}
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
			if successCnt == len(vs) {
				return nil
			}
		}
	}
}
