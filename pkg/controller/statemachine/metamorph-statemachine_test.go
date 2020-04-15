package controller

import (
	"bitbucket.com/metamorph/pkg/db/models/node"
	"fmt"
	"github.com/google/uuid"
	//	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)
//Use mockery -name nodedb -inpkg .
func TestMetamorphFSM(t *testing.T) {

	fmt.Println("Test Case")
	t.Log("Running testcase")

	nodeDB := &mockNodedb{}

	nodeDB.On("GetNodes").Return(
		[]node.Node{
			{Name: "node1", NodeUUID: uuid.New(), State: NEW},
			{Name: "node2", NodeUUID: uuid.New(), State: READY},
			{Name: "node3", NodeUUID: uuid.New(), State: SETUPREADY},
			{Name: "node4", NodeUUID: uuid.New(), State: DEPLOYED},
			//		{Name: "node5", NodeUUID: uuid.New(), State: FAILED},
		}, nil).WaitUntil(time.After(5 * time.Second))

	handler := &DBHandler{db: nodeDB}

	handler.StartMetamorphFSM(true)

}
