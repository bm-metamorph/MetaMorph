package idrac

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	//	"encoding/json"

	RFWrap "github.com/manojkva/go-redfish-api-wrapper/pkg/redfishwrap"
	redfish "opendev.org/airship/go-redfish/client"
	"go.uber.org/zap"
)

var RedfishSleepTimeSeconds int = 0
var JobCheckTimeoutMinutes int  = 0

func init() {
	RedfishSleepTimeSeconds, _ = strconv.Atoi(os.Getenv("REDFISH_SLEEPTIME_SECS"))
	JobCheckTimeoutMinutes,_ = strconv.Atoi(os.Getenv("REDFISH_JOBCHECKTIMEOUT_MTS"))
	if JobCheckTimeoutMinutes == 0{
		JobCheckTimeoutMinutes = 10 // five minutes assumed default
	}

}

type IdracRedfishClient struct {
	Username  string
	Password  string
	HostIP    string
	IDRAC_ver string
	Logger    *zap.Logger
}

func (a *IdracRedfishClient) createContext() context.Context {

	var auth = redfish.BasicAuth{UserName: a.Username,
		Password: a.Password,
	}
	ctx := context.WithValue(context.Background(), redfish.ContextBasicAuth, auth)
	ctx = context.WithValue(ctx, "logger", a.Logger)
	return ctx
}

func (a *IdracRedfishClient) UpgradeFirmware(filelocation string) bool {
	a.Logger.Info("UpgradeFirmware()")
	var imageURI string
	var err error
	var result bool = false

	ctx := a.createContext()

	httpPushURI := RFWrap.UpdateService(ctx, a.HostIP)

	a.Logger.Debug("HTTPURI retrieved :", zap.String("httpPushUri",fmt.Sprintf("%v", httpPushURI)))

	etag := RFWrap.GetETagHttpURI(ctx, a.HostIP)
	a.Logger.Debug("Etag revtreived :", zap.String("eTag",fmt.Sprintf("%v", etag)))
	if etag == "" {
		a.Logger.Error(fmt.Sprintf("Failed to extract ETAG..\n"))
		return false
	}
	if joblistEmpty, _ := a.IsJobListEmpty(); joblistEmpty == false {
		a.Logger.Error(fmt.Sprintf("Pending Jobs not empty.Hence returning\n"))
		return false
	}

	imageURI, err = RFWrap.HTTPUriDownload(ctx, a.HostIP, filelocation, etag)

	if (imageURI == "") || (err != nil) {
		a.Logger.Error(fmt.Sprintf("Failed to retrive ImageURI\n"))
		return false
	}
	a.Logger.Debug("ImageURI retrieved.", zap.String("imageURI",fmt.Sprintf("%v", imageURI)))

	jobID := RFWrap.SimpleUpdateRequest(ctx, a.HostIP, imageURI)

	a.Logger.Debug("JobId retreived", zap.String("JobID", fmt.Sprintf("%v", jobID)))
    //Check if Update succeeded. A scheduledl Job too is fine.
	result = a.CheckJobStatus(jobID, true)

	if result != false {
		if a.RebootServer(a.GetSystemID()) {
			result = a.CheckJobStatus(jobID,false)
		}
	}
	return result
}

func (a *IdracRedfishClient) GetPendingJobs() int {
	a.Logger.Info("GetPendingJobs()")
	ctx := a.createContext()
	statusCode, pendingJobCount := RFWrap.GetTaskList(ctx, a.HostIP)

	if statusCode == http.StatusOK {
		a.Logger.Debug("Pending Job count", zap.Int("Pending Job", pendingJobCount))
		return pendingJobCount
	}
	return 0
}

