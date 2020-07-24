package test

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	hclog "github.com/hashicorp/go-hclog"

	"github.com/hashicorp/go-plugin"
	"github.com/manojkva/metamorph-plugin/plugins/bmh"
)

func TestClientRequest(t *testing.T) {
	data, err := ioutil.ReadFile("../examples/nodeip.json")
	if err != nil {
		fmt.Printf("Could not read input config file\n")
		os.Exit(1)
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
		os.Exit(1)
	}

	raw, err := rpcClient.Dispense("bmh")
	if err != nil {
		fmt.Printf("Error %v\n", err)
		os.Exit(1)

	}
	service := raw.(bmh.Bmh)
	x, err := service.GetGUUID()
	fmt.Printf("%v\n", string(x))
}
