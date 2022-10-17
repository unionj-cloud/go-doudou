package yaml_test

import (
	"bytes"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/yaml"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestLoad_dev(t *testing.T) {
	defer os.Clearenv()
	_ = os.Chdir("testdata")
	yaml.Load("dev")
	port := os.Getenv("GDD_PORT")
	require.Equal(t, "8080", port)

	nacosServer := os.Getenv("GDD_NACOS_SERVER_ADDR")
	require.Equal(t, "http://localhost:8848/nacos", nacosServer)
}

func TestLoad_test(t *testing.T) {
	defer os.Clearenv()
	_ = os.Chdir("testdata")
	yaml.Load("test")
	port := os.Getenv("GDD_PORT")
	require.Equal(t, "6060", port)

	nacosServer := os.Getenv("GDD_NACOS_SERVER_ADDR")
	require.Equal(t, "", nacosServer)

	rootPath := os.Getenv("GDD_ROUTE_ROOT_PATH")
	require.Equal(t, "/api", rootPath)
}

func TestLoad_error(t *testing.T) {
	defer os.Clearenv()
	_ = os.Chdir("testdata")
	require.Panics(t, func() {
		yaml.Load("prod")
	})
}

func TestLoadReaderAsMap(t *testing.T) {
	defer os.Clearenv()
	_ = os.Chdir("testdata")
	data, err := ioutil.ReadFile("app.yml")
	require.NoError(t, err)
	result, err := yaml.LoadReaderAsMap(strings.NewReader(string(data)))
	require.NoError(t, err)
	require.Equal(t, float64(6060), result["gdd.port"])
	require.Equal(t, "go-doudou", result["gdd.tracing.metrics.root"])
}

func TestLoadFileAsMap(t *testing.T) {
	defer os.Clearenv()
	_ = os.Chdir("testdata")
	result, err := yaml.LoadFileAsMap("app.yml")
	require.NoError(t, err)
	require.Equal(t, float64(6060), result["gdd.port"])
	require.Equal(t, "go-doudou", result["gdd.tracing.metrics.root"])
}

func TestLoadReaderAsMapError(t *testing.T) {
	defer os.Clearenv()
	_, err := yaml.LoadReaderAsMap(ErrReader(errors.New("test error")))
	require.Error(t, err)
}

func TestLoadFileAsMapDotenvError(t *testing.T) {
	defer os.Clearenv()
	_ = os.Chdir("testdata")
	_, err := yaml.LoadFileAsMap(".env")
	require.Error(t, err)
}

func TestLoadReaderAsMapDotenvError(t *testing.T) {
	defer os.Clearenv()
	_ = os.Chdir("testdata")
	data, err := ioutil.ReadFile(".env")
	require.NoError(t, err)
	_, err = yaml.LoadReaderAsMap(strings.NewReader(string(data)))
	require.Error(t, err)
}

func TestLoadFileAsMapError(t *testing.T) {
	defer os.Clearenv()
	_, err := yaml.LoadFileAsMap("not_exist_file")
	require.Error(t, err)
}

func TestLoadReaderAsMapFromString(t *testing.T) {
	defer os.Clearenv()
	data := []byte("gdd:\n  port: 6060\n  tracing:\n    metrics:\n      root: \"go-doudou\"")
	result, err := yaml.LoadReaderAsMap(strings.NewReader(string(data)))
	require.NoError(t, err)
	require.Equal(t, float64(6060), result["gdd.port"])
	require.Equal(t, "go-doudou", result["gdd.tracing.metrics.root"])
}

func TestLoadReader(t *testing.T) {
	defer os.Clearenv()
	_ = os.Chdir("testdata")
	data, err := ioutil.ReadFile("app.yml")
	require.NoError(t, err)
	err = yaml.LoadReader(bytes.NewReader(data))
	require.NoError(t, err)
	require.Equal(t, "6060", os.Getenv("GDD_PORT"))
	require.Equal(t, "go-doudou", os.Getenv("GDD_TRACING_METRICS_ROOT"))
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

func TestLoadReaderError(t *testing.T) {
	err := yaml.LoadReader(ErrReader(errors.New("test error")))
	require.Error(t, err)
}

func TestLoadReaderDotenvError(t *testing.T) {
	_ = os.Chdir("testdata")
	data, err := ioutil.ReadFile(".env")
	require.NoError(t, err)
	err = yaml.LoadReader(bytes.NewReader(data))
	require.Error(t, err)
}
