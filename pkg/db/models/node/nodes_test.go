package node

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBondParameters(t *testing.T) {
	node := CreateTestNode()
	bondParameters, _ := GetBondParameters(node.NodeUUID.String())
	assert.Equal(t, len(bondParameters),6)
}
func TestGetKvmPolicy(t *testing.T) {
	node := CreateTestNode()
	kvmPolicy, _ := GetKvmPolicy(node.NodeUUID.String())
	assert.Equal(t, kvmPolicy.CpuAllocation, "1:1")
	assert.Equal(t, kvmPolicy.CpuPinning, "enabled")
	assert.Equal(t, kvmPolicy.CpuHyperthreading, "enabled")
}
func TestGetFilesystem(t *testing.T) {
	node := CreateTestNode()
	partitions, _ := GetPartitions(node.NodeUUID.String())
	for index, part := range partitions {
		filesystem, _ := GetFilesystem(part.ID)
		partitions[index].Filesystem = *filesystem
	}
	fmt.Printf("%v", partitions[0].Filesystem)
}

func  TestGetPlugins(t * testing.T){
	node := CreateTestNode()
	plugins, _ := GetPlugins(node.NodeUUID.String())
	fmt.Printf("%+v", plugins)
	apis,_ := GetPluginAPIs(plugins.ID)
	fmt.Printf("%+v", apis)
}
