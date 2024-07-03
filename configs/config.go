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
	configContainer = CreateYamlFactory()
}

func Get() *Config {
	return config
}

func GetContainer() TomlConfigInterf {
	return configContainer
}

// CreateYamlFactory 创建一个yaml配置文件工厂
// 两种方式存储配置文件，一种是解析到结构体，一种是使用KV
func CreateYamlFactory() TomlConfigInterf {
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

	return &ymlConfig{
		viper: v,
		mu:    new(sync.Mutex),
	}
}

type ymlConfig struct {
	viper *viper.Viper
	mu    *sync.Mutex
}

// ConfigFileChangeListen 监听文件变化
func (y *ymlConfig) ConfigFileChangeListen() {
	y.viper.OnConfigChange(func(changeEvent fsnotify.Event) {
		if time.Now().Sub(lastChangeTime).Seconds() >= 1 {
			if changeEvent.Op.String() == "WRITE" {
				y.clearCache()
				lastChangeTime = time.Now()
			}
		}

		if err := y.viper.Unmarshal(config); err != nil {
			panic(err)
		}
	})
	y.viper.WatchConfig()
}

// keyIsCache 判断相关键是否已经缓存
func (y *ymlConfig) keyIsCache(keyName string) bool {
	if _, exists := containerFactory.KeyIsExists(ConfigKeyPrefix + keyName); exists {
		return true
	} else {
		return false
	}
}

// 对键值进行缓存
func (y *ymlConfig) cache(keyName string, value interface{}) (bool, error) {
	// 避免缓存键、值时，程序提示键名已经被注册的日志输出
	y.mu.Lock()
	defer y.mu.Unlock()
	if _, exists := containerFactory.KeyIsExists(ConfigKeyPrefix + keyName); exists {
		return true, nil
	}
	return containerFactory.Set(ConfigKeyPrefix+keyName, value)
}

// 通过键获取缓存的值
func (y *ymlConfig) getValueFromCache(keyName string) interface{} {
	return containerFactory.Get(ConfigKeyPrefix + keyName)
}

// 清空已经缓存的配置项信息
func (y *ymlConfig) clearCache() {
	containerFactory.FuzzyDelete(ConfigKeyPrefix)
}

// Clone 允许 clone 一个相同功能的结构体
func (y *ymlConfig) Clone(fileName string) TomlConfigInterf {
	// 这里存在一个深拷贝，需要注意，避免拷贝的结构体操作对原始结构体造成影响
	var ymlC = *y
	var ymlConfViper = *(y.viper)
	(&ymlC).viper = &ymlConfViper

	(&ymlC).viper.SetConfigName(fileName)
	if err := (&ymlC).viper.ReadInConfig(); err != nil {
		// TODO: 在考虑要不要也加到code里
		panic("初始化配置文件发生错误: " + zap.Error(err).String)
	}
	return &ymlC
}

// Get 一个原始值
func (y *ymlConfig) Get(keyName string) interface{} {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName)
	} else {
		value := y.viper.Get(keyName)
		y.cache(keyName, value)
		return value
	}
}

// GetString 字符串格式返回值
func (y *ymlConfig) GetString(keyName string) string {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(string)
	} else {
		value := y.viper.GetString(keyName)
		y.cache(keyName, value)
		return value
	}

}

// GetBool 布尔格式返回值
func (y *ymlConfig) GetBool(keyName string) bool {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(bool)
	} else {
		value := y.viper.GetBool(keyName)
		y.cache(keyName, value)
		return value
	}
}

// GetInt 整数格式返回值
func (y *ymlConfig) GetInt(keyName string) int {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(int)
	} else {
		value := y.viper.GetInt(keyName)
		y.cache(keyName, value)
		return value
	}
}

// GetInt32 整数格式返回值
func (y *ymlConfig) GetInt32(keyName string) int32 {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(int32)
	} else {
		value := y.viper.GetInt32(keyName)
		y.cache(keyName, value)
		return value
	}
}

// GetInt64 整数格式返回值
func (y *ymlConfig) GetInt64(keyName string) int64 {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(int64)
	} else {
		value := y.viper.GetInt64(keyName)
		y.cache(keyName, value)
		return value
	}
}

// GetFloat64 小数格式返回值
func (y *ymlConfig) GetFloat64(keyName string) float64 {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(float64)
	} else {
		value := y.viper.GetFloat64(keyName)
		y.cache(keyName, value)
		return value
	}
}

// GetDuration 时间单位格式返回值
func (y *ymlConfig) GetDuration(keyName string) time.Duration {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(time.Duration)
	} else {
		value := y.viper.GetDuration(keyName)
		y.cache(keyName, value)
		return value
	}
}

// GetStringSlice 字符串切片数格式返回值
func (y *ymlConfig) GetStringSlice(keyName string) []string {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).([]string)
	} else {
		value := y.viper.GetStringSlice(keyName)
		y.cache(keyName, value)
		return value
	}
}

func (y *ymlConfig) i() {}
