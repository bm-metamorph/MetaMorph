package api

import(
	"fmt"
	"bitbucket.com/metamorph/proto"
	"github.com/gin-gonic/gin"
	"net/http"
	"log"
	"google.golang.org/grpc"
)

func grpcClient() ( proto.NodeServiceClient) {
	conn, err := grpc.Dial( "localhost:4040", grpc.WithInsecure() )
	if err != nil { panic(err) }
	client := proto.NewNodeServiceClient(conn) 
	return client
}

 func createNode(ctx *gin.Context){
	client := grpcClient()
	data,_ := ctx.GetRawData()
	req := &proto.Request{NodeSpec: data }
	if response, err := client.Create(ctx, req); err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"result": fmt.Sprint(response.Result),
		})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
}

func describeNode(ctx *gin.Context){
	client := grpcClient()
	node_id := ctx.Param("node_id")
	fmt.Println(node_id)
	req := &proto.Request{NodeID: string(node_id)}
	if response, err := client.Describe(ctx, req); err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"result": fmt.Sprint(response.Result),
		})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
}

func updateNode(ctx *gin.Context){
	client := grpcClient()
	node_id := ctx.Param("node_id")
	data,_ := ctx.GetRawData()
	req := &proto.Request{NodeID: string(node_id), NodeSpec: data }
	if response, err := client.Update(ctx, req); err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"result": fmt.Sprint(response.Result),
		})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
}

func deleteNode(ctx *gin.Context){
	client := grpcClient()
	node_id := ctx.Param("node_id")
	req := &proto.Request{NodeID: string(node_id)}
	if response, err := client.Delete(ctx, req); err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"result": fmt.Sprint(response.Result),
		})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
}

func deployNode(ctx *gin.Context){
	client := grpcClient()
	node_id := ctx.Param("node_id")
	req := &proto.Request{NodeID: string(node_id)}
	if response, err := client.Deploy(ctx, req); err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"result": fmt.Sprintf("Node %s deployed", response.Result),
		})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
}

func listNodes(ctx *gin.Context){
	client := grpcClient()
	req := &proto.Request{}
	if response, err := client.List(ctx, req); err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"result": fmt.Sprint(response.Result),
		})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

}


func Serve() {

	r :=  gin.Default()

	r.GET("/nodes", listNodes)
	node := r.Group("/node")
	{
		node.POST("/", createNode)
		node.GET("/:node_id", describeNode)
		node.PUT("/:node_id", updateNode)
		node.DELETE("/:node_id", deleteNode)
		node.POST("/deploy/:node_id", deployNode)
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to Run server: %v ", err)
	}

}

