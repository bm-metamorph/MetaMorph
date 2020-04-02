package main

import (
	 "github.com/gin-gonic/gin"
	 "fmt"
)


func createNode(c *gin.Context){

	c.JSON(200, gin.H{
		"message": "Created node",
	})
}

func updateNode(c *gin.Context){

	c.JSON(200, gin.H{
		"message": "Updated node",
	})
}


func describeNode(c *gin.Context){
	nodeID := c.Param("node_id")
	result := fmt.Sprintf("Describe  node %s", nodeID)
	c.JSON(200, gin.H{
		"message": result,
	})
}


func deleteNode(c *gin.Context){

	c.JSON(200, gin.H{
		"message": "deleted node",
	})
}


func deployNode(c *gin.Context){

	c.JSON(200, gin.H{
		"message": "deployed node",
	})
}


func main(){
	r := gin.Default()

	node := r.Group("/node")
	{
		node.POST("/", createNode)
		node.GET("/:node_id", describeNode)
		node.PUT("/", updateNode)
		node.DELETE("/", deleteNode)
		node.POST("/deploy/:node_id", deployNode)
	}

	r.GET("/nodes", func(c *gin.Context){

		c.JSON(200, gin.H{
			"message": "Node list",
		})
	})
	r.Run()
}