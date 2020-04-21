package isogen

import (
	"fmt"
	"io/ioutil"
	"testing"
	config "bitbucket.com/metamorph/pkg/config"
//	"bitbucket.com/metamorph/pkg/db/models/node"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

)

func createTestNode() *BMHNode{
	data, _ := ioutil.ReadFile(config.Get("testing.inputfile").(string))
	var bmhnode *BMHNode  = new(BMHNode)
	UUID, _ := uuid.NewRandom()
	_ = json.Unmarshal(data, &bmhnode)
	bmhnode.NodeUUID = UUID
	return bmhnode

}

func TestgetDiskSpaceinMB(t *testing.T) {
	dspace, maxdspace, _ := getDiskSpaceinMB(">300g")
	fmt.Printf("%s, %s", dspace, maxdspace)
}

func TestCreateGrubfile(t *testing.T){
	bmhnode  := createTestNode()
	bmhnode.CreateFileFromTemplate("/tmp","grub")
	assert.FileExists(t,"/tmp/grub.conf","")

}

func TestCreatePreseedfile(t *testing.T){
	bmhnode  := createTestNode()
	bmhnode.CreateFileFromTemplate("/tmp","preseed")
	assert.FileExists(t,"/tmp/preseed/hwe-ubuntu-server.seed","")
}

func TestCreateNetplanfile(t *testing.T){
	bmhnode  := createTestNode()
	bmhnode.CreateFileFromTemplate("/tmp","netplan")
	assert.FileExists(t,"/tmp/50-cloud-init.yaml")
}
func TestCreateInitfile(t *testing.T){
	bmhnode  := createTestNode()
	bmhnode.CreateFileFromTemplate("/tmp","init")
	assert.FileExists(t,"init")
}