func (a *IdracRedfishClient) CheckJobStatus(jobId string, setScheduledasTrue bool) bool {
	a.Logger.Info("CheckJobStatus()")
	ctx := a.createContext()
	start := time.Now()
	var result bool = false

	if jobId == "" {
		a.Logger.Error("Job ID is null. Returing Failed Job Status")
		return false
	}

	for {

		statusCode, jobInfo := RFWrap.GetTask(ctx, a.HostIP, jobId)

		timeelapsedInMinutes := time.Since(start).Minutes()

		if (statusCode == 202) || (statusCode == 200) {
			a.Logger.Debug(fmt.Sprintf("HTTP  status OK"))

		} else {
			a.Logger.Error("HTTP Request Failed ",zap.String("JobID", jobId), zap.Int("HTTP Status", statusCode))
			return false
		}

		if timeelapsedInMinutes >= float64(JobCheckTimeoutMinutes) {
			a.Logger.Error(fmt.Sprintf("\n- FAIL: Timeout of %v minute has been hit, update job should of already been marked completed. Check the iDRAC job queue and LC logs to debug the issue\n",JobCheckTimeoutMinutes))
			return false
		} else if jobInfo.Messages != nil {
			if strings.Contains(jobInfo.Messages[0].Message, "Failed") {
				a.Logger.Error("Job Failed", zap.String("JobID", jobId))
				return false

			} else if strings.Contains(jobInfo.Messages[0].Message, "scheduled") && (setScheduledasTrue) {
				a.Logger.Debug(fmt.Sprintf("\n- PASS, job ID %s successfully marked as scheduled, powering on or rebooting the server to apply the update",jobId ))
				result = true
				break

			} else if strings.Contains(jobInfo.Messages[0].Message, "completed successfully") {
				//		fmt.Prinln("\n- PASS, job ID %s successfully marked as scheduled, powering on or rebooting the server to apply the update" % data[u"Id"] ")
				a.Logger.Debug(fmt.Sprintf("Job successfully marked as scheduled."), zap.String("JobID", jobId))
				result = true
				break
			}
		} else {
			a.Logger.Debug("Sleeping for 5 Secs", zap.String("JobID", jobId))
			time.Sleep(time.Second * 5)
			continue
		}
	}
	return result
}

func (a *IdracRedfishClient) RebootServer(systemID string) bool {
	a.Logger.Info("RebootServer()")

	ctx := a.createContext()

	//Systems/System.Embedded.1/Actions/ComputerSystem.Reset
	resetRequestBody := redfish.ResetRequestBody{ResetType: redfish.RESETTYPE_FORCE_RESTART}

	return RFWrap.ResetServer(ctx, a.HostIP, systemID, resetRequestBody)

}

func (a *IdracRedfishClient) PowerOn(systemID string) bool {
	a.Logger.Info("PowerOn()")
	ctx := a.createContext()
	resetRequestBody := redfish.ResetRequestBody{ResetType: redfish.RESETTYPE_ON}

	return RFWrap.ResetServer(ctx, a.HostIP, systemID, resetRequestBody)

}

func (a *IdracRedfishClient) PowerOff(systemID string) bool {
	a.Logger.Info("PowerOff()")
	ctx := a.createContext()
	resetRequestBody := redfish.ResetRequestBody{ResetType: redfish.RESETTYPE_GRACEFUL_SHUTDOWN}

	return RFWrap.ResetServer(ctx, a.HostIP, systemID, resetRequestBody)
}

func (a *IdracRedfishClient) GetVirtualMediaStatus(managerID string, media string) bool {
	a.Logger.Info("GetVirtualMediaStatus()")
	ctx := a.createContext()
	return RFWrap.GetVirtualMediaConnectedStatus(ctx, a.HostIP, managerID, media)
}

func (a *IdracRedfishClient) EjectISO(managerID string, media string) bool {
	a.Logger.Info("EjectISO()")
	ctx := a.createContext()
	return RFWrap.EjectVirtualMedia(ctx, a.HostIP, managerID, media)
}

func (a *IdracRedfishClient) SetOneTimeBoot(systemID string) bool {
	a.Logger.Info("SetOneTimeBoot()")
	ctx := a.createContext()
	computeSystem := redfish.ComputerSystem{Boot: redfish.Boot{BootSourceOverrideEnabled: redfish.BOOTSOURCEOVERRIDEENABLED_ONCE,
		BootSourceOverrideTarget: redfish.BOOTSOURCE_CD}}

	return RFWrap.SetSystem(ctx, a.HostIP, systemID, computeSystem)

}

func (a *IdracRedfishClient) InsertISO(managerID string, mediaID string, imageURL string) bool {
	a.Logger.Info("InsertISO()")

	ctx := a.createContext()

	if a.GetVirtualMediaStatus(managerID, mediaID) {
		a.Logger.Error(fmt.Sprintf("Exiting .. Already connected"),zap.String("ImageURL", imageURL))
		return false
	}
	insertMediaReqBody := redfish.InsertMediaRequestBody{
		Image: imageURL,
	}
	return RFWrap.InsertVirtualMedia(ctx, a.HostIP, managerID, mediaID, insertMediaReqBody)

}

