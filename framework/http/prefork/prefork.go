package prefork

import (
	"fmt"
	logger "github.com/unionj-cloud/go-doudou/toolkit/zlogger"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/libp2p/go-reuseport"
)

const (
	envPreforkChildKey = "GDD_PREFORK_CHILD"
	envPreforkChildVal = "1"
	defaultNetwork     = "tcp4"
)

var (
	defaultLogger = Logger(log.New(os.Stderr, "", log.LstdFlags))
)

// IsChild determines if the current process is a child of Prefork
func IsChild() bool {
	return os.Getenv(envPreforkChildKey) == envPreforkChildVal
}

// Logger is used for logging formatted messages.
type Logger interface {
	// Printf must have the same semantics as log.Printf.
	Printf(format string, args ...interface{})
}

// Prefork implements fasthttp server prefork
//
// Preforks master process (with all cores) between several child processes
// increases performance significantly, because Go doesn't have to share
// and manage memory between cores
//
// WARNING: using prefork prevents the use of any global state!
// Things like in-memory caches won't work.
type Prefork struct {
	// The network must be "tcp", "tcp4" or "tcp6".
	//
	// By default is "tcp4"
	Network      string
	Addr         string
	ServeFunc    func(ln net.Listener) error
	ServeTLSFunc func(ln net.Listener, certFile, keyFile string) error
}

// New wraps the net/http server to run with preforked processes
func New(s *http.Server) *Prefork {
	return &Prefork{
		Network:      defaultNetwork,
		Addr:         s.Addr,
		ServeFunc:    s.Serve,
		ServeTLSFunc: s.ServeTLS,
	}
}

func (p *Prefork) listen(addr string) (net.Listener, error) {
	runtime.GOMAXPROCS(1)

	if p.Network == "" {
		p.Network = defaultNetwork
	}

	return reuseport.Listen(p.Network, addr)
}

func (p *Prefork) doCommand() (*exec.Cmd, error) {
	/* #nosec G204 */
	cmd := exec.Command(os.Args[0], os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("%s=%s", envPreforkChildKey, envPreforkChildVal),
	)
	return cmd, cmd.Start()
}

func (p *Prefork) prefork() (err error) {
	type procSig struct {
		pid int
		err error
	}

	goMaxProcs := runtime.GOMAXPROCS(0)
	sigCh := make(chan procSig, goMaxProcs)
	childProcs := make(map[int]*exec.Cmd)

	defer func() {
		for _, proc := range childProcs {
			_ = proc.Process.Kill()
		}
	}()

	for i := 0; i < goMaxProcs; i++ {
		var cmd *exec.Cmd
		if cmd, err = p.doCommand(); err != nil {
			logger.Error().Err(err).Msg("failed to start a child prefork process")
			return
		}

		childProcs[cmd.Process.Pid] = cmd
		go func() {
			sigCh <- procSig{cmd.Process.Pid, cmd.Wait()}
		}()
	}

	// return error if child crashes
	if err = (<-sigCh).err; err != nil {
		logger.Error().Err(err).Msg("")
	}
	return err
}

// ListenAndServe serves HTTP requests from the given TCP addr
func (p *Prefork) ListenAndServe() error {
	if IsChild() {
		ln, err := p.listen(p.Addr)
		if err != nil {
			return err
		}
		logger.Info().Msgf("Http server is listening at %v", p.Addr)
		// kill current child proc when master exits
		go watchMaster()
		defer func() {
			if e := ln.Close(); err == nil {
				err = e
			}
		}()
		if err = p.ServeFunc(ln); err != nil {
			logger.Error().Err(err).Msg("")
		}
		return err
	}

	var err error
	if err = p.prefork(); err != nil {
		logger.Error().Err(err).Msg("")
	}
	return err
}

// ListenAndServeTLS serves HTTPS requests from the given TCP addr
//
// certFile and keyFile are paths to TLS certificate and key files.
func (p *Prefork) ListenAndServeTLS(certKey, certFile string) error {
	if IsChild() {
		ln, err := p.listen(p.Addr)
		if err != nil {
			return err
		}
		logger.Info().Msgf("Http server is listening at %v", p.Addr)
		// kill current child proc when master exits
		go watchMaster()
		defer func() {
			if e := ln.Close(); err == nil {
				err = e
			}
		}()
		return p.ServeTLSFunc(ln, certFile, certKey)
	}

	return p.prefork()
}

// watchMaster watches child procs
func watchMaster() {
	if runtime.GOOS == "windows" {
		// finds parent process,
		// and waits for it to exit
		p, err := os.FindProcess(os.Getppid())
		if err == nil {
			_, _ = p.Wait()
		}
		os.Exit(1)
	}
	// if it is equal to 1 (init process ID),
	// it indicates that the master process has exited
	for range time.NewTicker(time.Millisecond * 500).C {
		if os.Getppid() == 1 {
			logger.Info().Msg("exit self as parent process exit")
			os.Exit(1)
		}
	}
}
