package status

import (
	"reflect"
	"testing"

	"github.com/funapy-sandbox/actions-jobkeeper/internal/github"
	"github.com/funapy-sandbox/actions-jobkeeper/internal/github/mock"
	"github.com/funapy-sandbox/actions-jobkeeper/internal/validators"
)

func stringPtr(str string) *string {
	return &str
}

func TestCreateValidator(t *testing.T) {
	type args struct {
		c    github.Client
		opts []Option
	}
	type test struct {
		args args
		want validators.Validator
	}
	tests := map[string]test{
		"returns Validator when option is empty": func() test {
			c := &mock.Client{}
			return test{
				args: args{
					c: c,
				},
				want: &statusValidator{
					client: c,
				},
			}
		}(),
		"returns Validator when option is not empty": func() test {
			c := &mock.Client{}
			return test{
				args: args{
					c: c,
					opts: []Option{
						WithGitHubOwnerAndRepo("test-owner", "test-repo"),
						WithGitHubRef("sha"),
						WithTargetJob("job"),
					},
				},
				want: &statusValidator{
					client:        c,
					owner:         "test-owner",
					repo:          "test-repo",
					ref:           "sha",
					targetJobName: "job",
				},
			}
		}(),
		"returns Validator when there are duplicate options": func() test {
			c := &mock.Client{}
			return test{
				args: args{
					c: c,
					opts: []Option{
						WithGitHubOwnerAndRepo("test", "test-repo"),
						WithGitHubRef("sha"),
						WithGitHubRef("sha-01"),
						WithTargetJob("job"),
						WithTargetJob("job-01"),
					},
				},
				want: &statusValidator{
					client:        c,
					owner:         "test",
					repo:          "test-repo",
					ref:           "sha-01",
					targetJobName: "job-01",
				},
			}
		}(),
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := CreateValidator(tt.args.c, tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateValidator() = %v, want %v", got, tt.want)
			}
		})
	}
}

