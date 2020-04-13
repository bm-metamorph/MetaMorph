package controller


import (
	"testing"
	"fmt"
	"bitbucket.com/metamorph/pkg/db/models/node"
	"github.com/google/uuid"
)


func TestMetamorphFSM(t *testing.T ){
	var nodelist = []node.Node {
		{ Name: "node1", NodeUUID: uuid.New(), State: NEW},
		{ Name: "node2", NodeUUID:  uuid.New(), State: READY},
		{ Name: "node3", NodeUUID:  uuid.New(), State: SETUPREADY},
		{ Name: "node4", NodeUUID:  uuid.New(), State: DEPLOYED},
		{ Name: "node5", NodeUUID:  uuid.New(), State: FAILED},
	}
	fmt.Println(nodelist)


}