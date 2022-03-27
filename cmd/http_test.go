package cmd_test

import (
	"github.com/unionj-cloud/go-doudou/cmd"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestHttpCmd(t *testing.T) {
	dir := filepath.Join(testDir, "testsvc")
	_ = os.Chdir(dir)
	// go-doudou svc http --handler -c go -o
	_, _, err := ExecuteCommandC(cmd.GetRootCmd(), []string{"svc", "http", "--handler", "-c", "go", "-o"}...)
	if err != nil {
		t.Fatal(err)
	}
	expect := `package httpsrv

import (
	"context"
	"encoding/json"
	"net/http"
	service "testsvc"
	"testsvc/vo"

	"github.com/pkg/errors"
)

type TestsvcHandlerImpl struct {
	testsvc service.Testsvc
}

func (receiver *TestsvcHandlerImpl) PageUsers(_writer http.ResponseWriter, _req *http.Request) {
	var (
		ctx   context.Context
		query vo.PageQuery
		code  int
		data  vo.PageRet
		err   error
	)
	ctx = _req.Context()
	if err := json.NewDecoder(_req.Body).Decode(&query); err != nil {
		http.Error(_writer, err.Error(), http.StatusBadRequest)
		return
	}
	defer _req.Body.Close()
	code, data, err = receiver.testsvc.PageUsers(
		ctx,
		query,
	)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			http.Error(_writer, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(_writer, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if err := json.NewEncoder(_writer).Encode(struct {
		Code int        ` + "`" + `json:"code,omitempty"` + "`" + `
		Data vo.PageRet ` + "`" + `json:"data,omitempty"` + "`" + `
	}{
		Code: code,
		Data: data,
	}); err != nil {
		http.Error(_writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func NewTestsvcHandler(testsvc service.Testsvc) TestsvcHandler {
	return &TestsvcHandlerImpl{
		testsvc,
	}
}
`
	handlerimplfile := filepath.Join(dir, "transport", "httpsrv", "handlerimpl.go")
	f, err := os.Open(handlerimplfile)
	if err != nil {
		t.Fatal(err)
	}
	content, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != expect {
		t.Errorf("want %s, got %s\n", expect, string(content))
	}
}
