package dotenv_test

import (
	"github.com/stretchr/testify/require"
	"github.com/unionj-cloud/go-doudou/toolkit/dotenv"
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
