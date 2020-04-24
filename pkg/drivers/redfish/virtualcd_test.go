package redfish

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"bitbucket.com/metamorph/pkg/db/models/node"
)

func TestReboot(t *testing.T) {
	bmhnode := &BMHNode { node.CreateTestNode()}
	res := bmhnode.Reboot()
	assert.Equal(t, res, true)
}

func TestISOInstallation(t *testing.T) {
	var res bool = false
	bmhnode := &BMHNode { node.CreateTestNode()}

	//Step 1 Eject Existing ISO
	res = bmhnode.EjectISO()
	if res != false {
		//Step 2 Insert Ubuntu ISO
		bmhnode.ImageURL = "http://32.68.220.23:31180/a451dcb7-9a17-45a8-8915-f5ab0a175cf6-ubuntu.iso"
		res = bmhnode.InsertISO()
		if res != false {
			//Step 3 Set Onetime boot to CD ROM
			res = bmhnode.SetOneTimeBoot()
			if res != false {
				//Step 4 Reboot

				res = bmhnode.Reboot()

			}
		}
	}
	assert.Equal(t, res, true)
}

func TestISODeploy(t *testing.T){
	var res bool
	bmhnode := &BMHNode { node.CreateTestNode()}
	bmhnode.ImageURL = "http://32.68.220.23:31180/a451dcb7-9a17-45a8-8915-f5ab0a175cf6-ubuntu.iso"

	res = bmhnode.DeployISO()
	assert.Equal(t,res,true)


}

func TestEjectISO(t *testing.T){
	bmhnode := &BMHNode { node.CreateTestNode()}
	res  := bmhnode.EjectISO()
	assert.Equal(t, res,true)

}