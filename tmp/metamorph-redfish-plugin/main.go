package main

import (
	config "github.com/manojkva/metamorph-plugin/pkg/config"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-plugin"
	"github.com/manojkva/metamorph-plugin/common/bmh"
	"github.com/manojkva/metamorph-plugin/pkg/logger"
	driver "github.com/manojkva/metamorph-redfish-plugin/pkg/redfish"
	"os"
	"encoding/base64"
	"go.uber.org/zap"
)

func main() {
	config.SetLoggerConfig("logger.plugins.redfishpluginpath")
	if len(os.Args) != 2 {
		logger.Log.Error(fmt.Sprintf("Usage metamorph-redfish-plugin <inputConfig>"))
		os.Exit(1)
	}
	data := os.Args[1]


	var bmhnode driver.BMHNode

        inputConfig,err :=  base64.StdEncoding.DecodeString(data)

	if  err != nil {
		logger.Log.Error(fmt.Sprintf("Failed to decode input config %v\n", data))
		logger.Log.Error(fmt.Sprintf("Error %v\n", err),zap.Error(err))
		os.Exit(1)
	}

	err = json.Unmarshal([]byte(inputConfig), &bmhnode)

	if err != nil {

		logger.Log.Error(fmt.Sprintf("Failed to decode input config %v\n", inputConfig))
		logger.Log.Error(fmt.Sprintf("Error %v\n", err),zap.Error(err))
		os.Exit(1)
	}
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: bmh.Handshake,
		Plugins: map[string]plugin.Plugin{
			"metamorph-redfish-plugin": &bmh.BmhPlugin{Impl: &bmhnode}},
		GRPCServer: plugin.DefaultGRPCServer})
}
