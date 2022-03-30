package svc

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/radovskyb/watcher"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/cmd/internal/astutils"
	"github.com/unionj-cloud/go-doudou/cmd/internal/executils"
	v3helper "github.com/unionj-cloud/go-doudou/cmd/internal/openapi/v3"
	"github.com/unionj-cloud/go-doudou/cmd/internal/openapi/v3/codegen/client"
	"github.com/unionj-cloud/go-doudou/cmd/internal/svc/codegen"
	"github.com/unionj-cloud/go-doudou/toolkit/constants"
	"github.com/unionj-cloud/go-doudou/toolkit/stringutils"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"time"
)

const (
	split = iota
	nosplit
)

// Svc wraps all config properties for commands
type Svc struct {
	// dir is project root path
	dir string
	// Handler indicates whether generate default http handler implementation code or not
	Handler bool
	// Client is client language name
	Client bool
	// Omitempty indicates whether omit empty when marshal structs to json
	Omitempty bool
	// Doc indicates whether generate OpenAPI 3.0 json doc file
	Doc bool
	// Jsonattrcase is attribute case converter name when marshal structs to json
	Jsonattrcase string
	// DocPath is OpenAPI 3.0 json doc file path used for generating client code
	DocPath string
	// Env is service base url environment variable name used for generating client code
	Env string
	// ClientPkg is client package name
	ClientPkg string

	cmd        *exec.Cmd
	restartSig chan int

	// for being compatible with legacy code purpose only
	RoutePatternStrategy int

	runner executils.Runner
	w      *watcher.Watcher

	// ModName is go module name
	ModName string

	// ImagePrefix is image name prefix string used for building and pushing docker image
	ImagePrefix string
}

func ValidateDataType(dir string) {
	astutils.BuildInterfaceCollector(filepath.Join(dir, "svc.go"), codegen.ExprStringP)
	vodir := filepath.Join(dir, "vo")
	var files []string
	_ = filepath.Walk(vodir, astutils.Visit(&files))
	for _, file := range files {
		astutils.BuildStructCollector(file, codegen.ExprStringP)
	}
}

func (receiver *Svc) SetWatcher(w *watcher.Watcher) {
	receiver.w = w
}

func (receiver *Svc) GetWatcher() *watcher.Watcher {
	return receiver.w
}

func (receiver *Svc) GetDir() string {
	return receiver.dir
}

// Http generates main function, config files, db connection function, http routes, http handlers, service interface and service implementation
// from the result of ast parsing svc.go file in the project root. It may panic if validation failed
func (receiver Svc) Http() {
	dir := receiver.dir
	ValidateDataType(dir)

	ic := astutils.BuildInterfaceCollector(filepath.Join(dir, "svc.go"), astutils.ExprString)
	ValidateRestApi(ic)

	codegen.GenConfig(dir)
	codegen.GenDb(dir)
	codegen.GenHttpMiddleware(dir)

	codegen.GenMain(dir, ic)
	codegen.GenHttpHandler(dir, ic, receiver.RoutePatternStrategy)
	var caseconvertor func(string) string
	switch receiver.Jsonattrcase {
	case "snake":
		caseconvertor = strcase.ToSnake
	default:
		caseconvertor = strcase.ToLowerCamel
	}
	codegen.ParseVo(dir, ic)
	if receiver.Handler {
		codegen.GenHttpHandlerImplWithImpl(dir, ic, receiver.Omitempty, caseconvertor)
	} else {
		codegen.GenHttpHandlerImpl(dir, ic)
	}
	if receiver.Client {
		codegen.GenGoIClient(dir, ic)
		codegen.GenGoClient(dir, ic, receiver.Env, receiver.RoutePatternStrategy, caseconvertor)
		codegen.GenGoClientProxy(dir, ic)
	}
	codegen.GenSvcImpl(dir, ic)
	if receiver.Doc {
		codegen.GenDoc(dir, ic, receiver.RoutePatternStrategy)
	}
}

