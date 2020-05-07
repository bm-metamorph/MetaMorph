package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"

	"bitbucket.com/metamorph/pkg/db/models/node"
	"bitbucket.com/metamorph/pkg/drivers/redfish"
	"bitbucket.com/metamorph/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{}

func Serve() {

	listner, err := net.Listen("tcp", ":4040")
	if err != nil {
		panic(err)
	}

	srv := grpc.NewServer()
	proto.RegisterNodeServiceServer(srv, &server{})
	reflection.Register(srv)

	if e := srv.Serve(listner); e != nil {
		panic(e)
	}

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
	err := node.UpdateRaw(nodeId,NodeSpec)
	if err == nil {
	    return &proto.Response{Result: "Update successful"}, nil

	}
	    return &proto.Response{Result: "Update failed"}, err 
}

func (s *server) Delete(ctx context.Context, request *proto.Request) (*proto.Response, error) {

	nodeId := request.GetNodeID()
	fmt.Println(nodeId)
	result := "Deleting  node"
	return &proto.Response{Result: result}, nil
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

