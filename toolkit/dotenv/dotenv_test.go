package dotenv_test

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/unionj-cloud/go-doudou/toolkit/dotenv"
	"io"
	"os"
	"testing"
)

func TestLoadAsMap(t *testing.T) {
	_ = os.Chdir("testdata")
	f, _ := os.Open(".env")
	result, err := dotenv.LoadAsMap(f)
	require.NoError(t, err)
	require.Equal(t, "6060", result["gdd.port"])
	require.Equal(t, "/api", result["gdd.route.root.path"])
}

func ErrReader(err error) io.Reader {
	return &errReader{err: err}
}

type errReader struct {
	err error
}

func (r *errReader) Read(p []byte) (int, error) {
	return 0, r.err
}

func TestLoadAsMapError(t *testing.T) {
	_, err := dotenv.LoadAsMap(ErrReader(errors.New("test error")))
	require.Error(t, err)
}

func TestLoad(t *testing.T) {
	_ = os.Chdir("testdata")
	dotenv.Load("")
	require.Equal(t, "6060", os.Getenv("GDD_PORT"))
	require.Equal(t, "/api", os.Getenv("GDD_ROUTE_ROOT_PATH"))
}
