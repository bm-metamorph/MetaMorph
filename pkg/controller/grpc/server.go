package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"

	"github.com/bm-metamorph/MetaMorph/pkg/db/models/node"
	"github.com/bm-metamorph/MetaMorph/pkg/drivers/redfish"
	"github.com/bm-metamorph/MetaMorph/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{}

type agent struct{}

func Serve() {

	listner, err := net.Listen("tcp", ":4040")
	if err != nil {
		panic(err)
	}

	srv := grpc.NewServer()
	proto.RegisterNodeServiceServer(srv, &server{})
	proto.RegisterAgentServiceServer(srv, &agent{})
	reflection.Register(srv)

	if e := srv.Serve(listner); e != nil {
		panic(e)
	}

}

func (a *agent) GetTasks(ctx context.Context, request *proto.Request) (*proto.Response, error) {

	nodeId := request.GetNodeID()
	bootactions, err := node.GetBootActions(nodeId)
	if err != nil {
		return &proto.Response{Res: nil}, err
	}
	return &proto.Response{Res: bootactions}, nil

}

func (a *agent) UpdateTaskStatus(ctx context.Context, request *proto.Request) (*proto.Response, error) {
	data := request.GetTask()
	var task node.BootAction
	err := json.Unmarshal(data, &task)
	if err != nil {
		fmt.Println(err)
	}

	err = node.UpdateTaskStatus(&task)
	if err != nil {
		return &proto.Response{Res: nil}, err
	}
	return &proto.Response{Result: "Task status updated"}, nil
}

func (s *agent) UpdateNodeState(ctx context.Context, request *proto.Request) (*proto.Response, error) {
	nodeId := request.GetNodeID()
	state := request.GetNodeState()
	result, err := node.Describe(nodeId)
	if err != nil {
		return &proto.Response{Res: nil}, err
	}
	var Node node.Node
	err = json.Unmarshal(result, &Node)
	Node.State = state
	err = node.Update(&Node)
	if err != nil {
		return &proto.Response{Res: nil}, err
	}
	return &proto.Response{Result: "Node State Updated"}, nil
}

func (s *server) Describe(ctx context.Context, request *proto.Request) (*proto.Response, error) {
	nodeId := request.GetNodeID()
	result, err := node.Describe(nodeId)
	if err != nil {
		return &proto.Response{Res: nil}, err
	}
	return &proto.Response{Res: result}, nil
}

func (s *server) Deploy(ctx context.Context, request *proto.Request) (*proto.Response, error) {

	nodeId := request.GetNodeID()
	fmt.Println(nodeId)
	result := nodeId
	return &proto.Response{Result: result}, nil
}

func (s *server) Create(ctx context.Context, request *proto.Request) (*proto.Response, error) {

	NodeSpec := request.GetNodeSpec()
	fmt.Println(string(NodeSpec))
	result := "Creating node"
	result, err := node.Create(NodeSpec)
	if err != nil {
		return &proto.Response{Result: ""}, err
	}
	return &proto.Response{Result: result}, nil
}

func (s *server) Update(ctx context.Context, request *proto.Request) (*proto.Response, error) {
	NodeSpec := request.GetNodeSpec()
	nodeId := request.GetNodeID()
	//fmt.Println(string(NodeSpec))
	fmt.Println(nodeId)
	var nodeUpdate node.Node
	err := json.Unmarshal(NodeSpec, &nodeUpdate)
	if err != nil {
		return &proto.Response{Result: "Update failed. Invalid JSON"}, err

	}
	nodeUpdate.NodeUUID, err = uuid.Parse(nodeId)

	err = node.Update(&nodeUpdate)
	if err == nil {
		return &proto.Response{Result: "Update successful"}, nil
	}
	return &proto.Response{Result: "Update failed"}, err
}

func (s *server) Delete(ctx context.Context, request *proto.Request) (*proto.Response, error) {

	var result string

	nodeId := request.GetNodeID()
	fmt.Println(nodeId)
	err := node.Delete(nodeId)
	if err != nil {
		result = "Failed"
	}
	result = "Successful"
	return &proto.Response{Result: result}, err
}

func (s *server) List(ctx context.Context, request *proto.Request) (*proto.Response, error) {

	result := "List of  nodes"
	return &proto.Response{Result: result}, nil

}

func (s *server) GetNodeUUID(ctx context.Context, request *proto.Request) (*proto.Response, error) {

	nodeInfo := new(struct {
		IPMIIP       string
		IPMIUser     string
		IPMIPassword string
	})

	var result string
	data := request.GetNodeSpec()

	err := json.Unmarshal(data, &nodeInfo)
	if err == nil {
		if uuid, _ := redfish.GetUUID(nodeInfo.IPMIIP, nodeInfo.IPMIUser, nodeInfo.IPMIPassword); uuid != "" {
			result = uuid
		} else {
			err = errors.New(fmt.Sprintf("Failed to retrieve UUID from node for IPMI IP : %v", nodeInfo.IPMIIP))
			result = ""
		}
	}
	return &proto.Response{Result: result}, err

}

func (s *server) GetHWStatus(ctx context.Context, request *proto.Request) (*proto.Response, error) {
	var result string = "Off"
	nodeId := request.GetNodeID()
	data, err := node.Describe(nodeId)
	if err != nil {
		return &proto.Response{Result: result}, err
	}
	var node node.Node
	err = json.Unmarshal(data, &node)
	if err != nil {
		return &proto.Response{Result: result}, err
	}
	redfishClient := &redfish.BMHNode{&node}
	status := redfishClient.GetPowerStatus()
	if status == true {
		result = "On"
	}

	return &proto.Response{Result: result}, nil

}
func (s *server) UpdateHWStatus(ctx context.Context, request *proto.Request) (*proto.Response, error) {
	data := request.GetNodeSpec()
	nodeId := request.GetNodeID()
	nodeInfoBytes, err := node.Describe(nodeId)
	if err != nil {
		return &proto.Response{Result: ""}, err
	}
	var node node.Node
	var status bool
	err = json.Unmarshal(nodeInfoBytes, &node)
	if err != nil {
		return &proto.Response{Result: ""}, err
	}
	redfishClient := &redfish.BMHNode{&node}

	hwInfo := new(struct {
		PowerState string
	})
	err = json.Unmarshal(data, &hwInfo)
	if err != nil {
		return &proto.Response{Result: ""}, err
	}
	if hwInfo.PowerState == "On" {
		status = redfishClient.PowerOn()
		if status == false {
			return &proto.Response{Result: ""}, errors.New("Failed to Power On")
		}
	} else if hwInfo.PowerState == "Off" {
		status = redfishClient.PowerOff()
		if status == false {
			return &proto.Response{Result: ""}, errors.New("Failed to Power Off")
		}

	}

	return &proto.Response{Result: "Successfull"}, nil
}
