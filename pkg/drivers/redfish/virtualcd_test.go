package redfish

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReboot(t *testing.T) {
	bmhnode := createTestNode()
	res := bmhnode.Reboot()
	assert.Equal(t, res, true)
}

func TestISOInstallation(t *testing.T) {
	var res bool = false

	bmhnode := createTestNode()
	//Step 1 Eject Existing ISO
	res = bmhnode.EjectISO()
	res = true
	if res != false {
		//Step 2 Insert Ubuntu ISO
		bmhnode.ImageURL = "http://32.68.220.23:31180/a451dcb7-9a17-45a8-8915-f5ab0a175cf6-ubuntu.iso"
		res = bmhnode.InsertISO()
		res  = true
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
