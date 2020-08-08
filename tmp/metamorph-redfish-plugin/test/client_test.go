package test

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"io/ioutil"

        hclog "github.com/hashicorp/go-hclog"

	"github.com/hashicorp/go-plugin"
	"github.com/manojkva/metamorph-plugin/common/bmh"
)

func TestClientRequest(t *testing.T) {
	data,err := ioutil.ReadFile("../examples/nodeip.json" )
	if  err  != nil{
	    fmt.Printf("Could not read input config file\n")
	    os.Exit(1)
	}
	inputConfig  :=  base64.StdEncoding.EncodeToString(data)

	logger := hclog.New(&hclog.LoggerOptions{
		  Name: "plugin",
		  Output: os.Stdout,
		  Level: hclog.Trace,})
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  bmh.Handshake,
		Plugins: map[string]plugin.Plugin{
			                        "metamorph-redfish-plugin": &bmh.BmhPlugin{}},
		Cmd: exec.Command("sh", "-c", "../metamorph-redfish-plugin " + string(inputConfig) ),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
	        Logger: logger,})
	defer client.Kill()

	rpcClient, err := client.Client()

	if err != nil {
		fmt.Printf("Error %v\n", err)
		os.Exit(1)
	}

	raw, err := rpcClient.Dispense("metamorph-redfish-plugin")
	if err != nil {
		fmt.Printf("Error %v\n", err)
		os.Exit(1)

	}
	service := raw.(bmh.Bmh)
  //x, err := service.GetPowerStatus()
  x, err := service.GetGUUID()
  fmt.Printf("%v\n", x)
}
