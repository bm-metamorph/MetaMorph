package redfish

import (
	config "bitbucket.com/metamorph/pkg/config"
	"fmt"
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
	result := redfishClient.InsertISO(config.Get("idrac.managerID").(string), "CD", bmhnode.ImageURL)
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

func (bmhnode *BMHNode) GetUUID() (string, bool) {
	redfishClient := getRedfishClient(bmhnode)
	uuid, result := redfishClient.GetNodeUUID(config.Get("idrac.systemID").(string))
	return uuid, result

}

func (bmhnode *BMHNode) DeployISO() bool {
	fmt.Printf("Entering DeplyISO for node %v\n", bmhnode.NodeUUID)
	var result bool
	// Setup Raid

	result = bmhnode.CreateVirtualDisks()

	if result == false {
		fmt.Printf("Failed to create virtual disk\n")
		return result
	}
	// Redfish steps or installation
	//Step 1 Eject CD
	fmt.Printf("Step 1 Eject CD\n")
	result = bmhnode.EjectISO()
	if result != false {
		//Step 2 Insert Ubuntu ISO
		fmt.Printf("Step 2 Insert Ubuntu ISO\n")
		result = bmhnode.InsertISO()
		if result != false {
			//Step 3 Set Onetime boot to CD ROM
			fmt.Printf("Step 3 Set Onetime boot to CD ROM\n\n")
			result = bmhnode.SetOneTimeBoot()
			if result != false {
				//Step 4 Reboot
				fmt.Printf("Step 4 Reboot\n")

				result = bmhnode.Reboot()

			}
		}
	}

	return result

}
