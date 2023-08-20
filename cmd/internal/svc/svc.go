package svc

import (
	"fmt"
	"github.com/radovskyb/watcher"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/openapi/v3/codegen/client"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/openapi/v3/codegen/server"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/svc/codegen"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/svc/codegen/database"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/svc/parser"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/svc/validate"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/assert"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/astutils"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/executils"
	v3 "github.com/unionj-cloud/go-doudou/v2/toolkit/protobuf/v3"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/stringutils"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"time"
)

//go:generate mockgen -destination ../../mock/mock_svc.go -package mock -source=./svc.go

type ISvc interface {
	SetWatcher(w *watcher.Watcher)
	GetWatcher() *watcher.Watcher
	GetDir() string
	Http()
	Init()
	Push(cfg PushConfig)
	Deploy(k8sfile string)
	Shutdown(k8sfile string)
	GenClient()
	DoRun()
	DoRestart()
	DoWatch()
	Run(watch bool)
	Upgrade(version string)
	Grpc()
}

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

	// PostmanCollectionPath is postman collection v2.1 compatible file disk path
	PostmanCollectionPath string
	// DotenvPath dotenv format config file disk path only for integration testing purpose
	DotenvPath string

	// AllowGetWithReqBody indicates whether allow get http request with request body.
	// If true, when you defined a get api with struct type parameter in svc.go file,
	// it will try to decode json format encoded request body.
	AllowGetWithReqBody bool

	DbConfig *DbConfig

	module         bool
	protoGenerator v3.ProtoGenerator

	JsonCase      string
	CaseConverter func(string) string
}

