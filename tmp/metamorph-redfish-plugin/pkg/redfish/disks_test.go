package redfish

import (
	"github.com/bm-metamorph/MetaMorph/pkg/db/models/node"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCleanVirtualDIskIfEExists(t *testing.T) {
	bmhnode := &BMHNode{node.CreateTestNode()}
	res := bmhnode.CleanVirtualDIskIfEExists()
	assert.Equal(t, res, true)

}
func TestConfigureRAID(t *testing.T) {
	bmhnode := &BMHNode{node.CreateTestNode()}
	err := bmhnode.ConfigureRAID()
	assert.Equal(t, err, nil)

}
