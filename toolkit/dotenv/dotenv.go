package dotenv

import (
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
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
