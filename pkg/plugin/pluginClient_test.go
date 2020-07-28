package plugin

import (
	"github.com/bm-metamorph/MetaMorph/pkg/db/models/node"
	//	"github.com/stretchr/testify/assert"
	"fmt"
	"testing"
)

func TestReadConfigFile(t *testing.T) {
	bmhnode := &BMHNode{node.CreateTestNode()}
	bmhnode.ReadConfigFile()

}

func TestDispenseClientRequest(t *testing.T) {
	bmhnode := &BMHNode{node.CreateTestNode()}
	bmhnode.ReadConfigFile()
	x, err := bmhnode.DispenseClientRequest("gethwinventory")
	if err == nil {
		fmt.Printf("%v\n", (x.(map[string]string)))
	} else {
		fmt.Printf("Error %v\n", err)
	}

}
