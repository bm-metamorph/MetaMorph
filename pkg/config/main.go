package config

import (  
	"fmt"
	"github.com/spf13/viper"
	"path"
	"github.com/bm-metamorph/MetaMorph/pkg/logger"
)

var Config string

func init() {
//	fmt.Println("Init Func")
	viper.AutomaticEnv()
	configPath := viper.GetString("METAMORPH_CONFIGPATH")

	if configPath == "" {
		gopathenv := viper.GetString("GOPATH")
		configPath = path.Join(gopathenv,"src/github.com/bm-metamorph/MetaMorph")
		viper.BindEnv("METAMORPH_CONFIGPATH",configPath)
	}
	
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)
	viper.SetEnvPrefix("metamorph")
	err := viper.ReadInConfig() 
	if err != nil { 
	  panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	Config = "Hello World"
	
}

func Get(key string) interface{} { return viper.Get(key) }

func Set(key string, value interface{}) { viper.Set(key, value) }

func SetLoggerConfig(filepathConfig string){
	loglevelString := viper.GetString("METAMORPH_LOG_LEVEL")
	level := logger.GetLogLevel(loglevelString)
	logger.InitLogger( level, path.Join(Get(filepathConfig).(string)))
}

func GetStringSlice(key string) []string  { return viper.GetStringSlice(key)}
func GetStringMapString(key string) map[string]string  { return viper.GetStringMapString(key)}

func GetStringMapStringSlice(key string) map[string][]string { return viper.GetStringMapStringSlice(key) }
