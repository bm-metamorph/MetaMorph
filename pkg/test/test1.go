package main
import (
    "github.com/gin-gonic/gin"
     "log"
//    "fmt"
)


func main() {
    r := gin.New()

    r.POST("/test", func(c *gin.Context) {
        // Body disappear on controller
        NrawBody, _ := c.GetRawData()
        log.Println("second", string(NrawBody))
    })

    // Listen and serve on 0.0.0.0:8080
    r.Run(":810")
}
