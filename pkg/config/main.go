package config

import "fmt"
import "github.com/spf13/viper"

var Config string

func init() {
	fmt.Println("Init Func")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
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
