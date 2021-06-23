package svc

import (
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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
	Dir          string
	Handler      bool
	Client       string
	Omitempty    bool
	Doc          bool
	Jsonattrcase string
}

func validateDataType(dir string) {
	astutils.BuildInterfaceCollector(filepath.Join(dir, "svc.go"), codegen.ExprStringP)
	vodir := filepath.Join(dir, "vo")
	var files []string
	err := filepath.Walk(vodir, astutils.Visit(&files))
	if err != nil {
		logrus.Panicln(err)
	}
	for _, file := range files {
		astutils.BuildStructCollector(file, codegen.ExprStringP)
	}
}

func (receiver Svc) Http() {
	var (
		err error
		dir string
	)
	dir = receiver.Dir
	if stringutils.IsEmpty(dir) {
		if dir, err = os.Getwd(); err != nil {
			panic(err)
		}
	}
	if receiver.Doc {
		validateDataType(dir)
	}

	ic := astutils.BuildInterfaceCollector(filepath.Join(dir, "svc.go"), astutils.ExprString)
	validateRestApi(ic)

	codegen.GenConfig(dir)
	codegen.GenDotenv(dir)
	codegen.GenDb(dir)
	codegen.GenHttpMiddleware(dir)

	codegen.GenMain(dir, ic)
	codegen.GenHttpHandler(dir, ic)
	if receiver.Handler {
		var caseconvertor func(string) string
		switch receiver.Jsonattrcase {
		case "snake":
			caseconvertor = strcase.ToSnake
		default:
			caseconvertor = strcase.ToLowerCamel
		}
		codegen.GenHttpHandlerImplWithImpl(dir, ic, receiver.Omitempty, caseconvertor)
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
	if receiver.Doc {
		codegen.GenDoc(dir, ic)
	}
}

// CheckIc is checking whether parameter types in each of service interface methods valid or not
// Only support at most one golang non-built-in type as parameter in a service interface method
// because go-doudou cannot put more than one parameter into request body except *multipart.FileHeader.
// If there are *multipart.FileHeader parameters, go-doudou will assume you want a multipart/form-data api
// Support struct, map[string]ANY, built-in type and corresponding slice only
func validateRestApi(ic astutils.InterfaceCollector) {
	if len(ic.Interfaces) == 0 {
		panic(errors.New("no service interface found"))
	}
	svcInter := ic.Interfaces[0]
	for _, method := range svcInter.Methods {
		// Append *multipart.FileHeader value to nonBasicTypes only once at most as multipart/form-data support multiple fields as file type
		var nonBasicTypes []string
		cpmap := make(map[string]int)
		for _, param := range method.Params {
			if param.Type == "context.Context" {
				continue
			}
			if !codegen.IsBuiltin(param) {
				ptype := param.Type
				if strings.HasPrefix(ptype, "[") || strings.HasPrefix(ptype, "*[") {
					elem := ptype[strings.Index(ptype, "]")+1:]
					if elem == "*multipart.FileHeader" {
						if _, exists := cpmap[elem]; !exists {
							cpmap[elem]++
							nonBasicTypes = append(nonBasicTypes, elem)
						}
						continue
					}
				}
				if ptype == "*multipart.FileHeader" {
					if _, exists := cpmap[ptype]; !exists {
						cpmap[ptype]++
						nonBasicTypes = append(nonBasicTypes, ptype)
					}
					continue
				}
				nonBasicTypes = append(nonBasicTypes, param.Type)
			}
		}
		if len(nonBasicTypes) > 1 {
			panic("Too many golang non-built-in type parameters, can't decide which one should be put into request body!")
		}
	}
}

func (receiver Svc) Init() {
	codegen.InitSvc(receiver.Dir)
}
