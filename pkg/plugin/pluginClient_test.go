package plugin

import (
	"github.com/bm-metamorph/MetaMorph/pkg/db/models/node"
	//	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadConfigFile(t *testing.T) {
	bmhnode := &BMHNode{node.CreateTestNode()}
	bmhnode.ReadConfigFile()

}

func TestDispense(t *testing.T) {
	bmhnode := &BMHNode{node.CreateTestNode()}
	bmhnode.ReadConfigFile()
	bmhnode.Dispense("getguuid")

}
