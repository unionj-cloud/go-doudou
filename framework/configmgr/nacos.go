package configmgr

import (
	"github.com/joho/godotenv"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/pkg/errors"
	"github.com/unionj-cloud/go-doudou/toolkit/stringutils"
	"github.com/unionj-cloud/go-doudou/toolkit/yaml"
	"os"
	"strings"
)

var NacosConfigClient config_client.IConfigClient

const (
	DotenvConfigFormat = "dotenv"
	YamlConfigFormat   = "yaml"
)

func fetchConfig(dataId, group string) (string, error) {
	content, err := NacosConfigClient.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})
	if err != nil {
		return "", err
	}
	return content, nil
}

func loadDotenv(dataId, group string) error {
	content, err := fetchConfig(dataId, group)
	if err != nil {
		return err
	}
	envMap, err := godotenv.Parse(strings.NewReader(content))
	if err != nil {
		return err
	}
	currentEnv := map[string]bool{}
	rawEnv := os.Environ()
	for _, rawEnvLine := range rawEnv {
		key := strings.Split(rawEnvLine, "=")[0]
		currentEnv[key] = true
	}
	for key, value := range envMap {
		if !currentEnv[key] {
			os.Setenv(key, value)
		}
	}
	return nil
}

func loadYaml(dataId, group string) error {
	content, err := fetchConfig(dataId, group)
	if err != nil {
		return err
	}
	return yaml.LoadReader(strings.NewReader(content))
}

func LoadFromNacos(env string, param vo.NacosClientParam, service, format, group string) error {
	if stringutils.IsEmpty(service) {
		return errors.New("service name is required")
	}
	var err error
	NacosConfigClient, err = clients.NewConfigClient(param)
	if err != nil {
		return errors.Wrap(err, "[go-doudou] failed to create nacos config client")
	}
	switch format {
	case YamlConfigFormat:
		err = loadYaml(service+"-"+env+".yml", group)
		if err != nil {
			return err
		}
		err = loadYaml(service+".yml", group)
		if err != nil {
			return err
		}
		err = loadYaml("app.yml", group)
		if err != nil {
			return err
		}
	default:
		err = loadDotenv(service+".env."+env, group)
		if err != nil {
			return err
		}
		err = loadDotenv(service+".env", group)
		if err != nil {
			return err
		}
		err = loadDotenv(".env", group)
		if err != nil {
			return err
		}
	}
	return nil
}
