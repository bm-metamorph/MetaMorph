package config

import (  
	"fmt"
	"github.com/spf13/viper"
	"os"
)

var Config string

func init() {
	fmt.Println("Init Func")
	configPath := os.Getenv("METAMORPH_CONFIGPATH")
	if configPath == ""{
		panic(fmt.Errorf("Fatal erro Config file path not found. Set METAMORPH_CONFIGPATH environment variable"))
	}
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)
	viper.SetEnvPrefix("metamorph")
	viper.AutomaticEnv()
	err := viper.ReadInConfig() 
	if err != nil { 
	  panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	Config = "Hello World"
	
}

func Get(key string) interface{} { return viper.Get(key) }

func Set(key string, value interface{}) { viper.Set(key, value) }