func (a *IdracRedfishClient) GetVirtualDisks(systemID string, controllerID string) []string {
	a.Logger.Info("GetVirtualDisks()")

	ctx := a.createContext()
	idrefs := RFWrap.GetVolumes(ctx, a.HostIP, systemID, controllerID)
	if idrefs == nil {
		return nil
	}
	virtualDisks := []string{}
	for _, id := range idrefs {

		a.Logger.Debug(fmt.Sprintf("VirtualDisk Info %v\n", id.OdataId))
		vd := strings.Split(id.OdataId, "/")
		if vd != nil {
			virtualDisks = append(virtualDisks, vd[len(vd)-1])
		}
	}
	return virtualDisks

}

func (a *IdracRedfishClient) DeletVirtualDisk(systemID string, storageID string) string {
	a.Logger.Info("DeleteVirtualDisk()")

	if joblistEmpty, _ := a.IsJobListEmpty(); joblistEmpty != false {
		ctx := a.createContext()
		return RFWrap.DeleteVirtualDisk(ctx, a.HostIP, systemID, storageID)
	}
	return "" //failure send JobID = ""
}

func (a *IdracRedfishClient) CreateVirtualDisk(systemID string, controllerID string, volumeType string, name string, urilist []string) string {
	a.Logger.Info("CreateVirtualDisk()")
	if joblistEmpty, _ := a.IsJobListEmpty(); joblistEmpty != false {
		ctx := a.createContext()

		drives := []redfish.IdRef{}

		for _, uri := range urilist {
			driveinfo := fmt.Sprintf("/redfish/v1/Systems/%s/Storage/Drives/%s", systemID, uri)
			drives = append(drives, redfish.IdRef{OdataId: driveinfo})
		}

		createvirtualBodyReq := redfish.CreateVirtualDiskRequestBody{
			VolumeType: redfish.VolumeType(volumeType),
			Name:       name,
			Drives:     drives,
		}

		return RFWrap.CreateVirtualDisk(ctx, a.HostIP, systemID, controllerID, createvirtualBodyReq)
	}
	a.Logger.Error("Failed to Create Virtual Disks")
	return "" // failure send JobID  = ""
}

func (a *IdracRedfishClient) IsJobListEmpty() (bool, error) {
	a.Logger.Info("IsJobListEmpty()")

	if RedfishSleepTimeSeconds == 0 {
		RedfishSleepTimeSeconds = 100 // if environment variable is not set
	}
	timeout := time.NewTimer(time.Second * time.Duration(RedfishSleepTimeSeconds))
	ticker := time.NewTicker(3000 * time.Millisecond)
	for {
		select {
		case <-timeout.C:
			timeout.Stop()
			ticker.Stop()
			return false, errors.New("Getting Pending Jobs timed out")
		case <-ticker.C:
			a.Logger.Debug(fmt.Sprintf("Checking Pending Jobs .."))
			pendingJobs := a.GetPendingJobs()
			if pendingJobs == 0 {
				timeout.Stop()
				ticker.Stop()
				return true, nil
			}
		}
	}
	timeout.Stop()
	ticker.Stop()
	a.Logger.Error("Failed to Get Pending Jobs")
	return false, errors.New("Failed to Get Pending Jobs")

}

func (a *IdracRedfishClient) CleanVirtualDisksIfAny(systemID string, controllerID string) bool {
	a.Logger.Info("CleanVirtualDisksIfAny()")

	var result bool = false

	// Get the list of VirtualDisks
	virtualDisks := a.GetVirtualDisks(systemID, controllerID)
	totalvirtualDisks := len(virtualDisks)
	var countofVDdeleted int = 0
	// for testing skip the OS Disk
	//virtualDisks = virtualDisks[1:]
	if totalvirtualDisks == 0 {
		a.Logger.Debug(fmt.Sprintf("No existing RAID disks found"))
		result = true
	} else {
		for _, vd := range virtualDisks {
			jobid := a.DeletVirtualDisk(systemID, vd)
			a.Logger.Debug(fmt.Sprintf("Delete Job ID %v\n", jobid))
			result = a.CheckJobStatus(jobid,false)

			if result == false {
				a.Logger.Error(fmt.Sprintf("Failed to delete virtual disk %v\n", vd))
				return result
			}
			//	time.Sleep(time.Second * time.Duration(RedfishSleepTimeSeconds)) //Sleep in between calls
			countofVDdeleted += 1

		}
	}
	if countofVDdeleted != totalvirtualDisks {
		a.Logger.Error("Number of VirtualDisk deleted is not equal to Total Virtual Disks")
		result = false
	}

	return result
}

