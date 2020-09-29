package api

import (
	"github.com/manojkva/metamorph-plugin/pkg/logger"
	"github.com/bm-metamorph/MetaMorph/proto"
	"fmt"
	"github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
        ctrlgRPCServer "github.com/bm-metamorph/MetaMorph/pkg/controller/grpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net/http"
	"time"
        "os"
)
func run_controller() {
  ctrlgRPCServer.Serve() 

}
func grpcClient() (proto.NodeServiceClient, *grpc.ClientConn) {

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
	return client, conn
}

func createNode(ctx *gin.Context) {
	logger.Log.Info("createNode()")
	client, conn := grpcClient()
	data, _ := ctx.GetRawData()
	req := &proto.Request{NodeSpec: data}
	if response, err := client.Create(ctx, req); err == nil {
		fmt.Println(response)
		ctx.JSON(http.StatusOK, gin.H{
			"result": fmt.Sprint(response.Result),
		})
	} else {

		logger.Log.Error("Failed to create Node", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
        defer conn.Close()
}

func describeNode(ctx *gin.Context) {
	logger.Log.Info("describeNode()")
	client, conn := grpcClient()
	node_id := ctx.Param("node_id")
	logger.Log.Debug("Node info from Request", zap.String("NodeID", string(node_id)))
	req := &proto.Request{NodeID: string(node_id)}
	if response, err := client.Describe(ctx, req); err == nil {
		ctx.Data(http.StatusOK, gin.MIMEJSON, response.Res)
	} else {
		logger.Log.Error("Failed to retrieve Node information", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
        defer conn.Close()
}

func updateNode(ctx *gin.Context) {
	logger.Log.Info("updateNode()")
	client, conn := grpcClient()
	node_id := ctx.Param("node_id")
	data, _ := ctx.GetRawData()
	req := &proto.Request{NodeID: string(node_id), NodeSpec: data}
	if response, err := client.Update(ctx, req); err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"result": fmt.Sprint(response.Result),
		})
	} else {
		logger.Log.Error("Failed to update Node", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
        defer conn.Close()
}

func deleteNode(ctx *gin.Context) {
	logger.Log.Info("deleteNode()")
	client, conn := grpcClient()
	node_id := ctx.Param("node_id")
	req := &proto.Request{NodeID: string(node_id)}
	if response, err := client.Delete(ctx, req); err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"result": fmt.Sprint(response.Result),
		})
	} else {
		logger.Log.Error("Failed to delete node", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
        defer conn.Close()
}

func deployNode(ctx *gin.Context) {
	logger.Log.Info("deployNode")
	client, conn := grpcClient()
	node_id := ctx.Param("node_id")
	req := &proto.Request{NodeID: string(node_id)}
	if response, err := client.Deploy(ctx, req); err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"result": fmt.Sprintf("Node %s deployed", response.Result),
		})
	} else {
		logger.Log.Error("Failed to deploy Node", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
        defer conn.Close()
}

func listNodes(ctx *gin.Context) {
	logger.Log.Info("listNodes()")
	client, conn := grpcClient()
	req := &proto.Request{}
	if response, err := client.List(ctx, req); err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"result": fmt.Sprint(response.Result),
		})
	} else {
		logger.Log.Error("Failed to list Nodes", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
        defer conn.Close()
}

func getUUID(ctx *gin.Context) {
	logger.Log.Info("getUUID")
	client, conn := grpcClient()
	data, _ := ctx.GetRawData()
	req := &proto.Request{NodeSpec: data}
	if response, err := client.GetNodeUUID(ctx, req); err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"result": fmt.Sprint(response.Result),
		})
	} else {
		logger.Log.Error("Failed to retrieve Node UUID", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
        defer conn.Close()

}

func getNodeHWStatus(ctx *gin.Context) {
	logger.Log.Info("getNodeHWStatus()")
	client, conn := grpcClient()
	node_id := ctx.Param("node_id")

	req := &proto.Request{NodeID: string(node_id)}
	if response, err := client.GetHWStatus(ctx, req); err == nil {
		ctx.JSON(http.StatusOK, gin.H{"result": fmt.Sprint(response.Result)})
	} else {
		logger.Log.Error("Failed to retrieve HW status from Node", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
        defer conn.Close()
}

func updateNodeHWStatus(ctx *gin.Context) {
	logger.Log.Info("updateNodeHWStatus()")
	client, conn := grpcClient()
	node_id := ctx.Param("node_id")
	data, _ := ctx.GetRawData()
	req := &proto.Request{NodeID: string(node_id), NodeSpec: data}
	if response, err := client.UpdateHWStatus(ctx, req); err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"result": fmt.Sprint(response.Result),
		})
	} else {
		logger.Log.Error("Failed to Update HW status", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
        defer conn.Close()

}

func Serve() *gin.Engine {
	logger.Log.Info("Serve()")

	r := gin.Default()
	r.Use(ginzap.Ginzap(logger.Log, time.RFC3339, false))
	r.Use(ginzap.RecoveryWithZap(logger.Log, true))

	r.GET("/nodes", listNodes)
	r.POST("/uuid", getUUID)
	r.GET("/hwstatus/:node_id", getNodeHWStatus)
	r.PUT("/hwstatus/:node_id", updateNodeHWStatus)
	node := r.Group("/node")
	{
		node.POST("/", createNode)
		node.GET("/:node_id", describeNode)
		node.PUT("/:node_id", updateNode)
		node.DELETE("/:node_id", deleteNode)
		node.POST("/deploy/:node_id", deployNode)
	}
	return r
}

func api() {
	        if err := Serve().Run(":8080"); err != nil {
                logger.Log.Fatal("Failed to Run Server", zap.Error(err))
                //log.Fatalf("Failed to Run server: %v ", err)
        }

}
