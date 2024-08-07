package configmgr

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
	"github.com/apolloconfig/agollo/v4/storage"
	"github.com/pkg/errors"
	logger "github.com/unionj-cloud/go-doudou/v2/toolkit/zlogger"
)

var onceApollo sync.Once
var ApolloClient agollo.Client
var StartWithConfig = agollo.StartWithConfig

func InitialiseApolloConfig(appConfig *config.AppConfig) {
	var err error
	ApolloClient, err = StartWithConfig(func() (*config.AppConfig, error) {
		return appConfig, nil
	})
	if err != nil {
		panic(errors.Wrap(err, "[go-doudou] failed to initialise apollo client"))
	}
	logger.Info().Msg("[go-doudou] initialise apollo client successfully")
}

func LoadFromApollo(appConfig *config.AppConfig) {
	onceApollo.Do(func() {
		InitialiseApolloConfig(appConfig)
	})
	currentEnv := map[string]bool{}
	namespaces := strings.Split(appConfig.NamespaceName, ",")
	for _, item := range namespaces {
		rawEnv := os.Environ()
		for _, rawEnvLine := range rawEnv {
			key := strings.Split(rawEnvLine, "=")[0]
			currentEnv[key] = true
		}
		cache := ApolloClient.GetConfigCache(item)
		cache.Range(func(key, value interface{}) bool {
			logger.Debug().Msgf("[go-doudou] key: %s, value: %s\n", key, value)
			upperK := strings.ToUpper(strings.ReplaceAll(key.(string), ".", "_"))
			if !currentEnv[upperK] {
				_ = os.Setenv(upperK, fmt.Sprint(value))
			}
			return true
		})
	}
}

type BaseApolloListener struct {
	SkippedFirstEvent bool
	Lock              sync.Mutex
}

func (c *BaseApolloListener) OnNewestChange(event *storage.FullChangeEvent) {
}
