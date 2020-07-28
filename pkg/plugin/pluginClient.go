package plugin

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	hclog "github.com/hashicorp/go-hclog"

	config "github.com/bm-metamorph/MetaMorph/pkg/config"
	"github.com/bm-metamorph/MetaMorph/pkg/db/models/node"
	"github.com/hashicorp/go-plugin"
	"github.com/manojkva/metamorph-plugin/common/bmh"
	"github.com/manojkva/metamorph-plugin/common/isogen"
)

type BMHNode struct {
	*node.Node
}

func (bmhnode *BMHNode) ReadConfigFile() error {

	var err error
	var plugins node.Plugins
	var pluginskeyname string = "plugins"

	pluginsConfig := config.GetStringMapString(pluginskeyname)

	if pluginsConfig != nil {

		pluginsNode, err := node.GetPlugins(bmhnode.NodeUUID.String())
		if err != nil {
			return err
		}
		apisNode, err := node.GetPluginAPIs(pluginsNode.ID)
		if err != nil {
			return err
		}

		for key, value := range pluginsConfig {
			var api node.API
			api.Name = key
			valueFromNode := node.GetPluginForAPI(apisNode, key)
			if valueFromNode == "" {
				api.Plugin = value
			} else {
				api.Plugin = valueFromNode
			}

			plugins.APIs = append(plugins.APIs, api)
		}
		err = node.Update(&node.Node{NodeUUID: bmhnode.NodeUUID, Plugins: plugins})
	} else {
		err = fmt.Errorf("Failed to retrieve Plugins information from config file")
	}

	return err
}

func (bmhnode *BMHNode) DispenseClientRequest(apiName string) error {

	pluginLocation := config.Get("pluginlocation").(string)

	if pluginLocation == "" {
		return fmt.Errorf("Failed to retrieve pluginlocation from config file")
	}
	pluginsNode, err := node.GetPlugins(bmhnode.NodeUUID.String())
	if err != nil {
		return err
	}
	apisNode, err := node.GetPluginAPIs(pluginsNode.ID)
	if err != nil {
		return err
	}

	pluginName := node.GetPluginForAPI(apisNode, apiName)

	data, err := json.Marshal(bmhnode)

	if err != nil {
		fmt.Printf("Error %v\n", err)
		return err
	}
	inputConfig := base64.StdEncoding.EncodeToString(data)

	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Trace})

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  bmh.Handshake,
		Plugins:          map[string]plugin.Plugin{pluginName: &bmh.BmhPlugin{}},
		Cmd:              exec.Command("sh", "-c", pluginLocation+"/"+pluginName+" "+string(inputConfig)),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
		Logger:           logger})
	defer client.Kill()

	rpcClient, err := client.Client()

	if err != nil {
		fmt.Printf("Error %v\n", err)
		return err
	}

	raw, err := rpcClient.Dispense(pluginName)
	if err != nil {
		fmt.Printf("Error %v\n", err)
		return err

	}
	switch apiNameLowerCase := strings.ToLower(apiName); apiNameLowerCase {
	case "getguuid":
		service := raw.(bmh.Bmh)
		var x []byte
		x, err = service.GetGUUID()
		fmt.Printf("%v\n", string(x))
	case "deployiso":
		service := raw.(bmh.Bmh)
		err = service.DeployISO()
	case "updatefirmware":
		service := raw.(bmh.Bmh)
		err = service.UpdateFirmware()
	case "configureraid":
		service := raw.(bmh.Bmh)
		err = service.ConfigureRAID()
	case "createiso":
		service := raw.(isogen.ISOgen)
		err = service.CreateISO()
	default:
		err = fmt.Errorf("%v not supported.", apiName)
	}
	return err
}
