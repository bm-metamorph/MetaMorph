package isogen

import (
	"github.com/bm-metamorph/MetaMorph/pkg/db/models/node"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetDiskSpaceinMB(t *testing.T) {
	dspace, maxdspace, _ := getDiskSpaceinMB(">300g")
	fmt.Printf("%s, %s", dspace, maxdspace)
}

func TestCreateGrubfile(t *testing.T) {
	bmhnode := &BMHNode{node.CreateTestNode()}
	bmhnode.CreateFileFromTemplate("/tmp", "grub")
	assert.FileExists(t, "/tmp/grub.conf", "")

}

func TestCreatePreseedfile(t *testing.T) {
	bmhnode := &BMHNode{node.CreateTestNode()}
	bmhnode.CreatePressedFileFromTemplate("/tmp", "preseed")
	assert.FileExists(t, "/tmp/preseed/hwe-ubuntu-server.seed", "")
}

func TestCreateNetplanfile(t *testing.T) {
	bmhnode := &BMHNode{node.CreateTestNode()}
	bmhnode.CreateNetplanFileFromTemplate("/tmp", "netplan")
	assert.FileExists(t, "/tmp/50-cloud-init.yaml")
}

func TestCreateNetplanFileFromString(t *testing.T){
       bmhnode := &BMHNode{node.CreateTestNode()}
       bmhnode.CreateNetplanFileFromString("/tmp","netplan")
        assert.FileExists(t, "/tmp/50-cloud-init.yaml")

}
func TestCreateInitfile(t *testing.T) {
	bmhnode := &BMHNode{node.CreateTestNode()}
	bmhnode.CreateFileFromTemplate("/tmp", "init")
	assert.FileExists(t, "/tmp/init.sh")
}
