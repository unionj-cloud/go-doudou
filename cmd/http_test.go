package cmd_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/unionj-cloud/go-doudou/v2/cmd"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestHttpCmd(t *testing.T) {
	dir := filepath.Join(testDir, "testsvc")
	_ = os.Chdir(dir)
	// go-doudou svc http -c -o
	_, _, err := ExecuteCommandC(cmd.GetRootCmd(), []string{"svc", "http", "-c", "-o"}...)
	if err != nil {
		t.Fatal(err)
	}
	expect := `package httpsrv

import (
	"context"
	"net/http"
	service "testsvc"
	"testsvc/vo"

	"encoding/json"

	"github.com/pkg/errors"
	"github.com/unionj-cloud/go-doudou/v2/framework/rest"
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
	if _req.Body == nil {
		http.Error(_writer, "missing request body", http.StatusBadRequest)
		return
	} else {
		if _err := json.NewDecoder(_req.Body).Decode(&query); _err != nil {
			http.Error(_writer, _err.Error(), http.StatusBadRequest)
			return
		} else {
			if _err := rest.ValidateStruct(query); _err != nil {
				http.Error(_writer, _err.Error(), http.StatusBadRequest)
				return
			}
		}
	}
	code, data, err = receiver.testsvc.PageUsers(
		ctx,
		query,
	)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			http.Error(_writer, err.Error(), http.StatusBadRequest)
		} else if _err, ok := err.(*rest.BizError); ok {
			http.Error(_writer, _err.Error(), _err.StatusCode)
		} else {
			http.Error(_writer, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if _err := json.NewEncoder(_writer).Encode(struct {
		Code int        ` + "`" + `json:"code,omitempty"` + "`" + `
		Data vo.PageRet ` + "`" + `json:"data,omitempty"` + "`" + `
	}{
		Code: code,
		Data: data,
	}); _err != nil {
		http.Error(_writer, _err.Error(), http.StatusInternalServerError)
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
	assert.Equal(t, expect, string(content))
}
