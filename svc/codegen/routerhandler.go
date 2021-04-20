package codegen

import (
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"text/template"
)

var routerHandlerTmpl = `package router

import (
	"context"
	"encoding/json"
	{{.Ic.Package.Name}} {{.IcPath}}
	"example/user-svc/vo"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/ddl/query"
	"net/http"
)

var userService = service.NewUserService()

func handleError(w http.ResponseWriter, err error, status ...int) {
	logrus.Errorln(fmt.Sprintf("%+v", err))
	if len(status) > 0 {
		w.WriteHeader(status[0])
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
	if err = json.NewEncoder(w).Encode(vo.Ret{
		Code: 1,
		Data: nil,
		Msg:  err.Error(),
	}); err != nil {
		panic(err)
	}
}

func postPageUsersHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var p vo.PageQuery
	if err = json.NewDecoder(r.Body).Decode(&p); err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var pr query.PageRet

	if pr, err = userService.PostPageUsers(context.Background(), p); err != nil {
		handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(vo.Ret{
		Code: 0,
		Data: pr,
		Msg:  "",
	}); err != nil {
		handleError(w, err)
	}
}
`

func GenRouterRouter(dir string) {
	var (
		err        error
		routerfile string
		f          *os.File
		tpl        *template.Template
		routerDir  string
	)
	routerDir = filepath.Join(dir, "router")
	if err = os.MkdirAll(routerDir, os.ModePerm); err != nil {
		panic(err)
	}

	routerfile = filepath.Join(routerDir, "router.go")
	if _, err = os.Stat(routerfile); os.IsNotExist(err) {
		if f, err = os.Create(routerfile); err != nil {
			panic(err)
		}
		defer f.Close()

		if tpl, err = template.New("router.go.tmpl").Parse(routerHandlerTmpl); err != nil {
			panic(err)
		}
		if err = tpl.Execute(f, nil); err != nil {
			panic(err)
		}
	} else {
		logrus.Warnf("file %s already exists", routerfile)
	}
}
