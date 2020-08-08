package bmh

import (
	"context"
	"github.com/hashicorp/go-plugin"
	"github.com/manojkva/metamorph-plugin/proto"
	"google.golang.org/grpc"
)

var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BMH_PLUGIN",
	MagicCookieValue: "1.0",
}

var PluginMap = map[string]plugin.Plugin{
	"bmh": &BmhPlugin{},
}

type Bmh interface {
	GetGUUID() ([]byte, error)
	DeployISO() (error)
	UpdateFirmware() (error)
	ConfigureRAID()  (error)
	GetHWInventory() (map[string]string, error)
	PowerOff()  (error)
	PowerOn()   (error)
	GetPowerStatus()(bool,error)
}

type BmhPlugin struct {
	Impl Bmh 
	plugin.Plugin
}

func (p *BmhPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterBmhServer(s, &GRPCServer{Impl: p.Impl})
	return nil
}

func (p *BmhPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{client: proto.NewBmhClient(c)}, nil

}
