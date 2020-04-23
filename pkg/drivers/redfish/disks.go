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
	for _, raiddisk := range bmhnode.VirtualDisks {

		result = redfishClient.CleanVirtualDisksIfAny(config.Get("idrac.systemID").(string), raiddisk.RaidController)
	}

	return result
}

func (bmhnode *BMHNode) CreateVirtualDisks() bool {

	var result bool

	if !bmhnode.CleanVirtualDIskIfEExists() {
		return false
	}

	redfishClient := getRedfishClient(bmhnode)
	raidLevelMap := getSupportedRAIDLevels()

	for _, vd := range bmhnode.VirtualDisks {

		var diskIDs []string
		for _, disk := range vd.PhysicalDisks {
			diskIDs = append(diskIDs, disk.PhysicalDisk)
		}

		volumeType := raidLevelMap[vd.RaidType]

		jobId := redfishClient.CreateVirtualDisk(config.Get("idrac.systemID").(string),
			vd.RaidController, volumeType, vd.DiskName, diskIDs)
			
		fmt.Printf("Job Id returned is %v\n", jobId)
		//check Job Status to decide on return value
        result  = redfishClient.CheckJobStatus(jobId)

		if result != true {
			return result
		}

	}
	return result
}
