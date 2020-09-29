package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"

	"github.com/bm-metamorph/MetaMorph/pkg/db/models/node"
	"github.com/bm-metamorph/MetaMorph/pkg/plugin"
	"github.com/bm-metamorph/MetaMorph/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"github.com/manojkva/metamorph-plugin/pkg/logger"
	"go.uber.org/zap"
)

type server struct{}

type agent struct{}

func Serve() {
	logger.Log.Info("Serve()")

	listner, err := net.Listen("tcp", ":4040")
	if err != nil {
		logger.Log.Error("Failed to receive network connection for GRPC server ", zap.Error(err))
		panic(err)
	}

	srv := grpc.NewServer()
	proto.RegisterNodeServiceServer(srv, &server{})
	proto.RegisterAgentServiceServer(srv, &agent{})
	reflection.Register(srv)

	if e := srv.Serve(listner); e != nil {
		logger.Log.Error("Failed to start GRPC Server", zap.Error(e))
		panic(e)
	}

}

func (a *agent) GetTasks(ctx context.Context, request *proto.Request) (*proto.Response, error) {

	logger.Log.Info("GetTasks()")

	nodeId := request.GetNodeID()
	bootactions, err := node.GetBootActions(nodeId)
	if err != nil {
		logger.Log.Error("Failed  GetBootActions() from DB", zap.Error(err))
		return &proto.Response{Res: nil}, err
	}
	return &proto.Response{Res: bootactions}, nil

}

func (a *agent) UpdateTaskStatus(ctx context.Context, request *proto.Request) (*proto.Response, error) {
	logger.Log.Info("UpdateTaskStatus()")
	data := request.GetTask()
	var task node.BootAction
	err := json.Unmarshal(data, &task)
	if err != nil {
		logger.Log.Error("Failed to UnMarshal json", zap.Error(err))
		fmt.Println(err)
	}

	err = node.UpdateTaskStatus(&task)
	if err != nil {
		logger.Log.Error("Failed  to update Task to DB", zap.Error(err))
		return &proto.Response{Res: nil}, err
	}
	return &proto.Response{Result: "Task status updated"}, nil
}

func (s *agent) UpdateNodeState(ctx context.Context, request *proto.Request) (*proto.Response, error) {
	logger.Log.Info("UpdateNodeState()")
	nodeId := request.GetNodeID()
	state := request.GetNodeState()
	result, err := node.Describe(nodeId)
	if err != nil {
		logger.Log.Error("Failed to retrieve node info ", zap.String("NodeId",nodeId), zap.Error(err)) 
		return &proto.Response{Res: nil}, err
	}
	var Node node.Node
	err = json.Unmarshal(result, &Node)
	Node.State = state
	err = node.Update(&Node)
	if err != nil {
		logger.Log.Error("Failed to update Node in DB", zap.Error(err))
		return &proto.Response{Res: nil}, err
	}
	return &proto.Response{Result: "Node State Updated"}, nil
}

func (s *server) Describe(ctx context.Context, request *proto.Request) (*proto.Response, error) {
	logger.Log.Info("Describe()")
	nodeId := request.GetNodeID()
	result, err := node.Describe(nodeId)
	if err != nil {
		logger.Log.Error("Failed to Describe Node", zap.String("NodeID", nodeId), zap.Error(err))
		return &proto.Response{Res: nil}, err
	}
	return &proto.Response{Res: result}, nil
}

func (s *server) Deploy(ctx context.Context, request *proto.Request) (*proto.Response, error) {

	logger.Log.Info("Deploy()")

	nodeId := request.GetNodeID()
	fmt.Println(nodeId)
	result := nodeId
	return &proto.Response{Result: result}, nil
}

func (s *server) Create(ctx context.Context, request *proto.Request) (*proto.Response, error) {

	logger.Log.Info("Create()")

	NodeSpec := request.GetNodeSpec()
	fmt.Println(string(NodeSpec))
	result := "Creating node"
	result, err := node.Create(NodeSpec)
	if err != nil {
		logger.Log.Error("Failed to create Node", zap.Error(err))
		return &proto.Response{Result: ""}, err
	}
	return &proto.Response{Result: result}, nil
}

func (s *server) Update(ctx context.Context, request *proto.Request) (*proto.Response, error) {
	logger.Log.Info("Update()")
	NodeSpec := request.GetNodeSpec()
	nodeId := request.GetNodeID()
	//fmt.Println(string(NodeSpec))
	//fmt.Println(nodeId)
	logger.Log.Debug("Node to be updated", zap.String("NodeID", nodeId))
	var nodeUpdate node.Node
	err := json.Unmarshal(NodeSpec, &nodeUpdate)
	if err != nil {
		logger.Log.Error("Update failed. Invalid JSON",zap.Error(err))
		return &proto.Response{Result: "Update failed. Invalid JSON"}, err

	}
	nodeUpdate.NodeUUID, err = uuid.Parse(nodeId)

	err = node.Update(&nodeUpdate)
	if err == nil {
		return &proto.Response{Result: "Update successful"}, nil
	}
	logger.Log.Error("Update failed", zap.Error(err))
	return &proto.Response{Result: "Update failed"}, err
}

