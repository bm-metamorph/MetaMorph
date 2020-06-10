package redfish

import (
	"crypto/tls"
	"encoding/json"
	"fmt"

	resty "github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	//"net/http"
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
	resultBody := make(map[string]interface{})
	restyClient := resty.New()
	restyClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	restyClient.SetBasicAuth(bmhnode.IPMIUser, bmhnode.IPMIPassword)
	restyClient.SetDebug(true)

	resp, err := restyClient.R().EnableTrace().
		SetHeader("Content-Type", "application/json").
		SetBody([]byte(data)).Post(endpointURL)
	fmt.Println("Trace Info:", resp.Request.TraceInfo())

	if (err == nil) && checkStatusCode(resp.StatusCode()) {
		if len(resp.Body()) != 0 {
			err = json.Unmarshal(resp.Body(), &resultBody)
		}
	} else {
		return nil, errors.Wrap(err, fmt.Sprintf("Post request failed : URL - %v, reqbody - %v", endpointURL, string(data)))
	}
	return resultBody, err
}

func (bmhnode *BMHNode) PatchRequestToRedfish(endpointURL string, data []byte) (map[string]interface{}, error) {
	resultBody := make(map[string]interface{})
	restyClient := resty.New()
	restyClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	restyClient.SetBasicAuth(bmhnode.IPMIUser, bmhnode.IPMIPassword)
	restyClient.SetDebug(true)

	resp, err := restyClient.R().EnableTrace().
		SetHeader("Content-Type", "application/json").
		SetBody([]byte(data)).Patch(endpointURL)
	fmt.Println("Trace Info:", resp.Request.TraceInfo())

	if (err == nil) && checkStatusCode(resp.StatusCode()) {
		err = json.Unmarshal(resp.Body(), &resultBody)
	} else {
		return nil, errors.Wrap(err, fmt.Sprintf("Patch request failed : URL - %v, reqbody - %v", endpointURL, string(data)))
	}
	return resultBody, err
}

func (bmhnode *BMHNode) GetRequestRedfish(endpointURL string) (map[string]interface{}, error) {

	resultBody := make(map[string]interface{})
	restyClient := resty.New()
	restyClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	restyClient.SetBasicAuth(bmhnode.IPMIUser, bmhnode.IPMIPassword)
	restyClient.SetDebug(true)

	resp, err := restyClient.R().EnableTrace().Get(endpointURL)
	fmt.Println("Trace Info:", resp.Request.TraceInfo())

	if (err == nil) && checkStatusCode(resp.StatusCode()) {
		err = json.Unmarshal(resp.Body(), &resultBody)
	} else {
		return nil, errors.Wrap(err, fmt.Sprintf("Get request failed : URL - %v", endpointURL))
	}
	return resultBody, err

}

func (bmhnode *BMHNode) EjectISOILO4(redfishManagerID string) bool {
	var err error

	redfishEndpoint := fmt.Sprintf(managerEndPoint, bmhnode.IPMIIP, redfishManagerID)
	redfishEndpoint += ilo4VirtualMediaAction + ejectVirtualMedia

	requestBody := "{}"

	fmt.Printf("Endpoint %+v", redfishEndpoint)

	_, err = bmhnode.PostRequestToRedfish(redfishEndpoint, []byte(requestBody))
	if err != nil {
		return false
	}
	return true
}

func (bmhnode *BMHNode) InsertISOILO4(redfishManagerID string) bool {
	var err error
	redfishEndpoint := fmt.Sprintf(managerEndPoint, bmhnode.IPMIIP, redfishManagerID)
	redfishEndpoint += ilo4VirtualMediaAction + insertVirtualMedia

	requestBody := fmt.Sprintf("{ \"Image\" : \"%v\" }", bmhnode.ImageURL)

	fmt.Printf("Endpoint %+v", redfishEndpoint)

	_, err = bmhnode.PostRequestToRedfish(redfishEndpoint, []byte(requestBody))
	if err != nil {
		return false
	}
	return true
}

func (bmhnode *BMHNode) SetOneTimeBootILO4() bool {
	var err error
	redfishEndpoint := fmt.Sprintf(systemEndPoint, bmhnode.IPMIIP, bmhnode.RedfishSystemID)

	requestBody := "{ \"Boot\": { \"BootSourceOverrideEnabled\": \"Once\", \"BootSourceOverrideTarget\": \"Cd\" } }"

	fmt.Printf("Endpoint %+v\n", redfishEndpoint)

	fmt.Printf("Request Body %+v\n", requestBody)

	_, err = bmhnode.PatchRequestToRedfish(redfishEndpoint, []byte(requestBody))
	if err != nil {
		fmt.Printf("Failed to set one time boot\n")
		return false
	}
	return true
}
