package svc

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/unionj-cloud/go-doudou/astutils"
	"github.com/unionj-cloud/go-doudou/stringutils"
	"github.com/unionj-cloud/go-doudou/svc/internal/codegen"
	"os"
	"path/filepath"
	"strings"
)

type SvcCmd interface {
	Init()
	Http()
}

type Svc struct {
	Dir       string
	Handler   bool
	Client    string
	Omitempty bool
}

func (receiver Svc) Http() {
	var (
		err     error
		svcfile string
		dir     string
	)
	dir = receiver.Dir
	if stringutils.IsEmpty(dir) {
		if dir, err = os.Getwd(); err != nil {
			panic(err)
		}
	}

	codegen.GenConfig(dir)
	codegen.GenDotenv(dir)
	codegen.GenDb(dir)
	codegen.GenHttpMiddleware(dir)

	svcfile = filepath.Join(dir, "svc.go")
	ic := codegen.BuildIc(svcfile)

	checkIc(ic)

	if len(ic.Interfaces) > 0 {
		codegen.GenMain(dir, ic)
		codegen.GenHttpHandler(dir, ic)
		if receiver.Handler {
			codegen.GenHttpHandlerImplWithImpl(dir, ic, receiver.Omitempty)
		} else {
			codegen.GenHttpHandlerImpl(dir, ic)
		}
		if stringutils.IsNotEmpty(receiver.Client) {
			switch receiver.Client {
			case "go":
				codegen.GenGoClient(dir, ic)
			}
		}
		codegen.GenSvcImpl(dir, ic)
		codegen.GenDoc(dir, ic)
	}
}

func checkIc(ic astutils.InterfaceCollector) {
	if len(ic.Interfaces) == 0 {
		panic(errors.New("no service interface found!"))
	}
	svcInter := ic.Interfaces[0]
	for _, method := range svcInter.Methods {
		var complexParams []string
		// Append *multipart.FileHeader value to complexTypes only once at most as multipart/form-data support multiple fields as file type
		var complexTypes []string
		cpmap := make(map[string]int)
		for _, param := range method.Params {
			if param.Type == "context.Context" {
				continue
			}
			if param.Type == "chan" || param.Type == "func" || stringutils.IsEmpty(param.Type) {
				panic(errors.New(fmt.Sprintf("Error occur at %s %s. Support struct, map, string, bool, numberic and corresponding slice type only!", param.Name, param.Type)))
			}
			if !codegen.IsSimple(param) {
				complexParams = append(complexParams, param.Name)
				ptype := param.Type
				if strings.HasPrefix(ptype, "[") || strings.HasPrefix(ptype, "*[") {
					elem := ptype[strings.Index(ptype, "]")+1:]
					if elem == "*multipart.FileHeader" {
						if _, exists := cpmap[elem]; !exists {
							cpmap[elem]++
							complexTypes = append(complexTypes, elem)
						}
						continue
					}
				}
				if ptype == "*multipart.FileHeader" {
					if _, exists := cpmap[ptype]; !exists {
						cpmap[ptype]++
						complexTypes = append(complexTypes, ptype)
					}
					continue
				}
				complexTypes = append(complexTypes, param.Type)
			}
		}
		if len(complexTypes) > 1 {
			panic("Too many struct/map/*multipart.FileHeader parameters, can't decide which one should be put into request body!")
		}
	}
}

func (receiver Svc) Init() {
	codegen.InitSvc(receiver.Dir)
}
