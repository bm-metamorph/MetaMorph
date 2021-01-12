package test

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	hclog "github.com/hashicorp/go-hclog"

	"encoding/base64"
	"github.com/hashicorp/go-plugin"
	"github.com/manojkva/metamorph-plugin/common/isogen"
	"io/ioutil"
)

func TestClientRequest(t *testing.T) {
	data, err := ioutil.ReadFile("../examples/node1_input.json")
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
		HandshakeConfig:  isogen.Handshake,
		Plugins: map[string]plugin.Plugin{ "metamorph-isogen-plugin": &isogen.ISOgenPlugin{}},
		Cmd:              exec.Command("sh", "-c", "../metamorph-isogen-plugin "+string(inputConfig)),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
		Logger:           logger})
	defer client.Kill()

	rpcClient, err := client.Client()

	if err != nil {
		fmt.Printf("Error %v\n", err)
		os.Exit(1)
	}

	raw, err := rpcClient.Dispense("metamorph-isogen-plugin")
	if err != nil {
		fmt.Printf("Error %v\n", err)
		os.Exit(1)

	}
	service := raw.(isogen.ISOgen)
	err = service.CreateISO()
	if err != nil {
		fmt.Printf("Erro %v\n", err)
	} else {
		fmt.Printf("Successfull ISO Creation")
	}
}
