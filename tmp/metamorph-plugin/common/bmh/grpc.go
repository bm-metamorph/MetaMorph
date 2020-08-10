package bmh

import (
	"github.com/manojkva/metamorph-plugin/proto"
	"golang.org/x/net/context"
)

type GRPCClient struct{ client proto.BmhClient }

func (m *GRPCClient) GetGUUID() ([]byte, error) {
	resp, err := m.client.GetGUUID(context.Background(), &proto.Empty{})
	if err != nil {
		return nil, err
	}

	return resp.Value, nil
}
func (m *GRPCClient) GetPowerStatus() (bool, error) {
	resp, err := m.client.GetPowerStatus(context.Background(), &proto.Empty{})
	return resp.Status, err
}

func (m *GRPCClient) UpdateFirmware() error {
	_, err := m.client.UpdateFirmware(context.Background(), &proto.Empty{})
	return err
}
func (m *GRPCClient) ConfigureRAID() error {
	_, err := m.client.ConfigureRAID(context.Background(), &proto.Empty{})
	return err
}
func (m *GRPCClient) DeployISO() error {
	_, err := m.client.DeployISO(context.Background(), &proto.Empty{})
	return err
}
func (m *GRPCClient) GetHWInventory() (map[string]string,error) {
	hwInfo, err := m.client.GetHWInventory(context.Background(), &proto.Empty{})
	return hwInfo.MapInfo,err
}
func (m *GRPCClient) PowerOff() error {
	_, err := m.client.PowerOff(context.Background(), &proto.Empty{})
	return err
}
func (m *GRPCClient) PowerOn() error {
	_, err := m.client.PowerOn(context.Background(), &proto.Empty{})
	return err
}

type GRPCServer struct {
	Impl Bmh
}

func (m *GRPCServer) GetGUUID(ctx context.Context, req *proto.Empty) (*proto.ResponseByte, error) {
	/* <TBD> Add check for required parameters and raise necessary errors if reqd*/
	v, err := m.Impl.GetGUUID()
	return &proto.ResponseByte{Value: v}, err
}
func (m *GRPCServer) GetPowerStatus(ctx context.Context, req *proto.Empty) (*proto.ResponseBool, error) {
	/* <TBD> Add check for required parameters and raise necessary errors if reqd*/
	s,err := m.Impl.GetPowerStatus()
	return &proto.ResponseBool{Status: s},err
}
func (m *GRPCServer) UpdateFirmware(ctx context.Context, req *proto.Empty) (*proto.Empty, error) {
	/* <TBD> Add check for required parameters and raise necessary errors if reqd*/
	err := m.Impl.UpdateFirmware()
	return &proto.Empty{}, err
}
func (m *GRPCServer) ConfigureRAID(ctx context.Context, req *proto.Empty) (*proto.Empty, error) {
	/* <TBD> Add check for required parameters and raise necessary errors if reqd*/
	err := m.Impl.ConfigureRAID()
	return &proto.Empty{}, err
}
func (m *GRPCServer) DeployISO(ctx context.Context, req *proto.Empty) (*proto.Empty, error) {
	/* <TBD> Add check for required parameters and raise necessary errors if reqd*/
	err := m.Impl.DeployISO()
	return &proto.Empty{}, err
}
func (m *GRPCServer) GetHWInventory(ctx context.Context, req *proto.Empty) (*proto.ResponseMap, error) {
	/* <TBD> Add check for required parameters and raise necessary errors if reqd*/
	hwInfo,err := m.Impl.GetHWInventory()
	return &proto.ResponseMap{MapInfo: hwInfo}, err
}
func (m *GRPCServer) PowerOff(ctx context.Context, req *proto.Empty) (*proto.Empty, error) {
	/* <TBD> Add check for required parameters and raise necessary errors if reqd*/
	err := m.Impl.PowerOff()
	return &proto.Empty{}, err
}
func (m *GRPCServer) PowerOn(ctx context.Context, req *proto.Empty) (*proto.Empty, error) {
	/* <TBD> Add check for required parameters and raise necessary errors if reqd*/
	err := m.Impl.PowerOn()
	return &proto.Empty{}, err
}
