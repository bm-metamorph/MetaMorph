package plugin

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	hclog "github.com/hashicorp/go-hclog"

	config "github.com/bm-metamorph/MetaMorph/pkg/config"
	"github.com/bm-metamorph/MetaMorph/pkg/db/models/node"
	"github.com/hashicorp/go-plugin"
	"github.com/manojkva/metamorph-plugin/plugins/bmh"
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

	var pluginslist []node.Plugin
	var pluginskeyname string = "plugins"

	plugins := config.GetStringMapString(pluginskeyname)

	if plugins != nil {

		for k, _ := range plugins {
			keyvalue := pluginskeyname + "." + k
			var plugin node.Plugin
			plugin.Module = k
			plugin.Name = config.Get(keyvalue + ".name").(string)
			plugin.Location = config.Get(keyvalue + ".location").(string)
			apis := config.GetStringSlice(keyvalue + ".apis")
			fmt.Printf("%+v", apis)
			for _, api := range apis {
				var api_node node.API
				api_node.Name = api
				plugin.APIs = append(plugin.APIs, api_node)
			}
			pluginslist = append(pluginslist, plugin)

		}
		fmt.Printf("%+v", pluginslist)
	}

	//now check against the  plugins already present in DB

	// Add if plugin in not present.
	return nil
}

func (bmhnode *BMHNode) ComparePluginData(pluginData node.Plugin) {


}



func (bmhnode *BMHNode) CreateClientRequest() error {

	data, err := json.Marshal(bmhnode)
	if err != nil {
		fmt.Printf("Error %v\n", err)
		return err
	}
	inputConfig := base64.StdEncoding.EncodeToString(data)

	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Debug})
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  bmh.Handshake,
		Plugins:          bmh.PluginMap,
		Cmd:              exec.Command("sh", "-c", "../metamorph-redfish-plugin "+string(inputConfig)),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
		Logger:           logger})
	defer client.Kill()

	rpcClient, err := client.Client()

	if err != nil {
		fmt.Printf("Error %v\n", err)
		return err
	}

	raw, err := rpcClient.Dispense("bmh")
	if err != nil {
		fmt.Printf("Error %v\n", err)
		return err

	}
	service := raw.(bmh.Bmh)
	x, err := service.GetGUUID()
	fmt.Printf("%v\n", string(x))
	return nil
}
