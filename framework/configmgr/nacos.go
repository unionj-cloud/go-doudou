package configmgr

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/unionj-cloud/go-doudou/toolkit/dotenv"
	"github.com/unionj-cloud/go-doudou/toolkit/maputils"
	"github.com/unionj-cloud/go-doudou/toolkit/yaml"
	logger "github.com/unionj-cloud/go-doudou/toolkit/zlogger"
	"github.com/wubin1989/nacos-sdk-go/v2/clients"
	"github.com/wubin1989/nacos-sdk-go/v2/clients/cache"
	"github.com/wubin1989/nacos-sdk-go/v2/clients/config_client"
	"github.com/wubin1989/nacos-sdk-go/v2/util"
	"github.com/wubin1989/nacos-sdk-go/v2/vo"
	"io"
	"os"
	"strings"
	"sync"
)

type nacosConfigType string

const (
	DotenvConfigFormat nacosConfigType = "dotenv"
	YamlConfigFormat   nacosConfigType = "yaml"
)

type NacosConfigMgr struct {
	dataIds     []string
	group       string
	format      nacosConfigType
	namespaceId string
	client      config_client.IConfigClient
	listeners   cache.ConcurrentMap
}

func (m *NacosConfigMgr) Listeners() cache.ConcurrentMap {
	return m.listeners
}

func NewNacosConfigMgr(dataIds []string, group string, format nacosConfigType, namespaceId string, client config_client.IConfigClient, listeners cache.ConcurrentMap) *NacosConfigMgr {
	return &NacosConfigMgr{dataIds: dataIds, group: group, format: format, namespaceId: namespaceId, client: client, listeners: listeners}
}

var NacosClient *NacosConfigMgr

type NacosChangeEvent struct {
	Namespace, Group, DataId string
	Changes                  map[string]maputils.Change
}

type NacosConfigListenerParam struct {
	DataId   string
	OnChange func(event *NacosChangeEvent)
}

func (m *NacosConfigMgr) fetchConfig(dataId string) (string, error) {
	content, err := m.client.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  m.group,
	})
	if err != nil {
		return "", err
	}
	return content, nil
}

func (m *NacosConfigMgr) loadDotenv(dataId string) error {
	content, err := m.fetchConfig(dataId)
	if err != nil {
		return err
	}
	envMap, err := godotenv.Parse(StringReader(content))
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
			_ = os.Setenv(key, value)
		}
	}
	return nil
}

func (m *NacosConfigMgr) loadYaml(dataId string) error {
	content, err := m.fetchConfig(dataId)
	if err != nil {
		return err
	}
	return yaml.LoadReader(strings.NewReader(content))
}

var onceNacos sync.Once
var NewConfigClient = clients.NewConfigClient

func InitialiseNacosConfig(param vo.NacosClientParam, dataId, format, group string) {
	client, err := NewConfigClient(param)
	if err != nil {
		panic(errors.Wrap(err, "[go-doudou] failed to create nacos config client"))
	}
	dataIds := strings.Split(dataId, ",")
	NacosClient = &NacosConfigMgr{
		dataIds:     dataIds,
		group:       group,
		format:      nacosConfigType(format),
		namespaceId: param.ClientConfig.NamespaceId,
		client:      client,
		listeners:   cache.NewConcurrentMap(),
	}
}

func LoadFromNacos(param vo.NacosClientParam, dataId, format, group string) error {
	onceNacos.Do(func() {
		InitialiseNacosConfig(param, dataId, format, group)
	})
	switch nacosConfigType(format) {
	case YamlConfigFormat:
		for _, item := range NacosClient.dataIds {
			if err := NacosClient.loadYaml(item); err != nil {
				return errors.Wrap(err, "[go-doudou] failed to load yaml config")
			}
		}
	case DotenvConfigFormat:
		for _, item := range NacosClient.dataIds {
			if err := NacosClient.loadDotenv(item); err != nil {
				return errors.Wrap(err, "[go-doudou] failed to load dotenv config")
			}
		}
	default:
		return fmt.Errorf("[go-doudou] unknown config format: %s\n", format)
	}
	NacosClient.listenConfig()
	return nil
}

var StringReader = func(s string) io.Reader {
	return strings.NewReader(s)
}

func (m *NacosConfigMgr) CallbackOnChange(namespace, group, dataId, data, old string) {
	var newData, oldData map[string]interface{}
	var err error
	switch m.format {
	case YamlConfigFormat:
		if newData, err = yaml.LoadReaderAsMap(StringReader(data)); err != nil {
			logger.Error().Err(err).Msg("[go-doudou] error from nacos config listener")
			return
		}
		if oldData, err = yaml.LoadReaderAsMap(StringReader(old)); err != nil {
			logger.Error().Err(err).Msg("[go-doudou] error from nacos config listener")
			return
		}
	case DotenvConfigFormat:
		if newData, err = dotenv.LoadAsMap(StringReader(data)); err != nil {
			logger.Error().Err(err).Msg("[go-doudou] error from nacos config listener")
			return
		}
		if oldData, err = dotenv.LoadAsMap(StringReader(old)); err != nil {
			logger.Error().Err(err).Msg("[go-doudou] error from nacos config listener")
			return
		}
	}
	changes := maputils.Diff(newData, oldData)
	m.onChange("__"+dataId+"__"+"registry", group, namespace, changes)
	m.onChange("__"+dataId+"__"+"ddhttp", group, namespace, changes)
	m.onChange(dataId, group, namespace, changes)
}

func (m *NacosConfigMgr) listenConfig() {
	for _, dataId := range m.dataIds {
		if err := m.client.ListenConfig(vo.ConfigParam{
			DataId:   dataId,
			Group:    m.group,
			OnChange: m.CallbackOnChange,
		}); err != nil {
			panic(err)
		}
	}
}

func (m *NacosConfigMgr) onChange(dataId, group, namespace string, changes map[string]maputils.Change) {
	key := util.GetConfigCacheKey(dataId, group, namespace)
	if v, ok := m.listeners.Get(key); ok {
		listener := v.(NacosConfigListenerParam)
		listener.OnChange(&NacosChangeEvent{
			Namespace: namespace,
			Group:     group,
			DataId:    dataId,
			Changes:   changes,
		})
	}
}

func (m *NacosConfigMgr) AddChangeListener(param NacosConfigListenerParam) {
	key := util.GetConfigCacheKey(param.DataId, m.group, m.namespaceId)
	if _, ok := m.listeners.Get(key); ok {
		logger.Warn().Msgf("[go-doudou] you have already add a config change listener for dataId: %s, you cannot override it", param.DataId)
		return
	}
	m.listeners.Set(key, param)
}
