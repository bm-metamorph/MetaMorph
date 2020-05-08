package controller

import(
	"context"
	"bitbucket.com/metamorph/proto"
	"bitbucket.com/metamorph/pkg/db/models/node"
	"net"
	"fmt"
	"encoding/json"
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




func(a *agent) GetTasks( ctx context.Context, request *proto.Request) (*proto.Response, error) {

	nodeId:= request.GetNodeID()
	bootactions, err := node.GetBootActions(nodeId)
	if err != nil {
		return &proto.Response{Res: nil}, err
	} 
	return &proto.Response{Res: bootactions}, nil


}

func(a *agent) UpdateTaskStatus(ctx context.Context, request *proto.Request) (*proto.Response, error) {
	data := request.GetTask()
	var task  node.BootAction
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

func (s *agent) UpdateNodeState( ctx context.Context, request *proto.Request ) ( *proto.Response, error) {
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




func (s *server) Describe( ctx context.Context, request *proto.Request ) ( *proto.Response, error) {

	nodeId:= request.GetNodeID()
	result, err := node.Describe(nodeId)
	if err != nil {
		return &proto.Response{Res: nil}, err
	}
	return &proto.Response{Res: result}, nil
}

func (s *server) Deploy( ctx context.Context, request *proto.Request ) ( *proto.Response, error) {

	nodeId:= request.GetNodeID()
	fmt.Println(nodeId)
	result := nodeId
	return &proto.Response{Result: result}, nil
}

func (s *server) Create( ctx context.Context, request *proto.Request ) ( *proto.Response, error) {

	NodeSpec:= request.GetNodeSpec()
	fmt.Println(string(NodeSpec))
	result := "Creating node"
	result, err := node.Create(NodeSpec)
	if err != nil {
		return &proto.Response{Result: ""}, err
	}
	return &proto.Response{Result: result}, nil
}

func (s *server) Update( ctx context.Context, request *proto.Request ) ( *proto.Response, error) {
	NodeSpec := request.GetNodeSpec()
	nodeId := request.GetNodeID()
	fmt.Println(string(NodeSpec))
	fmt.Println(nodeId)
	result := "Updating node"
	return &proto.Response{Result: result}, nil
}


func (s *server) Delete( ctx context.Context, request *proto.Request ) ( *proto.Response, error) {

	nodeId := request.GetNodeID()
	fmt.Println(nodeId)
	result := "Deleting  node"
	return &proto.Response{Result: result}, nil
}

func (s *server) List( ctx context.Context, request *proto.Request ) ( *proto.Response, error) {

	result := "List of  nodes"
	return &proto.Response{Result: result}, nil

}
