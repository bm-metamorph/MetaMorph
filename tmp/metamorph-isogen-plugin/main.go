package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	config "github.com/manojkva/metamorph-plugin/pkg/config"
	"github.com/hashicorp/go-plugin"
	driver "github.com/manojkva/metamorph-isogen-plugin/pkg/isogen"
	"github.com/manojkva/metamorph-plugin/common/isogen"
	"github.com/manojkva/metamorph-plugin/pkg/logger"
	"os"
)

func main() {
	config.SetLoggerConfig("logger.plugins.isogenpluginpath")
	if len(os.Args) != 2 {
		logger.Log.Error(fmt.Sprintf("Usage metamorph-isogen-plugin <uuid>"))
		os.Exit(1)
	}
	data := os.Args[1]

	var bmhnode driver.BMHNode

	inputConfig, err := base64.StdEncoding.DecodeString(data)

	if err != nil {
		logger.Log.Error(fmt.Sprintf("Failed to decode input config %v\n", data))
		logger.Log.Error(fmt.Sprintf("Error %v\n", err))
		os.Exit(1)
	}

	err = json.Unmarshal([]byte(inputConfig), &bmhnode)
	if err != nil {

		logger.Log.Error(fmt.Sprintf("Failed to decode input config %v\n", inputConfig))
		logger.Log.Error(fmt.Sprintf("Error %v\n", err))
		os.Exit(1)
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: isogen.Handshake,
		Plugins: map[string]plugin.Plugin{
			"metamorph-isogen-plugin": &isogen.ISOgenPlugin{Impl: &bmhnode}},
		GRPCServer: plugin.DefaultGRPCServer})
}