func (s *server) Delete(ctx context.Context, request *proto.Request) (*proto.Response, error) {

	logger.Log.Info("Delete()")

	var result string

	nodeId := request.GetNodeID()
	//fmt.Println(nodeId)
	logger.Log.Debug("Node to be deleted", zap.String("NodeID", nodeId))
	err := node.Delete(nodeId)
	if err != nil {
		logger.Log.Error("Failed to Delete Node",zap.String("NodeID", nodeId), zap.Error(err))
		result = "Failed"
	}
	result = "Successful"
	return &proto.Response{Result: result}, err
}

func (s *server) List(ctx context.Context, request *proto.Request) (*proto.Response, error) {

	logger.Log.Info("List()")

	result := "List of  nodes"
	return &proto.Response{Result: result}, nil

}

func (s *server) GetNodeUUID(ctx context.Context, request *proto.Request) (*proto.Response, error) {

	logger.Log.Info("GetNodeUUID()")

	redfishClient := plugin.BMHNode{&node.Node{}}

	var result string
	data := request.GetNodeSpec()

	err := json.Unmarshal(data, &redfishClient)
	if err == nil {
		if uuid, _ := redfishClient.DispenseClientRequest("getguuid"); uuid.(string) != "" {
			result = uuid.(string)
		} else {
			errString := fmt.Sprintf("Failed to retrieve UUID from node for IPMI IP : %v", redfishClient.IPMIIP)
			logger.Log.Error(errString)
			err = errors.New(errString)
			err = errors.New(fmt.Sprintf("Failed to retrieve UUID from node for IPMI IP : %v", redfishClient.IPMIIP))
			result = ""
		}
	}else {
		logger.Log.Error("Failed to decode JSON object", zap.Error(err))
	}

	return &proto.Response{Result: result}, err

}

func (s *server) GetHWStatus(ctx context.Context, request *proto.Request) (*proto.Response, error) {
	logger.Log.Info("GetHWStatus()")
	var result string = "Off"
	nodeId := request.GetNodeID()
	data, err := node.Describe(nodeId)
	if err != nil {
		logger.Log.Error("Failed to retrieve node info", zap.String("NodeID", nodeId),zap.Error(err))
		return &proto.Response{Result: result}, err
	}
	var node node.Node
	err = json.Unmarshal(data, &node)
	if err != nil {
		logger.Log.Error("Failed to decode JSON object", zap.Error(err))
		return &proto.Response{Result: result}, err
	}
	redfishClient := &plugin.BMHNode{&node}
	status, err := redfishClient.DispenseClientRequest("getpowerstatus")

	result = status.(string)

	return &proto.Response{Result: result}, nil

}
func (s *server) UpdateHWStatus(ctx context.Context, request *proto.Request) (*proto.Response, error) {
	logger.Log.Info("UpdateHWStatus()")
	data := request.GetNodeSpec()
	nodeId := request.GetNodeID()
	nodeInfoBytes, err := node.Describe(nodeId)
	if err != nil {
		logger.Log.Error("Failed to retrieve node info", zap.String("NodeID", nodeId),zap.Error(err))
		return &proto.Response{Result: ""}, err
	}
	var node node.Node
	err = json.Unmarshal(nodeInfoBytes, &node)
	if err != nil {
		logger.Log.Error("Failed to decode JSON object", zap.Error(err))
		return &proto.Response{Result: ""}, err
	}
	redfishClient := &plugin.BMHNode{&node}

	hwInfo := new(struct {
		PowerState string
	})
	err = json.Unmarshal(data, &hwInfo)
	if err != nil {
		logger.Log.Error("Failed to decode JSON object", zap.Error(err))
		return &proto.Response{Result: ""}, err
	}
	if hwInfo.PowerState == "On" {
		_, err = redfishClient.DispenseClientRequest("poweron")
		if err  != nil {
			logger.Log.Error("PowerOn() API failed", zap.Error(err))
			return &proto.Response{Result: ""},err
		}
	} else if hwInfo.PowerState == "Off" {
		_, err = redfishClient.DispenseClientRequest("poweroff")
		if err != nil {
			logger.Log.Error("PowerOff() API failed", zap.Error(err))
			return &proto.Response{Result: ""}, err 
		}

	}

	return &proto.Response{Result: "Successfull"}, nil
}
