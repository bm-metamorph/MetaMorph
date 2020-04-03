package main

import(
	"fmt"
	"bitbucket.com/metamorph/proto"
	"github.com/gin-gonic/gin"
	"net/http"
	//"strconv"
	"log"


	"google.golang.org/grpc"
	//"google.golang.org/grpc/reflection"
)






func main() {

	conn, err := grpc.Dial( "localhost:4040", grpc.WithInsecure() )
	if err != nil {
		panic(err)
	}

	client := proto.NewNodeServiceClient(conn) 
	g :=  gin.Default()

	g.POST("/node", func(ctx *gin.Context){

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
	})

	g.GET("/node/:node_id", func(ctx *gin.Context){
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

	})


	g.POST("/node/deploy/:node_id", func(ctx *gin.Context){
		node_id := ctx.Param("node_id")
		//fmt.Println(node_id)
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

	})

	if err := g.Run(":8080"); err != nil {
		log.Fatalf("Failed to Run server: %v ", err)
	}

}

