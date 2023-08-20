package codegen

import (
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/svc/parser"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/astutils"
)

func genHttp(dir string, ic astutils.InterfaceCollector, caseConvertor func(string) string) {
	GenConfig(dir)
	GenHttpMiddleware(dir)
	GenHttpHandler(dir, ic, 0)
	GenHttpHandlerImpl(dir, ic, GenHttpHandlerImplConfig{
		CaseConvertor: caseConvertor,
	})
	GenSvcImpl(dir, ic)
	parser.GenDoc(dir, ic, parser.GenDocConfig{})
}
