package redfish

import (
	"fmt"
	"github.com/manojkva/metamorph-plugin/pkg/logger"
	client "github.com/manojkva/go-redfish-api-wrapper/pkg/redfishwrap/idrac"
	"go.uber.org/zap"
	"os"
	"strconv"
	"time"
	"errors"
)

var PowerStateChangeTimeoutSeconds int = 0

func init() {
	PowerStateChangeTimeoutSeconds, _ = strconv.Atoi(os.Getenv("METAMORPH_POWERCHANGE_TIMEOUT"))
}


func (bmhnode *BMHNode) GetVirtualMediaStatus() bool {
	var result bool = false
	redfishClient := getRedfishClient(bmhnode)
	if bmhnode.RedfishVersion == "1.0.0" { //HP workaround only for ILO4
		result = redfishClient.GetVirtualMediaStatus(bmhnode.RedfishManagerID, "2")
	} else {
		result = redfishClient.GetVirtualMediaStatus(bmhnode.RedfishManagerID, "CD")
	}
	return result
}

func (bmhnode *BMHNode) InsertISO() bool {
	if bmhnode.GetVirtualMediaStatus() == true {
		logger.Log.Error(fmt.Sprintf("Skipping Iso Insert. CD already Attached\n"))
		return false
	}

	redfishClient := getRedfishClient(bmhnode)
	logger.Log.Debug(fmt.Sprintf("Image URL to be inserted = %v", bmhnode.ImageURL))
	result := false
	for retryCount := 0; ; retryCount++ {
		if bmhnode.RedfishVersion == "1.0.0" {
			result = bmhnode.InsertISOILO4(bmhnode.RedfishManagerID)

		} else {
			result = redfishClient.InsertISO(bmhnode.RedfishManagerID, "CD", bmhnode.ImageURL)
		}
		if result == true {
			break
		}
		if retryCount >= 2 {
			break
		}
		time.Sleep(time.Second * 5)
		logger.Log.Debug(fmt.Sprintf("Retrying after 5 seconds. Retry Count :%v", retryCount+1))

	}
	return result
}

func (bmhnode *BMHNode) SetOneTimeBoot() bool {
	var result bool
	//redfishClient := getRedfishClient(bmhnode)
	if bmhnode.RedfishVersion == "1.0.0" {
		result = bmhnode.SetOneTimeBootILO4()

	} else {
		//result = redfishClient.SetOneTimeBoot(bmhnode.RedfishSystemID)
		result = bmhnode.SetOneTimeBootIDRAC()
	}
	return result

}

func (bmhnode *BMHNode) Reboot() bool {
	redfishClient := getRedfishClient(bmhnode)
	redfishSystemID := redfishClient.GetSystemID()
	result := redfishClient.RebootServer(redfishSystemID)
	return result
}

func (bmhnode *BMHNode) PowerOff() error {
	var err error
	redfishClient := getRedfishClient(bmhnode)
	redfishSystemID := redfishClient.GetSystemID()
	result := redfishClient.PowerOff(redfishSystemID)
	if result == false{
	    err =  errors.New("Failed to Power Off")
	}
	return err
}

func (bmhnode *BMHNode) PowerOn() error {
	var err error
	redfishClient := getRedfishClient(bmhnode)
	redfishSystemID := redfishClient.GetSystemID()
	result := redfishClient.PowerOn(redfishSystemID)
	if result == false{
	    err =  errors.New("Failed to Power On")
	}
	return err
}

func (bmhnode *BMHNode) GetPowerStatus() (bool,error) {
	redfishClient := getRedfishClient(bmhnode)
	redfishSystemID := redfishClient.GetSystemID()
	result := redfishClient.GetPowerStatus(redfishSystemID)
	return result,nil
}

func (bmhnode *BMHNode) EjectISO() bool {
	var result bool
	redfishClient := getRedfishClient(bmhnode)
	redfishManagerID := redfishClient.GetManagerID()
	if bmhnode.GetVirtualMediaStatus() == true {
		if bmhnode.RedfishVersion == "1.0.0" {
			result = bmhnode.EjectISOILO4(redfishManagerID)
		} else {
			result = redfishClient.EjectISO(redfishManagerID, "CD")
		}
	} else {
		logger.Log.Debug(fmt.Sprintf("Skipping Eject . VirtualMedia not attached \n"))
		result = true
	}
	return result

}



func GetUUID(hostIP string, username string, password string) (string, bool) {
	redfishClient := client.IdracRedfishClient{
		Username: username,
		Password: password,
		HostIP:   hostIP,
	}
	redfishSystemID := redfishClient.GetSystemID()
	if redfishSystemID == "" {
		return "", false
	}
	uuid, result := redfishClient.GetNodeUUID(redfishSystemID)
	return uuid, result
}