// func Test_statusValidator_Validate(t *testing.T) {
// 	type fields struct {
// 		token         string
// 		repo          string
// 		owner         string
// 		ref           string
// 		targetJobName string
// 		client        github.Client
// 	}
// 	type test struct {
// 		fields  fields
// 		ctx     context.Context
// 		wantErr bool
// 	}
// 	tests := map[string]func(*testing.T) test{
// 		"returns nil when validation is success": func(*testing.T) test {
// 			t.Helper()
//
// 			var (
// 				testOwner     = "test-owner"
// 				testRepo      = "test-repo"
// 				testRef       = "sha"
// 				testOpts      = &github.ListOptions{}
// 				targetJobName = "target-job"
// 			)
//
// 			c := &mock.Client{
// 				ListStatusesFunc: func(ctx context.Context, owner, repo, ref string, opts *github.ListOptions) ([]*github.RepoStatus, *github.Response, error) {
// 					if testOwner != owner {
// 						t.Errorf("ListStatuses owner = %s, want %s", owner, testOwner)
// 					}
// 					if testRepo != repo {
// 						t.Errorf("ListStatuses repo = %s, want %s", repo, testRepo)
// 					}
// 					if testRef != ref {
// 						t.Errorf("ListStatuses ref = %s, want %s", ref, testRef)
// 					}
// 					if !reflect.DeepEqual(opts, testOpts) {
// 						t.Errorf("ListStatuses opts = %v, want %v", opts, testOpts)
// 					}
//
// 					return []*github.RepoStatus{
// 						{
// 							Context: stringPtr("job-01"),
// 							State:   stringPtr(successState),
// 						},
// 						{
// 							Context: stringPtr("job-02"),
// 							State:   stringPtr(successState),
// 						},
// 						{
// 							Context: stringPtr(targetJobName),
// 							State:   stringPtr(pendingState),
// 						},
// 					}, nil, nil
// 				},
// 			}
// 			return test{
// 				fields: fields{
// 					client:        c,
// 					targetJobName: targetJobName,
// 					owner:         testOwner,
// 					repo:          testRepo,
// 					ref:           testRef,
// 				},
// 				wantErr: false,
// 			}
// 		},
// 		"returns nil when there is one job": func(*testing.T) test {
// 			t.Helper()
//
// 			var (
// 				testOwner     = "test-owner"
// 				testRepo      = "test-repo"
// 				testRef       = "sha"
// 				targetJobName = "target-job"
// 			)
//
// 			c := &mock.Client{
// 				ListStatusesFunc: func(ctx context.Context, owner, repo, ref string, opts *github.ListOptions) ([]*github.RepoStatus, *github.Response, error) {
// 					return []*github.RepoStatus{
// 						{
// 							Context: stringPtr(targetJobName),
// 							State:   stringPtr(pendingState),
// 						},
// 					}, nil, nil
// 				},
// 			}
// 			return test{
// 				fields: fields{
// 					client:        c,
// 					targetJobName: targetJobName,
// 					owner:         testOwner,
// 					repo:          testRepo,
// 					ref:           testRef,
// 				},
// 				wantErr: false,
// 			}
// 		},
// 		"returns error when ListStatuses returns an error": func(*testing.T) test {
// 			c := &mock.Client{
// 				ListStatusesFunc: func(ctx context.Context, owner, repo, ref string, opts *github.ListOptions) ([]*github.RepoStatus, *github.Response, error) {
// 					return nil, nil, errors.New("err")
// 				},
// 			}
// 			return test{
// 				fields: fields{
// 					client: c,
// 				},
// 				wantErr: true,
// 			}
// 		},
// 		"returns validate error when the returned job status is empty": func(*testing.T) test {
// 			c := &mock.Client{
// 				ListStatusesFunc: func(ctx context.Context, owner, repo, ref string, opts *github.ListOptions) ([]*github.RepoStatus, *github.Response, error) {
// 					return []*github.RepoStatus{}, nil, nil
// 				},
// 			}
// 			return test{
// 				fields: fields{
// 					client: c,
// 				},
// 				wantErr: true,
// 			}
// 		},
// 		"returns validate error when the success job count is invalid": func(*testing.T) test {
// 			targetJobName := "target-job"
//
// 			c := &mock.Client{
// 				ListStatusesFunc: func(ctx context.Context, owner, repo, ref string, opts *github.ListOptions) ([]*github.RepoStatus, *github.Response, error) {
// 					return []*github.RepoStatus{
// 						{
// 							Context: stringPtr("job-01"),
// 							State:   stringPtr(successState),
// 						},
// 						{
// 							Context: stringPtr("job-02"),
// 							State:   stringPtr(errorState),
// 						},
// 						{
// 							Context: stringPtr(targetJobName),
// 							State:   stringPtr(pendingState),
// 						},
// 					}, nil, nil
// 				},
// 			}
// 			return test{
// 				fields: fields{
// 					client:        c,
// 					targetJobName: targetJobName,
// 				},
// 				wantErr: true,
// 			}
// 		},
// 		"returns validate error when the statuses contains context or state is nil and the success job count is invalid": func(*testing.T) test {
// 			targetJobName := "target-job"
//
// 			c := &mock.Client{
// 				ListStatusesFunc: func(ctx context.Context, owner, repo, ref string, opts *github.ListOptions) ([]*github.RepoStatus, *github.Response, error) {
// 					return []*github.RepoStatus{
// 						{
// 							State: stringPtr(successState),
// 						},
// 						{
// 							Context: stringPtr("job-02"),
// 						},
// 						{
// 							Context: stringPtr(targetJobName),
// 							State:   stringPtr(pendingState),
// 						},
// 					}, nil, nil
// 				},
// 			}
// 			return test{
// 				fields: fields{
// 					client:        c,
// 					targetJobName: targetJobName,
// 				},
// 				wantErr: true,
// 			}
// 		},
// 	}
// 	for name, fn := range tests {
// 		t.Run(name, func(t *testing.T) {
// 			tt := fn(t)
// 			sv := &statusValidator{
// 				token:         tt.fields.token,
// 				repo:          tt.fields.repo,
// 				owner:         tt.fields.owner,
// 				ref:           tt.fields.ref,
// 				targetJobName: tt.fields.targetJobName,
// 				client:        tt.fields.client,
// 			}
// 			if err := sv.Validate(tt.ctx); (err != nil) != tt.wantErr {
// 				t.Errorf("statusValidator.Validate() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }
