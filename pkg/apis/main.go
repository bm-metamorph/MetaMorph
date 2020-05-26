package api

import (
	"bitbucket.com/metamorph/pkg/logger"
	"bitbucket.com/metamorph/proto"
	"fmt"
	"github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net/http"
	"time"
)

func grpcClient() proto.NodeServiceClient {
	logger.Log.Info("grpcClient()")
	conn, err := grpc.Dial("localhost:4040", grpc.WithInsecure())
	if err != nil {
		logger.Log.Error("Failed to connect to GRPC server", zap.Error(err))
		panic(err)
	}
	client := proto.NewNodeServiceClient(conn)
	return client
}

func createNode(ctx *gin.Context) {
	logger.Log.Info("createNode()")
	client := grpcClient()
	data, _ := ctx.GetRawData()
	req := &proto.Request{NodeSpec: data}
	if response, err := client.Create(ctx, req); err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"result": fmt.Sprint(response.Result),
		})
	} else {

		logger.Log.Error("Failed to create Node", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
}

func describeNode(ctx *gin.Context) {
	logger.Log.Info("describeNode()")
	client := grpcClient()
	node_id := ctx.Param("node_id")
	fmt.Println(node_id)
	req := &proto.Request{NodeID: string(node_id)}
	if response, err := client.Describe(ctx, req); err == nil {
		ctx.Data(http.StatusOK, gin.MIMEJSON, response.Res)
	} else {
		logger.Log.Error("Failed to retrieve Node information", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
}

func updateNode(ctx *gin.Context) {
	logger.Log.Info("updateNode()")
	client := grpcClient()
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
}

func deleteNode(ctx *gin.Context) {
	logger.Log.Info("deleteNode()")
	client := grpcClient()
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
}

func deployNode(ctx *gin.Context) {
	logger.Log.Info("deployNode")
	client := grpcClient()
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
}

func listNodes(ctx *gin.Context) {
	logger.Log.Info("listNodes()")
	client := grpcClient()
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

}

func getUUID(ctx *gin.Context) {
	logger.Log.Info("getUUID")
	client := grpcClient()
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

}

func getNodeHWStatus(ctx *gin.Context) {
	logger.Log.Info("getNodeHWStatus()")
	client := grpcClient()
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
}

func updateNodeHWStatus(ctx *gin.Context) {
	logger.Log.Info("updateNodeHWStatus()")
	client := grpcClient()
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

}

func Serve() {
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

	if err := r.Run(":8080"); err != nil {
		logger.Log.Fatal("Failed to Run Server", zap.Error(err))
		//log.Fatalf("Failed to Run server: %v ", err)
	}

}
