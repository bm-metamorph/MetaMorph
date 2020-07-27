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
)

type BMHNode struct {
	*node.Node
}

/*
- function to read the config files and set the info in  node structure right. The Contents of node structure to
  override the condig info
- Allow for adding a new entry overiding the current API support
- if there are multiple definition for the API, the first one will be chosen.
- To what level should the input json file be allowed to override the config files plugin details.
- Save it to DB.
- Use the node info for all actions.


*/
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

func (bmhnode *BMHNode) CreateClientRequest(apiname string) (*interface{}, error) {

	pluginLocation := config.Get("pluginlocation").(string)

	if pluginLocation == "" {
		return nil, fmt.Errorf("Failed to retrieve pluginlocation from config file")
	}
	pluginsNode, err := node.GetPlugins(bmhnode.NodeUUID.String())
	if err != nil {
		return nil, err
	}
	apisNode, err := node.GetPluginAPIs(pluginsNode.ID)
	if err != nil {
		return nil, err
	}

	pluginName := node.GetPluginForAPI(apisNode, apiname)

	data, err := json.Marshal(bmhnode)

	if err != nil {
		fmt.Printf("Error %v\n", err)
		return nil, err
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
		return nil, err
	}

	raw, err := rpcClient.Dispense(pluginName)
	if err != nil {
		fmt.Printf("Error %v\n", err)
		return nil, err

	}
	return &raw, err
}

func (bmhnode *BMHNode) Dispense(apiName string) error {

	raw, err := bmhnode.CreateClientRequest(apiName)

	if err != nil {
		return err
	}

	if strings.ToLower(apiName) == "getguuid" {
		var x []byte
		service := (*raw).(bmh.Bmh)
		x, err  = service.GetGUUID()
		fmt.Printf("%v\n", string(x))

	}

	return err

}
