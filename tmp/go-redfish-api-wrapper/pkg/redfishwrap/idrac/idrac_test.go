package idrac

import (
	"fmt"
	//"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var client = &IdracRedfishClient{
	Username: "root",
	Password: "Abc.1234",
	HostIP:   "",
}

func TestUpgradeFirmware(t *testing.T) {
	filelocation := "/home/test/workspace/iDRAC-with-Lifecycle-Controller_Firmware_NKGJW_WN64_3.31.31.31_A00.EXE"
	client.UpgradeFirmware(filelocation)

}

func TestCheckJobStatus(t *testing.T) {
	jobId := ""
	client.CheckJobStatus(jobId)
}

func TestGetPendingTasks(t *testing.T) {
	numberOfJobs := client.GetPendingJobs()
	fmt.Println(numberOfJobs)
	assert.Equal(t, 0, numberOfJobs)

}

func TestCheckPendingJobs(t *testing.T) {

	anyjobs, err := client.IsJobListEmpty()
	assert.Equal(t, true, anyjobs)
	fmt.Printf("%v", err)
}

func TestGetVirtualDisks(t *testing.T) {
	systemID := "System.Embedded.1"
	controllerID := "RAID.Slot.6-1"
	client.GetVirtualDisks(systemID, controllerID)

}

func TestDeleteVirtualDisk(t *testing.T) {
	systemID := "System.Embedded.1"
	storageID := "Disk.Virtual.1:RAID.Slot.6-1"
	jobid := client.DeletVirtualDisk(systemID, storageID)
	t.Logf("Job ID %v", jobid)
	res := client.CheckJobStatus(jobid)
	assert.Equal(t, res, true)
}

func TestCleanVirtualDisksIfAny(t *testing.T) {
	systemID := "System.Embedded.1"
	controllerID := "RAID.Slot.6-1"
	client.CleanVirtualDisksIfAny(systemID, controllerID)

}

/*
name: ephemeral
          #          raid-type: 1
          #          disk:
          #            - Disk.Bay.8:Enclosure.Internal.0-1:RAID.Slot.6-1
          #            - Disk.Bay.9:Enclosure.Internal.0-1:RAID.Slot.6-1
*/

func TestCreateVirtualDisk(t *testing.T) {
	systemID := "System.Embedded.1"
	controllerID := "RAID.Slot.6-1"
	volumeType := "Mirrored"
	name := "ephemeral-1"
	drives := []string{"Disk.Bay.8:Enclosure.Internal.0-1:RAID.Slot.6-1",
		"Disk.Bay.9:Enclosure.Internal.0-1:RAID.Slot.6-1"}
	jobid := client.CreateVirtualDisk(systemID, controllerID, volumeType, name, drives)
	t.Logf("Job ID %v", jobid)
	res := client.CheckJobStatus(jobid)
	t.Logf("%v", res)
	assert.Equal(t, res, true)

}

func TestGetNodeUUID(t *testing.T) {
	systemID := "System.Embedded.1"
	uuid, _ := client.GetNodeUUID(systemID)
	t.Logf("UUID %v", uuid)

}
func TestGetPowerStatus(t *testing.T) {
	systemID := "System.Embedded.1"
	result := client.GetPowerStatus(systemID)
	t.Logf("Result %v", result)
	assert.Equal(t, result, true)

}

func TestPowerOff(t *testing.T) {
	systemID := "System.Embedded.1"
	result := client.PowerOff(systemID)
	assert.Equal(t, result, true)
}
func TestPowerOn(t *testing.T) {
	systemID := "System.Embedded.1"
	result := client.PowerOn(systemID)
	assert.Equal(t, result, true)
}
func TestRebootServer(t *testing.T) {
	systemID := "System.Embedded.1"
	result := client.RebootServer(systemID)
	assert.Equal(t, result, true)
}

func TestEjectISO(t *testing.T) {
	managerID := "iDRAC.Embedded.1"
	res := client.EjectISO(managerID, "CD")
	assert.Equal(t, res, true)

}
func TestInsertCD(t *testing.T) {

	managerID := "iDRAC.Embedded.1"
	imageURL := "http://32.168.220.23:31180/a451dcb7-9a17-45a8-8915-f5ab0a175cf6-ubuntu.iso"
	res := client.InsertISO(managerID, "CD", imageURL)
	assert.Equal(t, res, true)
}

func TestGetManagerID(t *testing.T) {
	managerId := client.GetManagerID()
	fmt.Printf("%+v", managerId)
}

func TestGetSystemID(t *testing.T) {
	systemId := client.GetSystemID()
	fmt.Printf("%+v", systemId)
}

func TestGetRedfishVer(t *testing.T) {
	redfishVer := client.GetRedfishVer()
	fmt.Printf("%+v", redfishVer)
}
