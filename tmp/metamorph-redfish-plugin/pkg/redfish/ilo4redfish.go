package redfish

import (
	"crypto/tls"
	"encoding/json"
	"fmt"

	resty "github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/manojkva/metamorph-plugin/pkg/logger"
	"go.uber.org/zap"
)

const ilo4VirtualMediaAction string = "/VirtualMedia/2/Actions/Oem/Hp/HpiLOVirtualMedia."
const ejectVirtualMedia string = "EjectVirtualMedia/"
const insertVirtualMedia string = "InsertVirtualMedia/"
const managerEndPoint string = "https://%v/redfish/v1/Managers/%v"
const systemEndPoint string = "https://%v/redfish/v1/Systems/%v/"

func checkStatusCode(statuscode int) bool {
	sucessCodes := []int{200, 204, 202}
	for _, x := range sucessCodes {
		if statuscode == x {
			return true
		}
	}
	return false
}

func (bmhnode *BMHNode) PostRequestToRedfish(endpointURL string, data []byte) (map[string]interface{}, error) {
	logger.Log.Info("PostRequestToRedfish")
	resultBody := make(map[string]interface{})
	restyClient := resty.New()
	restyClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	restyClient.SetBasicAuth(bmhnode.IPMIUser, bmhnode.IPMIPassword)
	restyClient.SetDebug(true)

	resp, err := restyClient.R().EnableTrace().
		SetHeader("Content-Type", "application/json").
		SetBody([]byte(data)).Post(endpointURL)
	logger.Log.Debug(fmt.Sprintf("Trace Info:", resp.Request.TraceInfo()))

	if (err == nil) && checkStatusCode(resp.StatusCode()) {
		if len(resp.Body()) != 0 {
			err = json.Unmarshal(resp.Body(), &resultBody)
		}
	} else {
		logger.Log.Error("Post request failed" , zap.String("endpointURL", endpointURL), zap.String("RequestBody", string(data)))
		return nil, errors.Wrap(err, fmt.Sprintf("Post request failed : URL - %v, reqbody - %v", endpointURL, string(data)))
	}
	return resultBody, err
}

func (bmhnode *BMHNode) PatchRequestToRedfish(endpointURL string, data []byte) (map[string]interface{}, error) {
	logger.Log.Info("PatchRequestToRedfish()")
	resultBody := make(map[string]interface{})
	restyClient := resty.New()
	restyClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	restyClient.SetBasicAuth(bmhnode.IPMIUser, bmhnode.IPMIPassword)
	restyClient.SetDebug(true)

	resp, err := restyClient.R().EnableTrace().
		SetHeader("Content-Type", "application/json").
		SetBody([]byte(data)).Patch(endpointURL)
	logger.Log.Debug(fmt.Sprintf("Trace Info:", resp.Request.TraceInfo()))

	if (err == nil) && checkStatusCode(resp.StatusCode()) {
		err = json.Unmarshal(resp.Body(), &resultBody)
	} else {
		logger.Log.Error("Patch request failed", zap.String("Endpoint URL", endpointURL), zap.String("Request Body", string(data)))
		return nil, errors.Wrap(err, fmt.Sprintf("Patch request failed : URL - %v, reqbody - %v", endpointURL, string(data)))
	}
	return resultBody, err
}

func (bmhnode *BMHNode) GetRequestRedfish(endpointURL string) (map[string]interface{}, error) {

	logger.Log.Info("GetRequestRedfish()")

	resultBody := make(map[string]interface{})
	restyClient := resty.New()
	restyClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	restyClient.SetBasicAuth(bmhnode.IPMIUser, bmhnode.IPMIPassword)
	restyClient.SetDebug(true)

	resp, err := restyClient.R().EnableTrace().Get(endpointURL)
	logger.Log.Debug(fmt.Sprintf("Trace Info:", resp.Request.TraceInfo()))

	if (err == nil) && checkStatusCode(resp.StatusCode()) {
		err = json.Unmarshal(resp.Body(), &resultBody)
	} else {
		logger.Log.Error("Get request failed", zap.String("URL", endpointURL))
		return nil, errors.Wrap(err, fmt.Sprintf("Get request failed : URL - %v", endpointURL))
	}
	return resultBody, err

}

func (bmhnode *BMHNode) EjectISOILO4(redfishManagerID string) bool {
	logger.Log.Info("EjectISOILO4()")
	var err error

	redfishEndpoint := fmt.Sprintf(managerEndPoint, bmhnode.IPMIIP, redfishManagerID)
	redfishEndpoint += ilo4VirtualMediaAction + ejectVirtualMedia

	requestBody := "{}"

	fmt.Sprintf("Endpoint %+v", redfishEndpoint)

	_, err = bmhnode.PostRequestToRedfish(redfishEndpoint, []byte(requestBody))
	if err != nil {
		logger.Log.Error("Failed to EjectISO",zap.Error(err))
		return false
	}
	return true
}

func (bmhnode *BMHNode) InsertISOILO4(redfishManagerID string) bool {
	logger.Log.Info("InsertISOILO4()")
	var err error
	redfishEndpoint := fmt.Sprintf(managerEndPoint, bmhnode.IPMIIP, redfishManagerID)
	redfishEndpoint += ilo4VirtualMediaAction + insertVirtualMedia

	requestBody := fmt.Sprintf("{ \"Image\" : \"%v\" }", bmhnode.ImageURL)

	fmt.Sprintf("Endpoint %+v", redfishEndpoint)

	_, err = bmhnode.PostRequestToRedfish(redfishEndpoint, []byte(requestBody))
	if err != nil {
		logger.Log.Error("Failed to InsertISO",zap.Error(err))
		return false
	}
	return true
}

func (bmhnode *BMHNode) SetOneTimeBootILO4() bool {
	logger.Log.Info("SetOneTimeBootILO4")
	var err error
	redfishEndpoint := fmt.Sprintf(systemEndPoint, bmhnode.IPMIIP, bmhnode.RedfishSystemID)

	requestBody := "{ \"Boot\": { \"BootSourceOverrideEnabled\": \"Once\", \"BootSourceOverrideTarget\": \"Cd\" } }"

	fmt.Sprintf("Endpoint %+v\n", redfishEndpoint)

	fmt.Sprintf("Request Body %+v\n", requestBody)

	_, err = bmhnode.PatchRequestToRedfish(redfishEndpoint, []byte(requestBody))
	if err != nil {
		logger.Log.Error(fmt.Sprintf("Failed to set one time boot\n"),zap.Error(err))
		return false
	}
	return true
}
