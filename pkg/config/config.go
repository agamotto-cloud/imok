package config

import (
	"fmt"
	"log"
	"sync"

	"gopkg.in/yaml.v2"

	"github.com/spf13/viper"
)

// Config config
type Config struct {
	Server ServerConfig `mapstructure:"server"`
}

// ServerConfig 服务配置
type ServerConfig struct {
	Listen         string `mapstructure:"listen" yaml:"listen"`
	ClientReadBuf  int    `mapstructure:"client_read_buf" yaml:"client_read_buf"`
	ClientWriteBuf int    `mapstructure:"client_Write_buf" yaml:"client_Write_buf"`
}

func (c Config) String() string {
	if b, err := yaml.Marshal(c); err == nil {
		return string(b)
	}
	return ""
}

//Show 以 yaml 格式输出配置信息
func (c Config) Show() {
	fmt.Printf(`Config
--------------------------------------------------------------
%v--------------------------------------------------------------
`, c)
}

// C 全局配置
var (
	C    Config
	once sync.Once
)

// Initialize 初始化
func Initialize() {
	once.Do(func() {
		log.Println("initialize configuration")

		v := viper.New()
		v.SetConfigName("service")
		v.AddConfigPath(".")
		v.SetConfigType("yaml")

		if err := v.ReadInConfig(); err != nil {
			log.Println("initialize config", err)
		}

		if err := v.Unmarshal(&C); err != nil {
			log.Println("unmarshal config", err)
		}
	})
}
