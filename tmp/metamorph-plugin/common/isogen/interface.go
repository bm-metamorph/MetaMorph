package isogen

import (
	"context"
	"github.com/hashicorp/go-plugin"
	"github.com/manojkva/metamorph-plugin/proto"
	"google.golang.org/grpc"
)

var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "ISOGEN_PLUGIN",
	MagicCookieValue: "1.0",
}

var PluginMap = map[string]plugin.Plugin{
	"isogen": &ISOgenPlugin{},
}

type ISOgen interface {
	CreateISO()  (error)
}

type ISOgenPlugin struct {
	Impl  ISOgen
	plugin.Plugin
}

func (p *ISOgenPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterIsogenServer(s, &GRPCServer{Impl: p.Impl})
	return nil
}

func (p *ISOgenPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{client: proto.NewIsogenClient(c)}, nil

}
