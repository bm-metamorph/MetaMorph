package  config


import (
	"testing"
	"fmt"
)


func TestConfig(t * testing.T){
	x := Get( "idrac.systemID" )
	fmt.Println(x)

}