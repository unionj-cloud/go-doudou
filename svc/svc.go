package svc

import (
	"context"
	"fmt"
	"github.com/Jeffail/gabs/v2"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/astutils"
	"github.com/unionj-cloud/go-doudou/constants"
	"github.com/unionj-cloud/go-doudou/esutils"
	"github.com/unionj-cloud/go-doudou/logutils"
	"github.com/unionj-cloud/go-doudou/stringutils"
	"github.com/unionj-cloud/go-doudou/svc/internal/codegen"
	"github.com/unionj-cloud/go-doudou/test"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
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

	DocPath   string
	Es        *esutils.Es
	ImageRepo string

	K8sfile string
	N       int
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
	dir := receiver.Dir
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

func (receiver Svc) Push() {
	ic := astutils.BuildInterfaceCollector(filepath.Join(receiver.Dir, "svc.go"), astutils.ExprString)
	validateRestApi(ic)
	svcname := strings.ToLower(ic.Interfaces[0].Name)

	cmd := exec.Command("docker", "build", "-t", svcname, ".")
	_, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	if stringutils.IsEmpty(receiver.ImageRepo) {
		logrus.Warnln("stopped as no private docker image repository address provided")
		return
	}
	image := fmt.Sprintf("%s/%s:%s", receiver.ImageRepo, svcname, fmt.Sprintf("v%s", time.Now().Local().Format(constants.FORMAT11)))
	cmd = exec.Command("docker", "tag", svcname, image)
	_, err = cmd.Output()
	if err != nil {
		panic(err)
	}

	cmd = exec.Command("docker", "push", image)
	_, err = cmd.Output()
	if err != nil {
		panic(err)
	}
	logrus.Infof("image %s has been pushed successfully\n", image)

	codegen.GenK8s(receiver.Dir, svcname, image)
	logrus.Infof("k8s yaml has been created/updated successfully. execute command 'go-doudou svc deploy' to deploy service %s to k8s cluster\n", svcname)
}

func (receiver Svc) Deploy() {
	ic := astutils.BuildInterfaceCollector(filepath.Join(receiver.Dir, "svc.go"), astutils.ExprString)
	validateRestApi(ic)
	svcname := strings.ToLower(ic.Interfaces[0].Name)
	k8sfile := receiver.K8sfile
	if stringutils.IsEmpty(k8sfile) {
		k8sfile = svcname + "_k8s.yaml"
	}
	cmd := exec.Command("kubectl", "apply", "-f", k8sfile)
	stdout, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	logrus.Infoln(string(stdout))
}

func (receiver Svc) Shutdown() {
	ic := astutils.BuildInterfaceCollector(filepath.Join(receiver.Dir, "svc.go"), astutils.ExprString)
	validateRestApi(ic)
	svcname := strings.ToLower(ic.Interfaces[0].Name)
	k8sfile := receiver.K8sfile
	if stringutils.IsEmpty(k8sfile) {
		k8sfile = svcname + "_k8s.yaml"
	}
	cmd := exec.Command("kubectl", "delete", "-f", k8sfile)
	stdout, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	logrus.Infoln(string(stdout))
}

func (receiver Svc) Scale() {
	ic := astutils.BuildInterfaceCollector(filepath.Join(receiver.Dir, "svc.go"), astutils.ExprString)
	validateRestApi(ic)
	svcname := strings.ToLower(ic.Interfaces[0].Name)
	cmd := exec.Command("kubectl", "scale", fmt.Sprintf("--replicas=%d", receiver.N), fmt.Sprintf("deployment/%s-deployment", svcname))
	stdout, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	logrus.Infoln(string(stdout))
}

func (receiver Svc) Publish() string {
	ic := astutils.BuildInterfaceCollector(filepath.Join(receiver.Dir, "svc.go"), astutils.ExprString)
	validateRestApi(ic)
	svcname := strings.ToLower(ic.Interfaces[0].Name)
	docpath := receiver.DocPath
	if stringutils.IsEmpty(docpath) {
		docpath = svcname + "_openapi3.json"
	}
	container, err := gabs.ParseJSONFile(docpath)
	if err != nil {
		panic(err)
	}
	version := container.Path("info.version").Data().(string)
	result, err := receiver.Es.SaveOrUpdate(context.Background(), struct {
		Api      string    `json:"api,omitempty"`
		CreateAt time.Time `json:"createAt,omitempty"`
		Service  string    `json:"service,omitempty"`
		Version  string    `json:"version,omitempty"`
	}{
		Api:      container.String(),
		CreateAt: time.Now().UTC(),
		Service:  svcname,
		Version:  version,
	})
	if err != nil {
		panic(err)
	}
	return result
}

func prepareTestEnvironment() (string, func()) {
	logger := logutils.NewLogger()
	var terminateContainer func() // variable to store function to terminate container
	var host string
	var port int
	var err error
	terminateContainer, host, port, err = test.SetupEs6Container(logger)
	if err != nil {
		logger.Panicln("failed to setup Elasticsearch container")
	}
	return fmt.Sprintf("http://%s:%d", host, port), terminateContainer
}
