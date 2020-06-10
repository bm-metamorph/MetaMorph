package redfish

import (
	//config "bitbucket.com/metamorph/pkg/config"
	"bitbucket.com/metamorph/pkg/logger"
	"fmt"
	client "github.com/manojkva/go-redfish-api-wrapper/pkg/redfishwrap/idrac"
	"go.uber.org/zap"
	"os"
	"strconv"
	"time"
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
		fmt.Printf("Skipping Iso Insert. CD already Attached\n")
		return false
	}

	redfishClient := getRedfishClient(bmhnode)
	fmt.Printf("Image URL to be inserted = %v", bmhnode.ImageURL)
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
		fmt.Printf("Retrying after 5 seconds. Retry Count :%v", retryCount+1)

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
	result := redfishClient.RebootServer(bmhnode.RedfishSystemID)
	return result
}

func (bmhnode *BMHNode) PowerOff() bool {
	redfishClient := getRedfishClient(bmhnode)
	result := redfishClient.PowerOff(bmhnode.RedfishSystemID)
	return result
}

func (bmhnode *BMHNode) PowerOn() bool {
	redfishClient := getRedfishClient(bmhnode)
	result := redfishClient.PowerOn(bmhnode.RedfishSystemID)
	return result
}

func (bmhnode *BMHNode) GetPowerStatus() bool {
	redfishClient := getRedfishClient(bmhnode)
	result := redfishClient.GetPowerStatus(bmhnode.RedfishSystemID)
	return result
}

func (bmhnode *BMHNode) EjectISO() bool {
	var result bool
	redfishClient := getRedfishClient(bmhnode)
	if bmhnode.GetVirtualMediaStatus() == true {
		if bmhnode.RedfishVersion == "1.0.0" {
			result = bmhnode.EjectISOILO4(bmhnode.RedfishManagerID)
		} else {
			result = redfishClient.EjectISO(bmhnode.RedfishManagerID, "CD")
		}
	} else {
		fmt.Printf("Skipping Eject . VirtualMedia not attached \n")
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

func (bmhnode *BMHNode) DeployISO() bool {
	logger.Log.Info("Entering DeplyISO for node", zap.String("Node Name", bmhnode.Name), zap.String("Node UUID", bmhnode.NodeUUID.String()))
	var result bool
	// Setup Raid

	if bmhnode.GetPowerStatus() == false {
		logger.Log.Warn("Node is in Powered Off state", zap.String("Node Name", bmhnode.Name))
		logger.Log.Info("Trying to Power On the node", zap.String("Node Name", bmhnode.Name))
		result = bmhnode.PowerOn()
		if result == false {
			logger.Log.Error("Failed to power on Node", zap.String("Node Name", bmhnode.Name))
			return result
		}
		if PowerStateChangeTimeoutSeconds == 0 {
			PowerStateChangeTimeoutSeconds = 300 // if environment variable is not set
		}
		time.Sleep(time.Second * time.Duration(PowerStateChangeTimeoutSeconds))

	}

	if bmhnode.RAID_reset {
		result = bmhnode.CreateVirtualDisks()

		if result == false {
			logger.Log.Error("Failed to create virtual disk", zap.String("Node Name", bmhnode.Name))
			return result
		}
	} else {
		logger.Log.Info("RAID Reset set to false. Skipping RAID Virtual Disk Creation", zap.String("Node Name", bmhnode.Name))
	}
	// Redfish steps or installation
	//Step 1 Eject CD
	fmt.Printf("Step 1 Eject CD\n")
	logger.Log.Info("Step 1 Eject CD", zap.String("Node Name", bmhnode.Name))
	result = bmhnode.EjectISO()
	if result != false {
		//Step 2 Insert Ubuntu ISO
		logger.Log.Info("Step 2 Insert Ubuntu ISO", zap.String("Node Name", bmhnode.Name))
		fmt.Printf("Step 2 Insert Ubuntu ISO\n")
		result = bmhnode.InsertISO()
		if result != false {
			//Step 3 Set Onetime boot to CD ROM
			logger.Log.Info("Step 3 Set Onetime boot to CD", zap.String("Node Name", bmhnode.Name))
			fmt.Printf("Step 3 Set Onetime boot to CD ROM\n\n")
			result = bmhnode.SetOneTimeBoot()
			if result != false {
				//Step 4 Reboot
				logger.Log.Info("Step 4 Reboot", zap.String("Node Name", bmhnode.Name))
				fmt.Printf("Step 4 Reboot\n")

				result = bmhnode.Reboot()

			}
		}
	}

	return result

}
