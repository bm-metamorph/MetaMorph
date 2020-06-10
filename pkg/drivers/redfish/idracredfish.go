package redfish

import (
	"fmt"
)

const setOnceRebootOEMIDRACendpoint = "/Actions/Oem/EID_674_Manager.ImportSystemConfiguration"

func (bmhnode *BMHNode) SetOneTimeBootIDRAC() bool {
	var err error
	redfishEndpoint := fmt.Sprintf(managerEndPoint, bmhnode.IPMIIP, bmhnode.RedfishManagerID)
	redfishEndpoint += setOnceRebootOEMIDRACendpoint

	requestBody := `{

             "ShareParameters": { 
                 "Target": "ALL" 
              },

            "ImportBuffer": 
                "<SystemConfiguration><Component FQDD=\"%v\">
                 <Attribute Name=\"ServerBoot.1#BootOnce\">Enabled</Attribute>
                  <Attribute Name=\"ServerBoot.1#FirstBootDevice\">VCD-DVD</Attribute></Component></SystemConfiguration>"
              }`
	requestBody = fmt.Sprintf(requestBody, bmhnode.RedfishManagerID)

	fmt.Printf("Endpoint %+v\n", redfishEndpoint)

	fmt.Printf("Request Body %+v\n", requestBody)

	_, err = bmhnode.PostRequestToRedfish(redfishEndpoint, []byte(requestBody))
	if err != nil {
		fmt.Printf("Failed to set one time boot for IDRAC \n")
		return false
	}
	return true
}
