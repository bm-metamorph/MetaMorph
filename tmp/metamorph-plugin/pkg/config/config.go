package config

import (
	"fmt"
	"github.com/manojkva/metamorph-plugin/pkg/logger"
	"github.com/spf13/viper"
	"path"
)

func init() {
	//	fmt.Println("Init Func")
	viper.AutomaticEnv()
	configPath := viper.GetString("METAMORPH_CONFIGPATH")

	if configPath == "" {
		gopathenv := viper.GetString("GOPATH")
		configPath = path.Join(gopathenv, "src/github.com/bm-metamorph/MetaMorph")
		viper.BindEnv("METAMORPH_CONFIGPATH", configPath)
	}
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)
	viper.SetEnvPrefix("metamorph")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

}

func Get(key string) interface{} { return viper.Get(key) }

func Set(key string, value interface{}) { viper.Set(key, value) }

func GetStringMapString(key string) map[string]string { return viper.GetStringMapString(key) }

func SetLoggerConfig(filepathConfig string){
	loglevelString := viper.GetString("METAMORPH_LOG_LEVEL")
	level := logger.GetLogLevel(loglevelString)
	logger.InitLogger( level, path.Join(Get(filepathConfig).(string)))
}
