package redfish

import (
	"fmt"
	"github.com/manojkva/metamorph-plugin/pkg/logger"
)

const setOnceRebootOEMIDRACendpoint = "/Actions/Oem/EID_674_Manager.ImportSystemConfiguration"

func (bmhnode *BMHNode) SetOneTimeBootIDRAC() bool {
	logger.Log.Info("SetOneTimeBootIDRAC()")
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

	logger.Log.Debug(fmt.Sprintf("Endpoint %+v\n", redfishEndpoint))

	logger.Log.Debug(fmt.Sprintf("Request Body %+v\n", requestBody))

	_, err = bmhnode.PostRequestToRedfish(redfishEndpoint, []byte(requestBody))
	if err != nil {
		logger.Log.Error(fmt.Sprintf("Failed to set one time boot for IDRAC \n"))
		return false
	}
	return true
}
