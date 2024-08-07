package dotenv

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

func Load(env string) {
	wd, _ := os.Getwd()
	_ = godotenv.Load(filepath.Join(wd, ".env."+env+".local"))
	if "test" != env {
		_ = godotenv.Load(filepath.Join(wd, ".env.local"))
	}
	_ = godotenv.Load(filepath.Join(wd, ".env."+env))
	_ = godotenv.Load() // The Original .env
}

func LoadAsMap(reader io.Reader) (map[string]interface{}, error) {
	envMap, err := godotenv.Parse(reader)
	if err != nil {
		return nil, err
	}
	result := make(map[string]interface{})
	for key, value := range envMap {
		result[strings.ToLower(strings.ReplaceAll(key, "_", "."))] = value
	}
	return result, nil
}
