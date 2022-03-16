package yaml

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/jeremywohl/flatten"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func load(file string) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	config := make(map[string]interface{})
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}
	flat, err := flatten.Flatten(config, "", flatten.UnderscoreStyle)
	if err != nil {
		panic(err)
	}
	currentEnv := map[string]bool{}
	rawEnv := os.Environ()
	for _, rawEnvLine := range rawEnv {
		key := strings.Split(rawEnvLine, "=")[0]
		currentEnv[key] = true
	}
	for k, v := range flat {
		upperK := strings.ToUpper(k)
		if !currentEnv[upperK] {
			_ = os.Setenv(upperK, fmt.Sprint(v))
		}
	}
}

func Load(env string) {
	wd, _ := os.Getwd()
	matches, err := filepath.Glob(filepath.Join(wd, fmt.Sprintf("app-%s-local.%s", env, "y*ml")))
	if err != nil {
		panic(err)
	}
	for _, item := range matches {
		load(item)
	}
	if "test" != env {
		matches, err = filepath.Glob(filepath.Join(wd, fmt.Sprintf("app-local.%s", "y*ml")))
		if err != nil {
			panic(err)
		}
		for _, item := range matches {
			load(item)
		}
	}
	matches, err = filepath.Glob(filepath.Join(wd, fmt.Sprintf("app-%s.%s", env, "y*ml")))
	if err != nil {
		panic(err)
	}
	for _, item := range matches {
		load(item)
	}
	matches, err = filepath.Glob(filepath.Join(wd, fmt.Sprintf("app.%s", "y*ml")))
	if err != nil {
		panic(err)
	}
	for _, item := range matches {
		load(item)
	}
}
