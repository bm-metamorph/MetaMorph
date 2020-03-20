package main

import (
	"fmt"
	"github.com/metamorph/pkg/apis"
	"github.com/metamorph/pkg/controller"
	"github.com/metamorph/pkg/drivers/isogen"
)


func main() {
	fmt.Println("hello world")
	apis.GetVersion()
	controller.ControllerInit()
	isogen.GenerateISO()
}
