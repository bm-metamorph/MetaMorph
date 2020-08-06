package plugin

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	hclog "github.com/hashicorp/go-hclog"

	"github.com/bm-metamorph/MetaMorph/pkg/db/models/node"
	"github.com/hashicorp/go-plugin"
	"github.com/manojkva/metamorph-plugin/common/bmh"
	"github.com/manojkva/metamorph-plugin/common/isogen"
	config "github.com/manojkva/metamorph-plugin/pkg/config"
	"github.com/manojkva/metamorph-plugin/pkg/logger"
	"go.uber.org/zap"
)

var  HandShakeMap =  map[string]plugin.HandshakeConfig {
	"metamorph-redfish-plugin": bmh.Handshake,
	"metamorph-isogen-plugin": isogen.Handshake,

}
type BMHNode struct {
	*node.Node
}

func (bmhnode *BMHNode) ReadConfigFile() error {

	var err error
	var plugins node.Plugins
	var pluginskeyname string = "plugins"
	var valueFromNode string
	var errorString string
	var apisNode []node.API

	logger.Log.Info("ReadConfigFile()")

	pluginsConfig := config.GetStringMapString(pluginskeyname)

	if pluginsConfig != nil {

		pluginsNode, err := node.GetPlugins(bmhnode.NodeUUID.String())
		if err == nil {
			apisNode, err = node.GetPluginAPIs(pluginsNode.ID)
		} else {
			logger.Log.Info("Plugin details not found in input json file. Info in config files will be used")
		}

		for key, value := range pluginsConfig {
			var api node.API
			api.Name = key
			if len(apisNode) > 0 {

				valueFromNode = node.GetPluginForAPI(apisNode, key)
			}
			if valueFromNode == "" {
				api.Plugin = value
			} else {
				logger.Log.Info(fmt.Sprintf("Value overridden for %v", key), zap.String("Old value", value), zap.String("New value", valueFromNode))
				api.Plugin = valueFromNode
			}

			plugins.APIs = append(plugins.APIs, api)
			err = node.Update(&node.Node{NodeUUID: bmhnode.NodeUUID, Plugins: plugins})

		}
	} else {
		errorString = "Failed to retrieve Plugins information from config file"
		logger.Log.Error(errorString, zap.String("NodeName", bmhnode.Name))
		err = fmt.Errorf(errorString)
	}

	return err
}

func (bmhnode *BMHNode) DispenseClientRequest(apiName string) (interface{}, error) {

	logger.Log.Info("DispenseClientRequest")

	var resultIntf interface{}

	pluginLocation := config.Get("pluginlocation").(string)

	if pluginLocation == "" {
		logger.Log.Error("Failed to retrieve pluginlocation from config file")
		return nil, fmt.Errorf("Failed to retrieve pluginlocation from config file")
	}
	pluginsNode, err := node.GetPlugins(bmhnode.NodeUUID.String())
	if err != nil {
		logger.Log.Error("Failed to retreive Plugins information from config file")
		return nil, err
	}
	apisNode, err := node.GetPluginAPIs(pluginsNode.ID)
	if err != nil {
		logger.Log.Error("Failed to retrieve supported APIS for plugin.", zap.String( "PluginID", fmt.Sprintf("%v",pluginsNode.ID)))
		return nil, err
	}

	pluginName := node.GetPluginForAPI(apisNode, apiName)

	data, err := json.Marshal(bmhnode)

	if err != nil {
		logger.Log.Error("Failed to Marshal JSON info", zap.Error(err))
		return nil, err
	}
	inputConfig := base64.StdEncoding.EncodeToString(data)

	hclogger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Trace})

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: HandShakeMap[pluginName],
		Plugins: map[string]plugin.Plugin{
			"metamorph-redfish-plugin": &bmh.BmhPlugin{},
			"metamorph-isogen-plugin":  &isogen.ISOgenPlugin{}},
		Cmd:              exec.Command("sh", "-c", pluginLocation+"/"+pluginName+" "+string(inputConfig)),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
		Logger:           hclogger})
	defer client.Kill()

	rpcClient, err := client.Client()

	if err != nil {
		logger.Log.Error("Failed to retrieve RPC Client", zap.Error(err))
		return nil, err
	}

	raw, err := rpcClient.Dispense(pluginName)
	if err != nil {
		logger.Log.Error("Failed to Dispense plugin", zap.Error(err), zap.String("PluginName", pluginName))
		return nil, err

	}
	switch apiNameLowerCase := strings.ToLower(apiName); apiNameLowerCase {
	case "getguuid":
		service := raw.(bmh.Bmh)
		resultIntf, err = service.GetGUUID()
		logger.Log.Debug("GetGUUID() ", zap.String("result",fmt.Sprintf("%v\n", resultIntf.([]byte))))
	case "deployiso":
		service := raw.(bmh.Bmh)
		err = service.DeployISO()
		resultIntf = nil
	case "updatefirmware":
		service := raw.(bmh.Bmh)
		err = service.UpdateFirmware()
		resultIntf = nil
	case "configureraid":
		service := raw.(bmh.Bmh)
		err = service.ConfigureRAID()
		resultIntf = nil
	case "createiso":
		service := raw.(isogen.ISOgen)
		err = service.CreateISO()
		resultIntf = nil
	case "gethwinventory":
		service := raw.(bmh.Bmh)
		resultIntf, err = service.GetHWInventory()
	case "poweron":
		service := raw.(bmh.Bmh)
		err = service.PowerOn()
		resultIntf = nil
	case "poweroff":
		service := raw.(bmh.Bmh)
		err = service.PowerOff()
		resultIntf = nil
	case "getpowerstatus":
		service := raw.(bmh.Bmh)
		status, _ := service.GetPowerStatus()
		if status == true {
			resultIntf = "On"
		} else {
			resultIntf = "Off"
		}
		err = nil
	default:
		logger.Log.Error("Unsupported API", zap.String("API Name", apiName))
		err = fmt.Errorf("%v not supported.", apiName)
	}
	return resultIntf, err
}
