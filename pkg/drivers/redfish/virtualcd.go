package redfish


import  (
	config "bitbucket.com/metamorph/pkg/config"
	"fmt"
)

func (bmhnode *BMHNode) GetVirtualMediaStatus() bool {
	redfishClient := getRedfishClient(bmhnode)
	result := redfishClient.GetVirtualMediaStatus(config.Get("idrac.managerID").(string),"CD")
	return result
}

func (bmhnode *BMHNode) InsertISO() bool {
	if bmhnode.GetVirtualMediaStatus() == true {
		fmt.Printf("Skipping Iso Insert. CD already Attached\n")
		return false
	}

	redfishClient := getRedfishClient(bmhnode)
	result  := redfishClient.InsertISO(config.Get("idrac.managerID").(string),"CD",bmhnode.ImageURL)
	return result
}

func (bmhnode * BMHNode) SetOneTimeBoot() bool {
	redfishClient := getRedfishClient(bmhnode)
	result  := redfishClient.SetOneTimeBoot(config.Get("idrac.systemID").(string))
	return result

}

func (bhmnode * BMHNode) Reboot() bool {
	redfishClient := getRedfishClient(bhmnode)
	result := redfishClient.RebootServer(config.Get("idrac.systemID").(string))
	return result
}

func (bmhnode *BMHNode) EjectISO()  bool {
	redfishClient := getRedfishClient(bmhnode)
	result := redfishClient.EjectISO(config.Get("idrac.managerID").(string),"CD")
	return result

}

func (bmhnode *BMHNode) GetUUID()(string, bool) {
	redfishClient := getRedfishClient(bmhnode)
	uuid, result := redfishClient.GetNodeUUID(config.Get("idrac.systemID").(string))
	return uuid,result

}