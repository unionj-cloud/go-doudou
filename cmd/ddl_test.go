package cmd

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/unionj-cloud/go-doudou/pathutils"
	"github.com/unionj-cloud/go-doudou/svc"
	"github.com/unionj-cloud/go-doudou/test"
	"io/ioutil"
	"os"
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

func TestDdlCmd(t *testing.T) {
	dir := testDir + "ddlcmd"
	receiver := svc.Svc{
		Dir: dir,
	}
	receiver.Init()
	defer os.RemoveAll(dir)
	err := os.Chdir(dir)
	if err != nil {
		t.Fatal(err)
	}
	logger := logrus.New()
	var terminateContainer func() // variable to store function to terminate container
	var host string
	var port int
	terminateContainer, host, port, err = test.SetupMySQLContainer(logger, pathutils.Abs("../test/sql"), "")
	defer terminateContainer() // make sure container will be terminated at the end
	if err != nil {
		logger.Error("failed to setup MySQL container")
		t.Fatal(err)
	}
	os.Setenv("DB_HOST", host)
	os.Setenv("DB_PORT", fmt.Sprint(port))
	os.Setenv("DB_USER", "root")
	os.Setenv("DB_PASSWD", "1234")
	os.Setenv("DB_SCHEMA", "test")
	os.Setenv("DB_CHARSET", "utf8mb4")
	// go-doudou ddl --dao --pre=ddl_ --domain=ddl/domain --env=ddl/.env
	_, _, err = ExecuteCommandC(rootCmd, []string{"ddl", "--reverse", "--dao", "--pre=ddl_"}...)
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
	domainfile := dir + "/domain/user.go"
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
