package yaml_test

import (
	"github.com/stretchr/testify/require"
	"github.com/unionj-cloud/go-doudou/toolkit/yaml"
	"os"
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

func TestLoadReaderAsMap(t *testing.T) {
	_ = os.Chdir("testdata")
	f, _ := os.Open("app.yml")
	result, err := yaml.LoadReaderAsMap(f)
	require.NoError(t, err)
	require.Equal(t, float64(6060), result["gdd.port"])
	require.Equal(t, "go-doudou", result["gdd.tracing.metrics.root"])
}

func TestLoadFileAsMap(t *testing.T) {
	_ = os.Chdir("testdata")
	result, err := yaml.LoadFileAsMap("app.yml")
	require.NoError(t, err)
	require.Equal(t, float64(6060), result["gdd.port"])
	require.Equal(t, "go-doudou", result["gdd.tracing.metrics.root"])
}
