// Package configs
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-02 21:31
package configs

import (
	"bytes"
	_ "embed"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/LLiuHuan/gin-template/pkg/env"
	"github.com/LLiuHuan/gin-template/pkg/file"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

//var _ TomlConfigInterf = (*TomlConfigInterf)(nil)

type TomlConfigInterf interface {
	i()

	ConfigFileChangeListen()
	Clone(fileName string) TomlConfigInterf
	Get(keyName string) interface{}
	GetString(keyName string) string
	GetBool(keyName string) bool
	GetInt(keyName string) int
	GetInt32(keyName string) int32
	GetInt64(keyName string) int64
	GetFloat64(keyName string) float64
	GetDuration(keyName string) time.Duration
	GetStringSlice(keyName string) []string
	Set(keyName string, value any)
	WriteConfig() error
}

type Config struct {
	Project  Project  `mapstructure:"project" json:"project" yaml:"project" toml:"project"`     // 项目信息
	Notify   Notify   `mapstructure:"notify" json:"notify" yaml:"notify" toml:"notify"`         // 通知配置
	Redis    Redis    `mapstructure:"redis" json:"redis" yaml:"redis" toml:"redis"`             // redis配置
	DataBase DataBase `mapstructure:"database" json:"database" yaml:"database" toml:"database"` // 数据库配置
	HashIds  HashIds  `mapstructure:"hash-ids" json:"hashIds" yaml:"hash-ids" toml:"hash-ids"`  // hashIds配置
}

var (
	lastChangeTime   time.Time
	containerFactory *containers
	configContainer  TomlConfigInterf
)

var config = new(Config)

var (
	//go:embed config_dev.toml
	devConfigs []byte

	//go:embed config_pro.toml
	proConfigs []byte

	//go:embed config_fat.toml
	fatConfigs []byte

	//go:embed config_uat.toml
	uatConfigs []byte
)

// init 初始化Viper，解析配置文件
func init() {
	lastChangeTime = time.Now()
	configContainer = CreateTomlFactory()
}

func Get() *Config {
	return config
}

func GetContainer() TomlConfigInterf {
	return configContainer
}

// CreateTomlFactory 创建一个yaml配置文件工厂
// 两种方式存储配置文件，一种是解析到结构体，一种是使用KV
func CreateTomlFactory() TomlConfigInterf {
	containerFactory = CreateContainersFactory()

	var r io.Reader

	switch env.Active().Value() {
	case "dev":
		r = bytes.NewReader(devConfigs)
	case "fat":
		r = bytes.NewReader(fatConfigs)
	case "uat":
		r = bytes.NewReader(uatConfigs)
	case "pro":
		r = bytes.NewReader(proConfigs)
	default:
		r = bytes.NewReader(devConfigs)
	}
	v := viper.New()
	v.SetConfigType("toml")

	if err := v.ReadConfig(r); err != nil {
		panic(err)
	}

	if err := v.Unmarshal(config); err != nil {
		panic(err)
	}

	v.SetConfigName("config_" + env.Active().Value() + ".toml")
	v.AddConfigPath("./configs")

	configFile := "./configs/config_" + env.Active().Value() + ".toml"
	_, ok := file.IsExists(configFile)
	if !ok {
		if err := os.MkdirAll(filepath.Dir(configFile), 0766); err != nil {
			panic(err)
		}

		f, err := os.Create(configFile)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		if err := v.WriteConfig(); err != nil {
			panic(err)
		}
	}

	if err := v.ReadInConfig(); err != nil {
		panic("配置文件初始化失败" + err.Error())
	}

	return &tomlConfig{
		viper: v,
		mu:    new(sync.Mutex),
	}
}

type tomlConfig struct {
	viper *viper.Viper
	mu    *sync.Mutex
}

// ConfigFileChangeListen 监听文件变化
func (t *tomlConfig) ConfigFileChangeListen() {
	t.viper.OnConfigChange(func(changeEvent fsnotify.Event) {
		if time.Now().Sub(lastChangeTime).Seconds() >= 1 {
			if changeEvent.Op.String() == "WRITE" {
				t.clearCache()
				lastChangeTime = time.Now()
			}
		}

		if err := t.viper.Unmarshal(config); err != nil {
			panic(err)
		}
	})
	t.viper.WatchConfig()
}

// keyIsCache 判断相关键是否已经缓存
func (t *tomlConfig) keyIsCache(keyName string) bool {
	if _, exists := containerFactory.KeyIsExists(ConfigKeyPrefix + keyName); exists {
		return true
	} else {
		return false
	}
}

// 对键值进行缓存
func (t *tomlConfig) cache(keyName string, value interface{}) (bool, error) {
	// 避免缓存键、值时，程序提示键名已经被注册的日志输出
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, exists := containerFactory.KeyIsExists(ConfigKeyPrefix + keyName); exists {
		return true, nil
	}
	return containerFactory.Set(ConfigKeyPrefix+keyName, value)
}

// 通过键获取缓存的值
func (t *tomlConfig) getValueFromCache(keyName string) interface{} {
	return containerFactory.Get(ConfigKeyPrefix + keyName)
}

// 清空已经缓存的配置项信息
func (t *tomlConfig) clearCache() {
	containerFactory.FuzzyDelete(ConfigKeyPrefix)
}

// Clone 允许 clone 一个相同功能的结构体
func (t *tomlConfig) Clone(fileName string) TomlConfigInterf {
	// 这里存在一个深拷贝，需要注意，避免拷贝的结构体操作对原始结构体造成影响
	var ymlC = *t
	var ymlConfViper = *(t.viper)
	(&ymlC).viper = &ymlConfViper

	(&ymlC).viper.SetConfigName(fileName)
	if err := (&ymlC).viper.ReadInConfig(); err != nil {
		// TODO: 在考虑要不要也加到code里
		panic("初始化配置文件发生错误: " + zap.Error(err).String)
	}
	return &ymlC
}

// Get 一个原始值
func (t *tomlConfig) Get(keyName string) interface{} {
	if t.keyIsCache(keyName) {
		return t.getValueFromCache(keyName)
	} else {
		value := t.viper.Get(keyName)
		t.cache(keyName, value)
		return value
	}
}

// GetString 字符串格式返回值
func (t *tomlConfig) GetString(keyName string) string {
	if t.keyIsCache(keyName) {
		return t.getValueFromCache(keyName).(string)
	} else {
		value := t.viper.GetString(keyName)
		t.cache(keyName, value)
		return value
	}

}

// GetBool 布尔格式返回值
func (t *tomlConfig) GetBool(keyName string) bool {
	if t.keyIsCache(keyName) {
		return t.getValueFromCache(keyName).(bool)
	} else {
		value := t.viper.GetBool(keyName)
		t.cache(keyName, value)
		return value
	}
}

// GetInt 整数格式返回值
func (t *tomlConfig) GetInt(keyName string) int {
	if t.keyIsCache(keyName) {
		return t.getValueFromCache(keyName).(int)
	} else {
		value := t.viper.GetInt(keyName)
		t.cache(keyName, value)
		return value
	}
}

// GetInt32 整数格式返回值
func (t *tomlConfig) GetInt32(keyName string) int32 {
	if t.keyIsCache(keyName) {
		return t.getValueFromCache(keyName).(int32)
	} else {
		value := t.viper.GetInt32(keyName)
		t.cache(keyName, value)
		return value
	}
}

// GetInt64 整数格式返回值
func (t *tomlConfig) GetInt64(keyName string) int64 {
	if t.keyIsCache(keyName) {
		return t.getValueFromCache(keyName).(int64)
	} else {
		value := t.viper.GetInt64(keyName)
		t.cache(keyName, value)
		return value
	}
}

// GetFloat64 小数格式返回值
func (t *tomlConfig) GetFloat64(keyName string) float64 {
	if t.keyIsCache(keyName) {
		return t.getValueFromCache(keyName).(float64)
	} else {
		value := t.viper.GetFloat64(keyName)
		t.cache(keyName, value)
		return value
	}
}

// GetDuration 时间单位格式返回值
func (t *tomlConfig) GetDuration(keyName string) time.Duration {
	if t.keyIsCache(keyName) {
		return t.getValueFromCache(keyName).(time.Duration)
	} else {
		value := t.viper.GetDuration(keyName)
		t.cache(keyName, value)
		return value
	}
}

// GetStringSlice 字符串切片数格式返回值
func (t *tomlConfig) GetStringSlice(keyName string) []string {
	if t.keyIsCache(keyName) {
		return t.getValueFromCache(keyName).([]string)
	} else {
		value := t.viper.GetStringSlice(keyName)
		t.cache(keyName, value)
		return value
	}
}

func (t *tomlConfig) Set(keyName string, value any) {
	t.viper.Set(keyName, value)
}

func (t *tomlConfig) WriteConfig() error {
	return t.viper.WriteConfig()
}

func (t *tomlConfig) i() {}
