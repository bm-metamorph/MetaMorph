package controller

import (
      "github.com/manojkva/metamorph-plugin/pkg/logger"
      "github.com/bm-metamorph/MetaMorph/proto"
	"github.com/bm-metamorph/MetaMorph/pkg/db/models/node"
	"github.com/google/uuid"
//	"google.golang.org/grpc/reflection"
        "fmt"
 //     "github.com/gin-contrib/zap"
//      "github.com/gin-gonic/gin"
      "go.uber.org/zap"
      "google.golang.org/grpc"
//      "net"
//      "time"
      "os"
          "encoding/json"
          "testing"
//        "reflect"
        "strings"
//      "bytes"
        "github.com/manojkva/metamorph-plugin/pkg/config"
//      "github.com/gin-gonic/gin"
           "net/http"
  //      "github.com/stretchr/testify/require"
   // "net/http/httptest"
    //    "github.com/stretchr/testify/assert"
	"github.com/gin-gonic/gin"
)


func TestServe(t *testing.T) {
	config.SetLoggerConfig("logger.apipath")
	go Serve()
        grpcServer := "localhost"
        if gs := os.Getenv("METMORPH_CONTROLLER_HOST"); gs != "" {
            grpcServer = gs
        }

        logger.Log.Info("grpcClient()")
        conn, err := grpc.Dial(fmt.Sprintf("%s:4040", grpcServer), grpc.WithInsecure())
       if err != nil {
                logger.Log.Error("Failed to connect to GRPC server", zap.Error(err))
                panic(err)
        }
        client := proto.NewNodeServiceClient(conn)
        fmt.Println(client, conn)

}


func TestCreate(t *testing.T) {
	config.SetLoggerConfig("logger.apipath")
	r := strings.NewReader(" { \"AllowFirwareUpgrade\": false }")
	req, _ := http.NewRequest("POST", "/node/", r)
	ctx := &gin.Context{
		Request: req,

	}
	data, _ := ctx.GetRawData()
	request := &proto.Request{NodeSpec: data}
	NodeSpec := request.GetNodeSpec()
        fmt.Println(string(NodeSpec))
        result, err := node.Create(NodeSpec)
	fmt.Println("printint uuid")
	fmt.Println(result,err)



}
func TestDescribe(t *testing.T) {
	//Creating node
	config.SetLoggerConfig("logger.apipath")
	r := strings.NewReader(" { \"AllowFirwareUpgrade\": false }")
	req, _ := http.NewRequest("POST", "/node/", r)
	ctx := &gin.Context{
		Request: req,

	}
	data, _ := ctx.GetRawData()
	request := &proto.Request{NodeSpec: data}
	NodeSpec := request.GetNodeSpec()
        fmt.Println(string(NodeSpec))
        nodeId, err := node.Create(NodeSpec)
	fmt.Println("printint uuid")
	fmt.Println(nodeId,err)

	//Describing node
	 result, err := node.Describe(nodeId)
	fmt.Println("Describing node")
	fmt.Println(result,err)

}


func TestDelete(t *testing.T) {
	//Creating node
	config.SetLoggerConfig("logger.apipath")
	r := strings.NewReader(" { \"AllowFirwareUpgrade\": false }")
	req, _ := http.NewRequest("POST", "/node/", r)
	ctx := &gin.Context{
		Request: req,

	}
	data, _ := ctx.GetRawData()
	request := &proto.Request{NodeSpec: data}
	NodeSpec := request.GetNodeSpec()
        fmt.Println(string(NodeSpec))
        nodeId, err := node.Create(NodeSpec)
	fmt.Println("printint uuid")
	fmt.Println(nodeId,err)

	//Describing node
	 err = node.Delete(nodeId)
	fmt.Println("Deleting node")
	fmt.Println(err)

}

func TestUpdate(t *testing.T) {
        //Creating node
        config.SetLoggerConfig("logger.apipath")
        r := strings.NewReader(" { \"AllowFirwareUpgrade\": false }")
        req, _ := http.NewRequest("POST", "/node/", r)
        ctx := &gin.Context{
                Request: req,

        }
        data, _ := ctx.GetRawData()
        request := &proto.Request{NodeSpec: data}
        NodeSpec := request.GetNodeSpec()
        fmt.Println(string(NodeSpec))
        nodeId, err := node.Create(NodeSpec)
        fmt.Println("printint uuid")
        fmt.Println(nodeId,err)

        //Describing node

	update_req := fmt.Sprintf("{\"AllowFirwareUpgrade\": true, \"NodeID\": \"%s\"}",nodeId)
	r = strings.NewReader(update_req)
	request_new, _ := http.NewRequest("PUT", fmt.Sprintf("/node/%s/", nodeId), r)
        ctx = &gin.Context{
                Request: request_new,

        }
        data, _ = ctx.GetRawData()
        request = &proto.Request{NodeSpec: data}
        NodeSpec = request.GetNodeSpec()


        var nodeUpdate node.Node
        err = json.Unmarshal(NodeSpec, &nodeUpdate)
        if err != nil {
                fmt.Println("Update failed. Invalid JSON 1")

        }
	fmt.Printf("%+v", nodeUpdate)
        nodeUpdate.NodeUUID, err = uuid.Parse(nodeId)

        err = node.Update(&nodeUpdate)
        if err == nil {
                fmt.Println("Update Successful")
        }

}