func (a *IdracRedfishClient) GetNodeUUID(systemID string) (string, bool) {
	a.Logger.Info("GetNodeUUID()")

	ctx := a.createContext()
	computerSystem, _ := RFWrap.GetSystem(ctx, a.HostIP, systemID)

	if computerSystem != nil {
		return computerSystem.UUID, true
	}
	a.Logger.Error("Failed to retrieve Node UUID")
	return "", false
}

func (a *IdracRedfishClient) GetPowerStatus(systemID string) bool {
	a.Logger.Info("GetPowerStatus()")

	ctx := a.createContext()
	computerSystem, _ := RFWrap.GetSystem(ctx, a.HostIP, systemID)

	if computerSystem != nil {
		if computerSystem.PowerState == "On" {
			return true
		}
	}
	a.Logger.Error("Failed to get Power Status")
	return false
}

func (a *IdracRedfishClient) GetManagerID() string {
	a.Logger.Info("GetManagerID()")
	ctx := a.createContext()
	managerList := RFWrap.ListManagers(ctx, a.HostIP)
	if managerList == nil {
		a.Logger.Error(fmt.Sprintf("Failed to retreive manager ID"))
		return ""
	}
	a.Logger.Debug(fmt.Sprintf("%+v", managerList))
	return managerList[0]
}
func (a *IdracRedfishClient) GetSystemID() string {
	a.Logger.Info("GetSystemID()")
	ctx := a.createContext()
	systemList := RFWrap.ListSystems(ctx, a.HostIP)
	if systemList == nil {
		a.Logger.Error(fmt.Sprintf("Failed to retreive system ID"))
		return ""
	}
	a.Logger.Debug(fmt.Sprintf("%+v", systemList))
	return systemList[0]
}

func (a *IdracRedfishClient) GetRedfishVer() string {
	a.Logger.Info("GetRedfishVer()")
	ctx := a.createContext()
	root := RFWrap.GetRoot(ctx, a.HostIP)
	if root == nil {
		a.Logger.Error(fmt.Sprintf("Failed to retreive RedfishVersion"))
		return ""
	}
	redfishVersion := root.RedfishVersion
	a.Logger.Debug(fmt.Sprintf("Redfish Version : %+v", redfishVersion))
	return redfishVersion
}

func (a *IdracRedfishClient) GetFirmwareDetails(firmwarename string) (name string, version string, updateable bool) {
	a.Logger.Info("GetFirmwareDetails()")


	ctx := a.createContext()

	firmwareInv := RFWrap.GetFirwareInventory(ctx, a.HostIP)
	if firmwareInv == nil {
		a.Logger.Error(fmt.Sprintf("Failed to retreive FirmwareInventory"))
		return "", "", false
	}
	for _, id := range firmwareInv.Members{
		var softwareId string

		a.Logger.Debug("OdataId ", zap.String("OdataID",fmt.Sprintf("%v", id.OdataId)))
		fd := strings.Split(id.OdataId, "/")
		if fd != nil {
			softwareId = fd[len(fd)-1]
			fmt.Sprintf("Software Id %v\n",softwareId)

			softwareInv  := RFWrap.GetSoftwareInventory(ctx,a.HostIP,softwareId)
			a.Logger.Debug(fmt.Sprintf("Software  Inv : %+v\n",softwareInv))
			name = softwareInv.Name
			version = *softwareInv.Version
			updateable = *softwareInv.Updateable
			a.Logger.Debug(fmt.Sprintf("%+v,%+v, %+v", name,version,updateable))

			if strings.Contains( strings.ToLower(name), strings.ToLower(firmwarename)){
				return name,version,updateable
			}

		}
	}
	a.Logger.Error("Failed to retrieve Firmware Inventory")
	return "", "", false

}