func (bmhnode *BMHNode) GetGUUID()([]byte, error){
	logger.Log.Info("Entering GetGUUID for node", zap.String("Node Name", bmhnode.Name), zap.String("Node UUID", bmhnode.NodeUUID.String()))
	var err error
	redfishClient := getRedfishClient(bmhnode)
	redfishSystemID := redfishClient.GetSystemID()
	if redfishSystemID == "" {
		logger.Log.Error("Failed to retrieve SystemID")
		return []byte(""),fmt.Errorf("Failed to retrieve SystemID")
	}
	uuid,res :=  redfishClient.GetNodeUUID(redfishSystemID)
	if res != true{
	    logger.Log.Error("Failed to retreive GUUID of Node")
	    err  = fmt.Errorf("Failed to retreive GUUID of Node")

	}
	return []byte(uuid), err
}



func (bmhnode *BMHNode) GetManagerID() string {
	redfishClient := getRedfishClient(bmhnode)
	return redfishClient.GetManagerID()
}
func (bmhnode *BMHNode) GetSystemID() string {
	redfishClient := getRedfishClient(bmhnode)
	return redfishClient.GetSystemID()
}

func (bmhnode *BMHNode) GetRedfishVersion() string {
	redfishClient := getRedfishClient(bmhnode)
	return redfishClient.GetRedfishVer()
}

func (bmhnode *BMHNode) DeployISO() error {
	logger.Log.Info("Entering DeplyISO for node", zap.String("Node Name", bmhnode.Name), zap.String("Node UUID", bmhnode.NodeUUID.String()))
	var errorString string
	var result bool
	// Setup Raid
        status , _  := bmhnode.GetPowerStatus() 
	if status == false {
		logger.Log.Warn("Node is in Powered Off state", zap.String("Node Name", bmhnode.Name))
		logger.Log.Info("Trying to Power On the node", zap.String("Node Name", bmhnode.Name))
		err := bmhnode.PowerOn()
		if err  != nil {
			errorString = fmt.Sprintf("Failed to power on Node %v", zap.String("Node Name", bmhnode.Name))
			logger.Log.Error(errorString)
			return fmt.Errorf("%v", errorString)
		}
		if PowerStateChangeTimeoutSeconds == 0 {
			PowerStateChangeTimeoutSeconds = 300 // if environment variable is not set
		}
		time.Sleep(time.Second * time.Duration(PowerStateChangeTimeoutSeconds))

	}

	if bmhnode.RAID_reset {
		err := bmhnode.ConfigureRAID()

		if err != nil {
			logger.Log.Error("Failed to create virtual disk", zap.String("Node Name", bmhnode.Name))
			return fmt.Errorf("Failed to create virtual Disk %v", err)
		}
	} else {
		logger.Log.Info("RAID Reset set to false. Skipping RAID Virtual Disk Creation", zap.String("Node Name", bmhnode.Name))
	}
	// Redfish steps or installation
	//Step 1 Eject CD
	//fmt.Sprintf("Step 1 Eject CD\n")
	logger.Log.Info("Step 1 Eject CD", zap.String("Node Name", bmhnode.Name))
	result = bmhnode.EjectISO()
	if result != false {
		//Step 2 Insert Ubuntu ISO
		logger.Log.Info("Step 2 Insert Ubuntu ISO", zap.String("Node Name", bmhnode.Name))
		//fmt.Sprintf("Step 2 Insert Ubuntu ISO\n")
		result = bmhnode.InsertISO()
		if result != false {
			//Step 3 Set Onetime boot to CD ROM
			logger.Log.Info("Step 3 Set Onetime boot to CD", zap.String("Node Name", bmhnode.Name))
			//fmt.Sprintf("Step 3 Set Onetime boot to CD ROM\n\n")
			result = bmhnode.SetOneTimeBoot()
			if result != false {
				//Step 4 Reboot
				logger.Log.Info("Step 4 Reboot", zap.String("Node Name", bmhnode.Name))
				//fmt.Sprintf("Step 4 Reboot\n")

				result = bmhnode.Reboot()

			}
		}
	}

	if result != true {
		return fmt.Errorf("Failed to deploy ISO")
	}
	return nil
}

func (bmhnode *BMHNode) SetRedfishIDs() {
	if bmhnode.RedfishVersion == "" {
		bmhnode.RedfishVersion = bmhnode.GetRedfishVersion()
	}
	if bmhnode.RedfishManagerID == "" {
		bmhnode.RedfishManagerID = bmhnode.GetManagerID()
	}
	if bmhnode.RedfishSystemID == "" {

		bmhnode.RedfishSystemID = bmhnode.GetSystemID()
	}
}

func (bmhnode  *BMHNode) GetHWInventory() (map[string]string,error){

	var err error
	var hwInfo   = make(map[string]string)

	bmhnode.SetRedfishIDs()

	hwInfo["RedfishVersion"] =  bmhnode.RedfishVersion
	hwInfo["RedfishManagerID"] = bmhnode.RedfishManagerID
	hwInfo["RedfishSystemID"] = bmhnode.RedfishSystemID

	for  k,v := range hwInfo{
		if v == ""{
			err =  fmt.Errorf("Failed to retrieve %v", k)
		}
	}

	return hwInfo, err
}
