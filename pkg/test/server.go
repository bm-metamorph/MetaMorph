package main

import (
    "log"
    "io/ioutil"
    "bytes"
    "github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        var bodyBytes []byte
        if c.Request.Body != nil {
          bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
        }
       c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

        log.Println("first", string(bodyBytes))
    }
}

func main() {
    r := gin.New()
    r.Use(Logger())

    r.POST("/test", func(c *gin.Context) {
        // Body disappear on controller
        NrawBody, _ := c.GetRawData()
        log.Println("second", string(NrawBody))
    })

    // Listen and serve on 0.0.0.0:8080
    r.Run(":81")
}
