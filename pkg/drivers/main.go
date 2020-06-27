package main

import (
	"bitbucket.com/metamorph/pkg/db/models/node"
	"bitbucket.com/metamorph/pkg/drivers/redfish"
	"fmt"
//	"os"
)

func checkISOFuntionality(){
        var res bool

	bmhnode := &redfish.BMHNode{node.CreateTestNode()}
	bmhnode.ImageURL = "http://test/xyz.iso"
	bmhnode.RedfishVersion = bmhnode.GetRedfishVersion()

	fmt.Scanln()

	bmhnode.RedfishManagerID = bmhnode.GetManagerID()
	fmt.Scanln()

	bmhnode.RedfishSystemID = bmhnode.GetSystemID()
	fmt.Scanln()
/*
	fmt.Println("Insert ISO ")

	res := bmhnode.InsertISO()
	if !res {
		fmt.Println("Failed to insert ISO")
		os.Exit(1)
	}
	fmt.Scanln()
	fmt.Println("Check Status of virtual media")
	bmhnode.GetVirtualMediaStatus()
	fmt.Scanln()
	fmt.Println("Eject ISO")
	res = bmhnode.EjectISO()
	if !res {
		fmt.Println("Failed to eject ISO")
		os.Exit(1)
	}
	fmt.Scanln()
	fmt.Println("Check Status of virtual media")
	bmhnode.GetVirtualMediaStatus()
*/
	fmt.Scanln()
	fmt.Println("Set Onetime Boot the server")
	res = bmhnode.SetOneTimeBoot()
	if !res {
		fmt.Println("Failed to set onetime boot")
	}

	fmt.Scanln()
/*
	fmt.Println("Reboot the server")
	res = bmhnode.Reboot()
	if !res {
		fmt.Println("Reboot server failed")

	}
*/
}

func TestFirmwareUpdate(){
 var  filepath string
 var res bool
 fmt.Println("Provide the absolute path of the firmware file")
 fmt.Scanln(&filepath)
 fmt.Println(filepath)
 bmhnode := &redfish.BMHNode{node.CreateTestNode()}
 res = bmhnode.UpgradeFirmware(filepath)
 if res == false{
	 fmt.Println("Failed to upgrade")
 }

}

// For Testing Purpose only !!!

func main() {
	TestFirmwareUpdate()
}
