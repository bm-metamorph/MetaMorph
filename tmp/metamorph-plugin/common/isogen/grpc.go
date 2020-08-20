package isogen

import (
	"github.com/manojkva/metamorph-plugin/proto"
	"golang.org/x/net/context"
)

type GRPCClient struct{ client proto.IsogenClient }


func (m *GRPCClient)  CreateISO() error {
	_, err := m.client.CreateISO(context.Background(), &proto.Empty{})
	 return  err
}

type GRPCServer struct {
	Impl ISOgen
}

func (m *GRPCServer)  CreateISO(ctx context.Context, req *proto.Empty) (*proto.Empty,error) {
	/* <TBD> Add check for required parameters and raise necessary errors if reqd*/
	 err :=  m.Impl.CreateISO()
	 return &proto.Empty{},err
}
