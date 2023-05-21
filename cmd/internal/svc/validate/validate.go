package validate

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/svc/codegen"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/astutils"
	v3helper "github.com/unionj-cloud/go-doudou/v2/toolkit/openapi/v3"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func ValidateDataType(dir string) {
	astutils.BuildInterfaceCollector(filepath.Join(dir, "svc.go"), codegen.ExprStringP)
	vodir := filepath.Join(dir, "vo")
	var files []string
	if _, err := os.Stat(vodir); !os.IsNotExist(err) {
		_ = filepath.Walk(vodir, astutils.Visit(&files))
		for _, file := range files {
			astutils.BuildStructCollector(file, codegen.ExprStringP)
		}
	}
	dtodir := filepath.Join(dir, "dto")
	if _, err := os.Stat(dtodir); !os.IsNotExist(err) {
		files = nil
		_ = filepath.Walk(dtodir, astutils.Visit(&files))
		for _, file := range files {
			astutils.BuildStructCollector(file, codegen.ExprStringP)
		}
	}
}

// ValidateRestApi is checking whether parameter types in each of service interface methods valid or not
// Only support at most one golang non-built-in type as parameter in a service interface method
// because go-doudou cannot put more than one parameter into request body except v3.FileModel.
// If there are v3.FileModel parameters, go-doudou will assume you want a multipart/form-data api
// Support struct, map[string]ANY, built-in type and corresponding slice only
// Not support anonymous struct as parameter
func ValidateRestApi(dir string, ic astutils.InterfaceCollector) {
	if len(ic.Interfaces) == 0 {
		panic(errors.New("no service interface found"))
	}
	if len(v3helper.SchemaNames) == 0 && len(v3helper.Enums) == 0 {
		codegen.ParseDto(dir, "vo")
		codegen.ParseDto(dir, "dto")
	}
	svcInter := ic.Interfaces[0]
	re := regexp.MustCompile(`anonystruct«(.*)»`)
	for _, method := range svcInter.Methods {
		nonBasicTypes := getNonBasicTypes(method.Params)
		if len(nonBasicTypes) > 1 {
			panic(fmt.Sprintf("Too many golang non-builtin type parameters in method %s, can't decide which one should be put into request body!", method))
		}
		for _, param := range method.Results {
			if re.MatchString(param.Type) {
				panic("not support anonymous struct as parameter")
			}
		}
	}
}

func getNonBasicTypes(params []astutils.FieldMeta) []string {
	var nonBasicTypes []string
	cpmap := make(map[string]int)
	re := regexp.MustCompile(`anonystruct«(.*)»`)
	for _, param := range params {
		if param.Type == "context.Context" {
			continue
		}
		if re.MatchString(param.Type) {
			panic("not support anonymous struct as parameter")
		}
		if !v3helper.IsBuiltin(param) {
			ptype := param.Type
			if strings.HasPrefix(ptype, "[") || strings.HasPrefix(ptype, "*[") {
				elem := ptype[strings.Index(ptype, "]")+1:]
				if elem == "*v3.FileModel" || elem == "v3.FileModel" || elem == "*multipart.FileHeader" {
					elem = "file"
					if _, exists := cpmap[elem]; !exists {
						cpmap[elem]++
						nonBasicTypes = append(nonBasicTypes, elem)
					}
					continue
				}
			}
			if ptype == "*v3.FileModel" || ptype == "v3.FileModel" || ptype == "*multipart.FileHeader" {
				ptype = "file"
				if _, exists := cpmap[ptype]; !exists {
					cpmap[ptype]++
					nonBasicTypes = append(nonBasicTypes, ptype)
				}
				continue
			}
			nonBasicTypes = append(nonBasicTypes, param.Type)
		}
	}
	return nonBasicTypes
}
