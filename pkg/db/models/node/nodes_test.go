package node


import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGetBondParameters (t *testing.T){
	node  := CreateTestNode()
	bondParameters, _  := GetBondParameters(node.NodeUUID.String())
	assert.Equal(t,bondParameters.Mode,"802.3ad")
	assert.Equal(t,bondParameters.LacpRate,"fast")
}
func TestGetKvmPolicy(t *testing.T){
	node  := CreateTestNode()
	kvmPolicy, _  := GetKvmPolicy(node.NodeUUID.String())
	assert.Equal(t,kvmPolicy.CpuAllocation,"1:1")
	assert.Equal(t,kvmPolicy.CpuPinning,"enabled")
	assert.Equal(t,kvmPolicy.CpuHyperthreading,"enabled")
}