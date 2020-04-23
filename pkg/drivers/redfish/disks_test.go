package redfish

import (
	"testing"
	"io/ioutil"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"bitbucket.com/metamorph/pkg/config"

)

func createTestNode() *BMHNode{
	data, _ := ioutil.ReadFile(config.Get("testing.inputfile").(string))
	var bmhnode *BMHNode  = new(BMHNode)
	UUID, _ := uuid.NewRandom()
	_ = json.Unmarshal(data, &bmhnode)
	bmhnode.NodeUUID = UUID
	return bmhnode

}


func TestCleanVirtualDIskIfEExists(t *testing.T){
	bmhnode := createTestNode()
	res := bmhnode.CleanVirtualDIskIfEExists()
    assert.Equal(t, res,true)    

}
func TestCreateVirtualDIsks(t *testing.T){
	bmhnode := createTestNode()
	res := bmhnode.CreateVirtualDisks()
    assert.Equal(t, res,true)    

}