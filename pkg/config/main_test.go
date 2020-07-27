package  config


import (
	"testing"
	"fmt"
)


func TestConfig(t * testing.T){
	x := GetStringMapString( "plugins" )
	fmt.Printf("%v\n",x)

}