type DbConfig struct {
	Driver string
	Dsn    string
	// or schema for pg
	TablePrefix string
	TableGlob   string
	Orm         string
	Soft        string
	Grpc        bool
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

func (receiver *Svc) SetRunner(runner executils.Runner) {
	receiver.runner = runner
}

// Http generates main function, config files, db connection function, http routes, http handlers, service interface and service implementation
// from the result of ast parsing svc.go file in the project root. It may panic if validation failed
func (receiver *Svc) Http() {
	dir := receiver.dir
	parser.ParseDto(dir, "vo")
	parser.ParseDto(dir, "dto")
	validate.DataType(dir)

	ic := astutils.BuildInterfaceCollector(filepath.Join(dir, "svc.go"), astutils.ExprString)
	validate.RestApi(dir, ic)

	codegen.GenConfig(dir)
	codegen.GenHttpMiddleware(dir)

	codegen.GenMain(dir, ic)
	codegen.GenHttpHandler(dir, ic, receiver.RoutePatternStrategy)
	codegen.GenHttpHandlerImpl(dir, ic, codegen.GenHttpHandlerImplConfig{
		Omitempty:           receiver.Omitempty,
		AllowGetWithReqBody: receiver.AllowGetWithReqBody,
		CaseConvertor:       receiver.CaseConverter,
	})
	if receiver.Client {
		codegen.GenGoIClient(dir, ic)
		codegen.GenGoClient(dir, ic, codegen.GenGoClientConfig{
			Env:                  receiver.Env,
			RoutePatternStrategy: receiver.RoutePatternStrategy,
			AllowGetWithReqBody:  receiver.AllowGetWithReqBody,
			CaseConvertor:        receiver.CaseConverter,
		})
		codegen.GenGoClientProxy(dir, ic)
	}
	codegen.GenSvcImpl(dir, ic)
	parser.GenDoc(dir, ic, parser.GenDocConfig{
		RoutePatternStrategy: receiver.RoutePatternStrategy,
		AllowGetWithReqBody:  receiver.AllowGetWithReqBody,
	})
	runner := receiver.runner
	if runner == nil {
		runner = executils.CmdRunner{}
	}
	runner.Run("go", "mod", "tidy")
}

// Init inits a project
func (receiver *Svc) Init() {
	codegen.InitProj(codegen.InitProjConfig{
		Dir:            receiver.dir,
		ModName:        receiver.ModName,
		Runner:         receiver.runner,
		GenSvcGo:       receiver.DbConfig == nil,
		Module:         receiver.module,
		ProtoGenerator: receiver.protoGenerator,
		CaseConverter:  receiver.CaseConverter,
		JsonCase:       receiver.JsonCase,
	})
	// generate or overwrite svc.go file
	if receiver.DbConfig != nil {
		gen := database.GetOrmGenerator(database.OrmKind(receiver.DbConfig.Orm))
		assert.NotNil(gen, "Unknown orm kind")
		gen.Initialize(database.OrmGeneratorConfig{
			Driver:        receiver.DbConfig.Driver,
			Dsn:           receiver.DbConfig.Dsn,
			TablePrefix:   receiver.DbConfig.TablePrefix,
			TableGlob:     receiver.DbConfig.TableGlob,
			CaseConverter: receiver.CaseConverter,
			Dir:           receiver.dir,
			Soft:          receiver.DbConfig.Soft,
			Grpc:          receiver.DbConfig.Grpc,
		})
		gen.GenService()
	} else if !receiver.module {
		if stringutils.IsEmpty(receiver.DocPath) {
			matches, _ := filepath.Glob(filepath.Join(receiver.dir, "*_openapi3.json"))
			if len(matches) > 0 {
				receiver.DocPath = matches[0]
			}
		}
		if stringutils.IsNotEmpty(receiver.DocPath) {
			server.GenSvcGo(receiver.dir, receiver.DocPath)
		}
	}
}

type SvcOption func(svc *Svc)

func WithRunner(runner executils.Runner) SvcOption {
	return func(svc *Svc) {
		svc.runner = runner
	}
}

func WithModName(modName string) SvcOption {
	return func(svc *Svc) {
		svc.ModName = modName
	}
}

func WithDocPath(docfile string) SvcOption {
	return func(svc *Svc) {
		svc.DocPath = docfile
	}
}

func WithDbConfig(dbConfig *DbConfig) SvcOption {
	return func(svc *Svc) {
		svc.DbConfig = dbConfig
	}
}

func WithModule(module bool) SvcOption {
	return func(svc *Svc) {
		svc.module = module
	}
}

func WithCaseConverter(fn func(string) string) SvcOption {
	return func(svc *Svc) {
		svc.CaseConverter = fn
	}
}

func WithJsonCase(jsonCase string) SvcOption {
	return func(svc *Svc) {
		svc.JsonCase = jsonCase
	}
}

func WithProtoGenerator(protoGenerator v3.ProtoGenerator) SvcOption {
	return func(svc *Svc) {
		svc.protoGenerator = protoGenerator
	}
}

// NewSvc new Svc instance
func NewSvc(dir string, opts ...SvcOption) ISvc {
	ret := Svc{
		dir:        dir,
		runner:     executils.CmdRunner{},
		restartSig: make(chan int),
	}
	for _, opt := range opts {
		opt(&ret)
	}
	return &ret
}

type PushConfig struct {
	Repo   string
	Prefix string
	Ver    string
}

// Push executes go mod vendor command first, then build docker image and push to remote image repository
// It also generates deployment kind and statefulset kind yaml files for kubernetes deploy, if these files already exist,
// it will only change the image version in each file, so you can edit these files manually to fit your need.
func (receiver *Svc) Push(cfg PushConfig) {
	ic := astutils.BuildInterfaceCollector(filepath.Join(receiver.dir, "svc.go"), astutils.ExprString)
	svcname := strings.ToLower(ic.Interfaces[0].Name)
	imageName := fmt.Sprintf("%s%s", cfg.Prefix, svcname)
	if stringutils.IsNotEmpty(cfg.Ver) {
		imageName += ":" + cfg.Ver
	}
	loginUser, _ := user.Current()
	var err error
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

	if stringutils.IsNotEmpty(cfg.Repo) {
		remoteImageName := cfg.Repo + "/" + imageName
		err = receiver.runner.Run("docker", "tag", imageName, remoteImageName)
		if err != nil {
			panic(err)
		}
		err = receiver.runner.Run("docker", "push", remoteImageName)
		if err != nil {
			panic(err)
		}
		logrus.Infof("image %s has been pushed successfully\n", remoteImageName)
		imageName = remoteImageName
	}

	codegen.GenK8sDeployment(receiver.dir, svcname, imageName)
	codegen.GenK8sStatefulset(receiver.dir, svcname, imageName)
	logrus.Infof("k8s yaml has been created/updated successfully. execute command 'go-doudou svc deploy' to deploy service %s to k8s cluster\n", svcname)
}

// Deploy deploys project to kubernetes. If k8sfile flag not set, it will be deployed as deployment kind using *_deployment.yaml file in the project root,
func (receiver *Svc) Deploy(k8sfile string) {
	ic := astutils.BuildInterfaceCollector(filepath.Join(receiver.dir, "svc.go"), astutils.ExprString)
	svcname := strings.ToLower(ic.Interfaces[0].Name)
	if stringutils.IsEmpty(k8sfile) {
		k8sfile = filepath.Join(receiver.dir, svcname+"_deployment.yaml")
	}
	logrus.Infof("Execute command: kubectl apply -f %s\n", k8sfile)
	if err := receiver.runner.Run("kubectl", "apply", "-f", k8sfile); err != nil {
		panic(err)
	}
}

// Shutdown stops and removes the project from kubernetes. If k8sfile flag not set, it will use *_deployment.yaml file in the project root,
// so if you had already set k8sfile flag when you deploy the project, you should set the same k8sfile flag.
func (receiver *Svc) Shutdown(k8sfile string) {
	ic := astutils.BuildInterfaceCollector(filepath.Join(receiver.dir, "svc.go"), astutils.ExprString)
	svcname := strings.ToLower(ic.Interfaces[0].Name)
	if stringutils.IsEmpty(k8sfile) {
		k8sfile = filepath.Join(receiver.dir, svcname+"_deployment.yaml")
	}
	logrus.Infof("Execute command: kubectl delete -f %s\n", k8sfile)
	if err := receiver.runner.Run("kubectl", "delete", "-f", k8sfile); err != nil {
		panic(err)
	}
}

// GenClient generates http client code from OpenAPI3.0 description json file, only support Golang currently.
func (receiver *Svc) GenClient() {
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

// GenIntegrationTestingCode generates integration testing code from postman collection v2.1 compatible file
func (receiver *Svc) GenIntegrationTestingCode() {
	ic := astutils.BuildInterfaceCollector(filepath.Join(receiver.dir, "svc.go"), astutils.ExprString)
	codegen.GenHttpIntegrationTesting(receiver.dir, ic, receiver.PostmanCollectionPath, receiver.DotenvPath)
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

func (receiver *Svc) DoWatch() {
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
func (receiver *Svc) Run(watch bool) {
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
func (receiver *Svc) Upgrade(version string) {
	fmt.Printf("go install -v github.com/unionj-cloud/go-doudou/v2@%s\n", version)
	if err := receiver.runner.Run("go", "install", "-v", fmt.Sprintf("github.com/unionj-cloud/go-doudou/v2@%s", version)); err != nil {
		panic(err)
	}
}

func (receiver *Svc) Grpc() {
	dir := receiver.dir
	validate.DataType(dir)
	ic := astutils.BuildInterfaceCollector(filepath.Join(dir, "svc.go"), astutils.ExprString)
	validate.RestApi(dir, ic)
	codegen.GenConfig(dir)
	parser.ParseDtoGrpc(dir, receiver.protoGenerator, "vo")
	parser.ParseDtoGrpc(dir, receiver.protoGenerator, "dto")
	grpcSvc, protoFile := codegen.GenGrpcProto(dir, ic, receiver.protoGenerator)
	// protoc --proto_path=. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative transport/grpc/helloworld.proto
	if err := receiver.runner.Run("protoc", "--proto_path=.",
		"--go_out=.",
		"--go_opt=paths=source_relative",
		"--go-grpc_out=.",
		"--go-grpc_opt=paths=source_relative",
		protoFile); err != nil {
		panic(err)
	}
	codegen.GenSvcImplGrpc(dir, ic, grpcSvc)
	codegen.GenMainGrpc(dir, ic, grpcSvc)
	codegen.FixModGrpc(dir)
	codegen.GenMethodAnnotationStore(dir, ic)
	runner := receiver.runner
	if runner == nil {
		runner = executils.CmdRunner{}
	}
	runner.Run("go", "mod", "tidy")
}