// ValidateRestApi is checking whether parameter types in each of service interface methods valid or not
// Only support at most one golang non-built-in type as parameter in a service interface method
// because go-doudou cannot put more than one parameter into request body except v3.FileModel.
// If there are v3.FileModel parameters, go-doudou will assume you want a multipart/form-data api
// Support struct, map[string]ANY, built-in type and corresponding slice only
// Not support anonymous struct as parameter
func ValidateRestApi(ic astutils.InterfaceCollector) {
	if len(ic.Interfaces) == 0 {
		panic(errors.New("no service interface found"))
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

// Init inits a project
func (receiver Svc) Init() {
	codegen.InitProj(receiver.dir, receiver.ModName, receiver.runner)
}

type SvcOption func(svc *Svc)

func WithRunner(runner executils.Runner) SvcOption {
	return func(svc *Svc) {
		svc.runner = runner
	}
}

// NewSvc new Svc instance
func NewSvc(dir string, opts ...SvcOption) Svc {
	ret := Svc{
		dir:        dir,
		runner:     executils.CmdRunner{},
		restartSig: make(chan int),
	}
	for _, opt := range opts {
		opt(&ret)
	}
	return ret
}

// Push executes go mod vendor command first, then build docker image and push to remote image repository
// It also generates deployment kind(for monolithic) and statefulset kind(for microservice) yaml files for kubernetes deploy, if these files already exist,
// it will only change the image version in each file, so you can edit these files manually to fit your need.
func (receiver Svc) Push(repo string) {
	ic := astutils.BuildInterfaceCollector(filepath.Join(receiver.dir, "svc.go"), astutils.ExprString)
	err := receiver.runner.Run("go", "mod", "vendor")
	if err != nil {
		panic(err)
	}

	svcname := strings.ToLower(ic.Interfaces[0].Name)
	imageName := fmt.Sprintf("%s%s", receiver.ImagePrefix, svcname)
	loginUser, _ := user.Current()
	if loginUser != nil {
		err = receiver.runner.Run("docker", "build", "--build-arg", fmt.Sprintf("user=%s", loginUser.Username), "-t", imageName, ".")
		if err != nil {
			panic(err)
		}
	} else {
		err = receiver.runner.Run("docker", "build", "-t", imageName, ".")
		if err != nil {
			panic(err)
		}
	}

	if stringutils.IsEmpty(repo) {
		logrus.Warnln("no private docker image repository address provided")
		return
	}
	image := fmt.Sprintf("%s/%s:%s", repo, imageName, fmt.Sprintf("v%s", time.Now().Local().Format(constants.FORMAT11)))
	err = receiver.runner.Run("docker", "tag", imageName, image)
	if err != nil {
		panic(err)
	}
	err = receiver.runner.Run("docker", "push", image)
	if err != nil {
		panic(err)
	}
	logrus.Infof("image %s has been pushed successfully\n", image)

	codegen.GenK8sDeployment(receiver.dir, svcname, image)
	codegen.GenK8sStatefulset(receiver.dir, svcname, image)
	logrus.Infof("k8s yaml has been created/updated successfully. execute command 'go-doudou svc deploy' to deploy service %s to k8s cluster\n", svcname)
}

// Deploy deploys project to kubernetes. If k8sfile flag not set, it will be deployed as statefulset kind using statefulset.yaml file in the project root,
// so if you want to deploy a monolithic project, please set k8sfile flag.
func (receiver Svc) Deploy(k8sfile string) {
	ic := astutils.BuildInterfaceCollector(filepath.Join(receiver.dir, "svc.go"), astutils.ExprString)
	svcname := strings.ToLower(ic.Interfaces[0].Name)
	if stringutils.IsEmpty(k8sfile) {
		k8sfile = filepath.Join(receiver.dir, svcname+"_statefulset.yaml")
	}
	logrus.Infof("Execute command: kubectl apply -f %s\n", k8sfile)
	if err := receiver.runner.Run("kubectl", "apply", "-f", k8sfile); err != nil {
		panic(err)
	}
}

// Shutdown stops and removes the project from kubernetes. If k8sfile flag not set, it will use statefulset.yaml file in the project root,
// so if you had already set k8sfile flag when you deploy the project, you should set the same k8sfile flag.
func (receiver Svc) Shutdown(k8sfile string) {
	ic := astutils.BuildInterfaceCollector(filepath.Join(receiver.dir, "svc.go"), astutils.ExprString)
	svcname := strings.ToLower(ic.Interfaces[0].Name)
	if stringutils.IsEmpty(k8sfile) {
		k8sfile = filepath.Join(receiver.dir, svcname+"_statefulset.yaml")
	}
	logrus.Infof("Execute command: kubectl delete -f %s\n", k8sfile)
	if err := receiver.runner.Run("kubectl", "delete", "-f", k8sfile); err != nil {
		panic(err)
	}
}

// GenClient generates http client code from OpenAPI3.0 description json file, only support Golang currently.
func (receiver Svc) GenClient() {
	docpath := receiver.DocPath
	if stringutils.IsEmpty(docpath) {
		matches, _ := filepath.Glob(filepath.Join(receiver.dir, "*_openapi3.json"))
		if len(matches) > 0 {
			docpath = matches[0]
		}
	}
	if stringutils.IsEmpty(docpath) {
		panic("openapi 3.0 spec json file path is empty")
	}
	client.GenGoClient(receiver.dir, docpath, receiver.Omitempty, receiver.Env, receiver.ClientPkg)
}

func (receiver *Svc) DoRun() {
	err := receiver.runner.Run("go", "build", filepath.FromSlash("cmd/main.go"))
	if err != nil {
		panic(err)
	}
	start, err := receiver.runner.Start(filepath.FromSlash("./main"))
	if err != nil {
		panic(err)
	}
	receiver.cmd = start
}

//func terminateWinProc(pid int) error {
//	dll, err := windows.LoadDLL("kernel32.dll")
//	if err != nil {
//		return err
//	}
//	defer dll.Release()
//	f, err := dll.FindProc("GenerateConsoleCtrlEvent")
//	if err != nil {
//		return err
//	}
//	r1, _, err := f.Call(windows.CTRL_BREAK_EVENT, uintptr(pid))
//	if r1 == 0 {
//		return err
//	}
//	return nil
//}

// TODO there is a bug here on windows
func (receiver *Svc) DoRestart() {
	//if runtime.GOOS == "windows" {
	//	if err := terminateWinProc(receiver.Cmd.Process.Pid); err != nil {
	//		panic(err)
	//	}
	//} else {
	//	if err := receiver.Cmd.Process.Signal(syscall.SIGINT); err != nil {
	//		panic(err)
	//	}
	//}
	if receiver.cmd != nil {
		if err := receiver.cmd.Process.Signal(syscall.SIGINT); err != nil {
			panic(err)
		}
	}
	receiver.DoRun()
}

func (receiver Svc) DoWatch() {
	w := receiver.w
	// SetMaxEvents to 1 to allow at most 1 event's to be received
	// on the Event channel per watching cycle.
	//
	// If SetMaxEvents is not set, the default is to send all events.
	w.SetMaxEvents(1)

	// Only notify write events.
	w.FilterOps(watcher.Write)

	// Only files that match the regular expression during file listings
	// will be watched.
	r := regexp.MustCompile("\\.go$")
	w.AddFilterHook(watcher.RegexFilterHook(r, false))

	go func() {
		for {
			select {
			case event := <-w.Event:
				fmt.Println(event) // Print the event's info.
				if err := receiver.runner.Run("go", "build", "cmd/main.go"); err != nil {
					logrus.Warnln(err)
					continue
				}
				_ = os.Remove("main")
				receiver.restartSig <- 1
			case err := <-w.Error:
				logrus.Panicln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	// Watch this folder for changes.
	if err := w.AddRecursive(receiver.dir); err != nil {
		logrus.Panicln(err)
	}

	// Print a list of all of the files and folders currently
	// being watched and their paths.
	for path, f := range w.WatchedFiles() {
		logrus.Tracef("%s: %s\n", path, f.Name())
	}

	// Start the watching process - it'll check for changes every 100ms.
	if err := w.Start(time.Millisecond * 100); err != nil {
		logrus.Panicln(err)
	}
}

// Run runs the project locally. Recommend to set watch flag to enable watch mode for rapid development.
func (receiver Svc) Run(watch bool) {
	receiver.DoRun()
	if watch {
		if receiver.w == nil {
			receiver.w = watcher.New()
		}
		go receiver.DoWatch()
		for {
			select {
			case <-receiver.restartSig:
				receiver.DoRestart()
			default:
				time.Sleep(100 * time.Millisecond)
			}
		}
	} else {
		if err := receiver.cmd.Wait(); err != nil {
			panic(err)
		}
	}
}

// Upgrade upgrades go-doudou to latest release version
func (receiver Svc) Upgrade(version string) {
	fmt.Printf("go install -v github.com/unionj-cloud/go-doudou@%s\n", version)
	if err := receiver.runner.Run("go", "install", "-v", fmt.Sprintf("github.com/unionj-cloud/go-doudou@%s", version)); err != nil {
		panic(err)
	}
}
