package isogen

import (
	"fmt"
	"testing"
	"bitbucket.com/metamorph/pkg/db/models/node"
	"github.com/stretchr/testify/assert"

)

func TestgetDiskSpaceinMB(t *testing.T) {
	dspace, maxdspace, _ := getDiskSpaceinMB(">300g")
	fmt.Printf("%s, %s", dspace, maxdspace)
}

func TestCreateGrubfile(t *testing.T){
	bmhnode := &BMHNode{node.CreateTestNode()}
	bmhnode.CreateFileFromTemplate("/tmp","grub")
	assert.FileExists(t,"/tmp/grub.conf","")

}

func TestCreatePreseedfile(t *testing.T){
	bmhnode := &BMHNode{node.CreateTestNode()}
	bmhnode.CreateFileFromTemplate("/tmp","preseed")
	assert.FileExists(t,"/tmp/preseed/hwe-ubuntu-server.seed","")
}

func TestCreateNetplanfile(t *testing.T){
	bmhnode := &BMHNode{node.CreateTestNode()}
	bmhnode.CreateFileFromTemplate("/tmp","netplan")
	assert.FileExists(t,"/tmp/50-cloud-init.yaml")
}
func TestCreateInitfile(t *testing.T){
	bmhnode := &BMHNode{node.CreateTestNode()}
	bmhnode.CreateFileFromTemplate("/tmp","init")
	assert.FileExists(t,"init")
}
