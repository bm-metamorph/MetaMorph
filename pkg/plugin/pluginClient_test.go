package plugin

import (
	"github.com/bm-metamorph/MetaMorph/pkg/db/models/node"
	"github.com/manojkva/metamorph-plugin/pkg/config"
	//	"github.com/stretchr/testify/assert"
	"fmt"
	"testing"
)

func TestReadConfigFile(t *testing.T) {
	bmhnode := &BMHNode{node.CreateTestNode()}
	bmhnode.ReadConfigFile()

}

func TestDispenseClientRequest(t *testing.T) {
        config.SetLoggerConfig("logger.pluginpath")
	bmhnode := &BMHNode{node.CreateTestNode()}
	err := bmhnode.ReadConfigFile()
	if err == nil{
	x, err := bmhnode.DispenseClientRequest("gethwinventory")
	if err == nil {
		fmt.Printf("%+v\n", (x.(map[string]string)))
	} else {
		fmt.Printf("Error %v\n", err)
	}
}else{
	fmt.Println("Failed to read Config file")
}


}
