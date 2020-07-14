package redfish

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/bm-metamorph/MetaMorph/pkg/db/models/node"

)



func TestCleanVirtualDIskIfEExists(t *testing.T){
	bmhnode := &BMHNode { node.CreateTestNode()}
	res := bmhnode.CleanVirtualDIskIfEExists()
    assert.Equal(t, res,true)    

}
func TestCreateVirtualDIsks(t *testing.T){
	bmhnode := &BMHNode { node.CreateTestNode()}
	res := bmhnode.CreateVirtualDisks()
    assert.Equal(t, res,true)    

}