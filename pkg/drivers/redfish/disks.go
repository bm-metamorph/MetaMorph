package redfish

import (
	config "bitbucket.com/metamorph/pkg/config"
	"bitbucket.com/metamorph/pkg/db/models/node"
	"fmt"
	client "github.com/manojkva/go-redfish-API-Wrapper/pkg/redfishwrap/idrac"
)

type BMHNode node.Node

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

func (bmhnode *BMHNode) cleanVirtualDIskIfEExists() bool {
	redfishClient := getRedfishClient(bmhnode)
	result := redfishClient.CleanVirtualDisksIfAny(config.Get("idrac.systemID").(string),
		config.Get("idrac.controllerID").(string))
	return result
}

func (bmhnode *BMHNode) CreateVirtualDisks() bool {

	var result bool

	if !bmhnode.cleanVirtualDIskIfEExists() {
		return false
	}

	redfishClient := getRedfishClient(bmhnode)
	raidLevelMap := getSupportedRAIDLevels()

	for _, vd := range bmhnode.VirtualDisks {

		var diskIDs []string
		for _, disk := range vd.PhysicalDisks {
			diskIDs = append(diskIDs, disk.PhysicalDisk)
		}

		jobId := redfishClient.CreateVirtualDisk(config.Get("idrac.systemID").(string),
			config.Get("idrac.controllerID").(string), raidLevelMap[vd.RaidType], vd.DiskName, diskIDs)
		fmt.Printf("Job Id returned is %v\n", jobId)
		//check Job Status to decide on return value
		result = true

	}
	return result
}
