package main

import(
	"context"
	"bitbucket.com/metamorph/proto"
	"net"
	"fmt"


	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{}

func main() {

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

func (s *server) Describe( ctx context.Context, request *proto.NodeID ) ( *proto.Response, error) {

	nodeId:= request.GetNodeID()
	fmt.Println(nodeId)
	result := nodeId
	return &proto.Response{Result: result}, nil


}