package redfish

import (
	config "bitbucket.com/metamorph/pkg/config"
	"bitbucket.com/metamorph/pkg/db/models/node"
	"fmt"
	client "github.com/manojkva/go-redfish-api-wrapper/pkg/redfishwrap/idrac"
)

type BMHNode struct {
	*node.Node
}

func getSupportedRAIDLevels() map[int]string {

	return map[int]string{
		1:  "Mirrored",
		5:  "StripedWithParity",
		10: "SpannedMirrors",
		50: "SpannedStripesWithParity",
	}
}

func getRedfishClient(bmhnode *BMHNode) client.IdracRedfishClient {
	redfishClient := client.IdracRedfishClient{
		Username: bmhnode.IPMIUser,
		Password: bmhnode.IPMIPassword,
		HostIP:   bmhnode.IPMIIP,
	}
	return redfishClient

}

func (bmhnode *BMHNode) CleanVirtualDIskIfEExists() bool {
	var result bool = false
	redfishClient := getRedfishClient(bmhnode)
	virtualdisklist, err := node.GetVirtualDisks(bmhnode.NodeUUID.String())
	if err != nil {
		fmt.Printf("Virtual disk list is empty with err %v\n", err)
		return false
	}
	for _, raiddisk := range virtualdisklist {

		result = redfishClient.CleanVirtualDisksIfAny(config.Get("idrac.systemID").(string), raiddisk.RaidController)
		if result == false {
			fmt.Printf("Failed to clean up Virtual Disk %v\n", raiddisk)
			return result
		}
	}

	return result
}

func (bmhnode *BMHNode) CreateVirtualDisks() bool {

	fmt.Printf("Inside Create Virtual Disk function\n")

	var result bool

	if !bmhnode.CleanVirtualDIskIfEExists() {
		return false
	}

	redfishClient := getRedfishClient(bmhnode)
	raidLevelMap := getSupportedRAIDLevels()

	virtualdisklist, err := node.GetVirtualDisks(bmhnode.NodeUUID.String())
	if err != nil {
		fmt.Printf("Virtual disk list is empty with err %v\n", err)
		return false
	}

	for _, vd := range virtualdisklist {

		var diskIDs []string
		physicaldisklist, err := node.GetPhysicalDisks(vd.ID)

		if err != nil {
			fmt.Printf("Failed to retrieve physical disks with error %v", err)
			return false
		}
		for _, disk := range physicaldisklist {
			diskIDs = append(diskIDs, disk.PhysicalDisk)
		}

		volumeType := raidLevelMap[vd.RaidType]

		jobId := redfishClient.CreateVirtualDisk(config.Get("idrac.systemID").(string),
			vd.RaidController, volumeType, vd.DiskName, diskIDs)

		fmt.Printf("Job Id returned is %v\n", jobId)
		//check Job Status to decide on return value
		if jobId != "" {
			result = redfishClient.CheckJobStatus(jobId)
		} else {
			result = false
		}

		if result != true {
			return result
		}

	}
	return result
}
