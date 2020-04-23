package redfish

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"bitbucket.com/metamorph/pkg/db/models/node"

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