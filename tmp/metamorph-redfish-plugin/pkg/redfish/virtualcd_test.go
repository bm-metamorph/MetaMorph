package redfish

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"github.com/bm-metamorph/MetaMorph/pkg/db/models/node"
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
		bmhnode.ImageURL = "http://12.168.220.23:31180/4c4c4544-004a-5910-804d-c2c04f435032-ubuntu.iso"
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
	bmhnode := &BMHNode { node.CreateTestNode()}
	bmhnode.ImageURL = "http://12.168.220.23:31180/4c4c4544-004a-5910-804d-c2c04f435032-ubuntu.iso"

	err := bmhnode.DeployISO()
	assert.Equal(t,err,nil)


}

func TestEjectISO(t *testing.T){
	bmhnode := &BMHNode { node.CreateTestNode()}
	res  := bmhnode.EjectISO()
	assert.Equal(t, res,true)

}

func TestPowerOff(t *testing.T){
	bmhnode := &BMHNode { node.CreateTestNode()}
	res := bmhnode.PowerOff()
	assert.Equal(t, res, true)

}

func TestPowerOn(t *testing.T){
	bmhnode := &BMHNode { node.CreateTestNode()}
	res := bmhnode.PowerOn()
	assert.Equal(t, res, true)

}

func TestGetRedfishVersion(t *testing.T){
	bmhnode := &BMHNode { node.CreateTestNode()}
	rfvers := bmhnode.GetRedfishVersion()
	assert.Equal(t,rfvers,"1.4.0")

}
func TestGetGUUID(t *testing.T){
	bmhnode := &BMHNode { node.CreateTestNode()}
	uuid := bmhnode.GetGUUID()
	assert.Equal(t,uuid,"4c4c4544-004a-5910-804d-c2c04f435032")

}
