package redfish

import (
	"github.com/bm-metamorph/MetaMorph/pkg/db/models/node"
	"fmt"
	client "github.com/manojkva/go-redfish-api-wrapper/pkg/redfishwrap/idrac"
	logger "github.com/manojkva/metamorph-plugin/pkg/logger"
	"go.uber.org/zap"
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
		Logger:   logger.Log,
	}
	return redfishClient

}

func (bmhnode *BMHNode) CleanVirtualDIskIfEExists() bool {
	logger.Log.Info("CleanVirtualDIskIfEExists()")
	var result bool = false
	redfishClient := getRedfishClient(bmhnode)
	virtualdisklist, err := node.GetVirtualDisks(bmhnode.NodeUUID.String())
	if err != nil {
		logger.Log.Error(fmt.Sprintf("Virtual disk list is empty\n", zap.Error (err)))
		return false
	}
	for _, raiddisk := range virtualdisklist {

		result = redfishClient.CleanVirtualDisksIfAny(bmhnode.RedfishSystemID, raiddisk.RaidController)
		if result == false {
			logger.Log.Error(fmt.Sprintf("Failed to clean up Virtual Disk %v\n", raiddisk))
			return result
		}
	}

	return result
}

func (bmhnode *BMHNode) ConfigureRAID() error {
	logger.Log.Info("ConfigureRAID()")


	var result bool

	if !bmhnode.CleanVirtualDIskIfEExists() {
		logger.Log.Error("Failed to delete existing Virtual Disk")
	        return  fmt.Errorf("Failed to delete existing Virtual Disk")
	}

	redfishClient := getRedfishClient(bmhnode)
	raidLevelMap := getSupportedRAIDLevels()

	virtualdisklist, err := node.GetVirtualDisks(bmhnode.NodeUUID.String())
	if err != nil {
		logger.Log.Error("Virtual disk list is empty",zap.Error(err))
		return fmt.Errorf("Virtual disk list is empty with err %v", err)
	}

	for _, vd := range virtualdisklist {

		var diskIDs []string
		physicaldisklist, err := node.GetPhysicalDisks(vd.ID)

		if err != nil {
			logger.Log.Error("Failed to retrieve Physical Disk", zap.Error(err))
			return fmt.Errorf("Failed to retrieve physical disks with error %v", err)
		}
		for _, disk := range physicaldisklist {
			diskIDs = append(diskIDs, disk.PhysicalDisk)
		}

		volumeType := raidLevelMap[vd.RaidType]

		jobId := redfishClient.CreateVirtualDisk(bmhnode.RedfishSystemID,
			vd.RaidController, volumeType, vd.DiskName, diskIDs)

		logger.Log.Debug(fmt.Sprintf("Job Id returned is %v\n", jobId))
		//check Job Status to decide on return value
		if jobId != "" {
			result = redfishClient.CheckJobStatus(jobId,false)
		} else {
			logger.Log.Error("Failed to retrieve JobID")
			result = false
		}

		if result != true {
			return  fmt.Errorf("Failed to retrieve Job Status for Job ID %v", jobId)
		}

	}
	return nil
}
