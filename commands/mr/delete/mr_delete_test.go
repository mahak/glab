package delete

import (
	"fmt"
	"github.com/acarl005/stripansi"
	"github.com/profclems/glab/commands/cmdtest"
	"github.com/profclems/glab/pkg/api"
	"github.com/stretchr/testify/assert"
	"github.com/xanzy/go-gitlab"
	"strings"
	"testing"
)

// TODO: test by mocking the appropriate api function
func TestMain(m *testing.M) {
	cmdtest.InitTest(m, "mr_delete_test")
}

func Test_deleteMergeRequest(t *testing.T) {
	t.Parallel()
	oldDeleteMR := api.DeleteMR

	api.DeleteMR = func(client *gitlab.Client, projectID interface{}, mrID int) error {
		if projectID == "" || projectID == "WRONG_REPO" || projectID == "expected_err" || mrID == 0 {
			return fmt.Errorf("error expected")
		}
		return nil
	}

	tests := []struct {
		name       string
		args       []string
		wantErr    bool
		errMsg     string
		assertFunc func(t *testing.T, out string)
	}{
		{
			name:    "delete",
			args:    []string{"0"},
			wantErr: true,

			assertFunc: func(t *testing.T, out string) {
				assert.Contains(t, out, "error expected")
			},
		},
		{
			name:    "id exists",
			args:    []string{"1"},
			wantErr: false,
			assertFunc: func(t *testing.T, out string) {
				assert.Contains(t, out, "- Deleting Merge Request !1\n")
				assert.Contains(t, out, "✔ Merge request !1 deleted\n")
			},
		},
		{
			name:    "delete on different repo",
			args:    []string{"1", "-R", "profclems/glab"},
			wantErr: false,
			assertFunc: func(t *testing.T, out string) {
				assert.Contains(t, out, "- Deleting Merge Request !1\n")
				assert.Contains(t, out, "✔ Merge request !1 deleted\n")
			},
		},
		{
			name:    "delete no args",
			wantErr: true,
			assertFunc: func(t *testing.T, out string) {
				assert.Contains(t, out, "accepts 1 arg(s), received 0")
			},
		},
	}

	cmd := NewCmdDelete(cmdtest.StubFactory(""))

	cmd.Flags().StringP("repo", "R", "", "")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			cli := strings.Join(tt.args, " ")
			t.Log(cli)
			output, err := cmdtest.RunCommand(cmd, cli)
			if !tt.wantErr {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
			}

			out := stripansi.Strip(output.String())
			outErr := stripansi.Strip(output.Stderr())

			tt.assertFunc(t, out)
			assert.Contains(t, outErr, tt.errMsg)
		})
	}

	api.DeleteMR = oldDeleteMR
}
