package redfish

import (
	config "bitbucket.com/metamorph/pkg/config"
	"fmt"
	client "github.com/manojkva/go-redfish-api-wrapper/pkg/redfishwrap/idrac"
	"time"
        "go.uber.org/zap"
        "bitbucket.com/metamorph/pkg/logger"
)

func (bmhnode *BMHNode) GetVirtualMediaStatus() bool {
	redfishClient := getRedfishClient(bmhnode)
	result := redfishClient.GetVirtualMediaStatus(config.Get("idrac.managerID").(string), "CD")
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
		result = redfishClient.InsertISO(config.Get("idrac.managerID").(string), "CD", bmhnode.ImageURL)
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
	redfishClient := getRedfishClient(bmhnode)
	result := redfishClient.SetOneTimeBoot(config.Get("idrac.systemID").(string))
	return result

}

func (bhmnode *BMHNode) Reboot() bool {
	redfishClient := getRedfishClient(bhmnode)
	result := redfishClient.RebootServer(config.Get("idrac.systemID").(string))
	return result
}

func (bhmnode *BMHNode) PowerOff() bool {
	redfishClient := getRedfishClient(bhmnode)
	result := redfishClient.PowerOff(config.Get("idrac.systemID").(string))
	return result
}

func (bhmnode *BMHNode) PowerOn() bool {
	redfishClient := getRedfishClient(bhmnode)
	result := redfishClient.PowerOn(config.Get("idrac.systemID").(string))
	return result
}

func (bhmnode *BMHNode) GetPowerStatus() bool {
	redfishClient := getRedfishClient(bhmnode)
	result := redfishClient.GetPowerStatus(config.Get("idrac.systemID").(string))
	return result
}

func (bmhnode *BMHNode) EjectISO() bool {
	var result bool
	redfishClient := getRedfishClient(bmhnode)
	if bmhnode.GetVirtualMediaStatus() == true {
		result = redfishClient.EjectISO(config.Get("idrac.managerID").(string), "CD")
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
	uuid, result := redfishClient.GetNodeUUID(config.Get("idrac.systemID").(string))
	return uuid, result

}

func (bmhnode *BMHNode) DeployISO() bool {
        logger.Log.Info("Entering DeplyISO for node", zap.String("Node Name", bmhnode.Name), zap.String("Node UUID", bmhnode.NodeUUID.String()))
	var result bool
	// Setup Raid

        if  bmhnode.GetPowerStatus() == false {
                logger.Log.Warn("Node is in Powered Off state", zap.String("Node Name",  bmhnode.Name))
                logger.Log.Info("Trying to Power On the node", zap.String("Node Name",  bmhnode.Name))
                result = bmhnode.PowerOn()
                if result == false {
                  logger.Log.Error("Failed to power on Node", zap.String("Node Name", bmhnode.Name))
                  return result
          }
        }

        if  bmhnode.RAID_reset {
              result = bmhnode.CreateVirtualDisks()

	      if result == false {
                logger.Log.Error("Failed to create virtual disk",zap.String("Node Name", bmhnode.Name))
		return result
	       }
         } else{
              logger.Log.Info("RAID Reset set to false. Skipping RAID Virtual Disk Creation", zap.String("Node Name",  bmhnode.Name))
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
