package cmd_test

import (
	"bytes"
	"github.com/spf13/cobra"
	"github.com/unionj-cloud/go-doudou/cmd"
	"github.com/unionj-cloud/go-doudou/cmd/internal/ddl/table"
	"github.com/unionj-cloud/go-doudou/cmd/internal/svc"
	"github.com/unionj-cloud/go-doudou/toolkit/pathutils"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

var testDir string

func init() {
	testDir = pathutils.Abs("testdata")
}

func ExecuteCommandC(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	c, err = root.ExecuteC()

	return c, buf.String(), err
}

// NewMockSvc new Svc instance for unit test purpose
func NewMockSvc(dir string, opts ...svc.SvcOption) svc.ISvc {
	return svc.NewSvc(dir, svc.WithRunner(mockRunner{}))
}

type mockRunner struct {
}

func (r mockRunner) Output(command string, args ...string) ([]byte, error) {
	return []byte("go version go1.17.8 darwin/amd64"), nil
}

func (r mockRunner) Run(command string, args ...string) error {
	cs := []string{"-test.run=TestHelperProcess", "--"}
	cs = append(cs, args...)
	c := exec.Command(os.Args[0], cs...)
	c.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		panic(err)
	}
	return nil
}

func (r mockRunner) Start(command string, args ...string) (*exec.Cmd, error) {
	cs := []string{"-test.run=TestHelperProcess", "--"}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		panic(err)
	}
	return cmd, nil
}

func TestDdlCmd(t *testing.T) {
	dir := testDir + "/ddlcmd"
	receiver := NewMockSvc(dir)
	receiver.Init()
	defer os.RemoveAll(dir)
	err := os.Chdir(dir)
	if err != nil {
		t.Fatal(err)
	}
	terminator, db, err := table.Setup()
	if err != nil {
		panic(err)
	}
	defer terminator()
	defer db.Close()
	// go-doudou ddl --dao --pre=ddl_ --domain=ddl/domain --env=ddl/.env
	_, _, err = ExecuteCommandC(cmd.GetRootCmd(), []string{"ddl", "--reverse", "--dao", "--pre=ddl_"}...)
	if err != nil {
		t.Fatal(err)
	}
	expect := `package domain

import (
	"time"
)

//dd:table
type User struct {
	Id        int        ` + "`" + `dd:"pk;auto;type:int"` + "`" + `
	Name      string     ` + "`" + `dd:"type:varchar(255);default:'jack';index:name_phone_idx,2,asc"` + "`" + `
	Phone     string     ` + "`" + `dd:"type:varchar(255);default:'13552053960';extra:comment 'mobile phone';index:name_phone_idx,1,asc"` + "`" + `
	Age       int        ` + "`" + `dd:"type:int;index:age_idx,1,asc"` + "`" + `
	No        int        ` + "`" + `dd:"type:int;unique:no_idx,1,asc"` + "`" + `
	School    *string    ` + "`" + `dd:"type:varchar(255);default:'harvard';extra:comment 'school'"` + "`" + `
	IsStudent int8       ` + "`" + `dd:"type:tinyint"` + "`" + `
	CreateAt  *time.Time ` + "`" + `dd:"type:datetime;default:CURRENT_TIMESTAMP"` + "`" + `
	UpdateAt  *time.Time ` + "`" + `dd:"type:datetime;default:CURRENT_TIMESTAMP;extra:on update CURRENT_TIMESTAMP"` + "`" + `
	DeleteAt  *time.Time ` + "`" + `dd:"type:datetime"` + "`" + `
}
`
	domainfile := filepath.Join(dir, "domain", "user.go")
	f, err := os.Open(domainfile)
	if err != nil {
		t.Fatal(err)
	}
	content, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != expect {
		t.Errorf("want %s, go %s\n", expect, string(content))
	}
}
