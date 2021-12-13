package parse

import (
	logging "github.com/ipfs/go-log/v2"
	"github.com/spf13/viper"
	"os"
	"time"
)

var (
	runtime_viper = viper.New()
	confLog       = logging.Logger("conf")
	c             Config
)

type Config struct {
	Notify  Notify  `yaml:"notify"`
	Monitor Monitor `yaml:"monitor"`
}

type Notify struct {
	Platform string `yaml:"platform"`
	Webhook  string `yaml:"webhook"`
}

type Monitor struct {
	Fullnode_api_info string `yaml:"fullnode_api_info"`
	MinerId           string `yaml:"minerId"`
}

func Init() {
	runtime_viper.SetConfigName("robin")
	runtime_viper.SetConfigType("yaml")
	runtime_viper.AddConfigPath("./conf")
	runtime_viper.AddConfigPath("../conf")
	runtime_viper.AddConfigPath("../../conf")
	err := runtime_viper.ReadInConfig()
	if err != nil {
		confLog.Errorw("init:ReadInConfig fail", "err", err)
		os.Exit(1)
	}

	err = runtime_viper.Unmarshal(&c)
	if err != nil {
		confLog.Errorw("init:Unmarshal fail", "err", err)
		os.Exit(1)
	}
}

func GetRobinConfig() Config {
	return c
}

func WatchConf(cch chan<- Config) {
	for {
		time.Sleep(1 * time.Minute)

		err := runtime_viper.ReadInConfig()
		if err != nil {
			confLog.Errorw("WatchConf:ReadInConfig fail", "err", err)
			continue
		}

		nowConfig := new(Config)
		err = runtime_viper.Unmarshal(nowConfig)
		if err != nil {
			confLog.Errorw("WatchConf:Unmarshal fail", "err", err)
			continue
		}

		if c != *nowConfig {
			c = *nowConfig
			cch <- c
		}
	}
}
