package redfishwrap

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	redfish "opendev.org/airship/go-redfish/client"

	//     "reflect"
	//	"github.com/antihax/optional"
	_nethttp "net/http"
	"os"
	"regexp"

	//     "io/ioutil"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var log = logf.Log.WithName("RedfishAPI")

type RedfishAPIWrapper interface {
	UpgradeFirmware(string) error
	CheckJobStatus(string)
	RebootServer(string) bool
	PowerOn(string) bool
	PowerOff(string) bool
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func checkStatusCodeforGet(statuscode int) bool {
	sucessCodes := []int{200, 204, 202}
	for _, x := range sucessCodes {
		if statuscode == x {
			return true
		}
	}
	return false
}

var tr *http.Transport = &http.Transport{
	MaxIdleConns:       10,
	IdleConnTimeout:    30 * time.Second,
	DisableCompression: true,
	TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
}

func createAPIClient(HeaderInfo map[string]string, hostIPV4addr string) *redfish.DefaultApiService {
	client := &http.Client{Transport: tr}
	cfg := &redfish.Configuration{
		BasePath:      "https://" + hostIPV4addr,
		DefaultHeader: make(map[string]string),
		UserAgent:     "go-redfish/client",
		HTTPClient:    client,
	}

	if len(HeaderInfo) != 0 {

		for key, value := range HeaderInfo {
			cfg.DefaultHeader[key] = value
		}
	}
	return redfish.NewAPIClient(cfg).DefaultApi
}

func GetTask(ctx context.Context, hostIPV4addr string, taskID string) (int, redfish.Task) {
	redfishApi := createAPIClient(make(map[string]string), hostIPV4addr)
	sl, response, err := redfishApi.GetTask(ctx, taskID)
	fmt.Printf("%+v %+v %+v", prettyPrint(sl), response, err)
	return response.StatusCode, sl
}

func GetTaskList(ctx context.Context, hostIPV4addr string) (int, int) {
	redfishApi := createAPIClient(make(map[string]string), hostIPV4addr)
	sl, response, err := redfishApi.GetTaskList(ctx)
	fmt.Printf("%+v %+v %+v", prettyPrint(sl), response, err)
	return response.StatusCode, len(sl.Members)
}

func GetVirtualMediaConnectedStatus(ctx context.Context, hostIPV4addr string, managerID string, media string) bool {
	redfishApi := createAPIClient(make(map[string]string), hostIPV4addr)
	//	sl, response, err := redfishApi.GetManagerVirtualMedia(ctx, "iDRAC.Embedded.1", "CD")
	sl, response, err := redfishApi.GetManagerVirtualMedia(ctx, managerID, media)
	fmt.Printf("%+v %+v %+v", prettyPrint(sl), response, err)
	if err != nil || sl.ConnectedVia == "NotConnected" {
		return false
	}
	return true
}

func UpdateService(ctx context.Context, hostIPV4addr string) string {
	redfishApi := createAPIClient(make(map[string]string), hostIPV4addr)
	// call the UpdateService and get the HttpPushURi
	sl, response, err := redfishApi.UpdateService(ctx)
	fmt.Printf("%+v %+v %+v", prettyPrint(sl), response, err)
	return sl.HttpPushUri
}

func HTTPUriDownload(ctx context.Context, hostIPV4addr string, filePath string, etag string) (string, error) {
	filehandle, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}
	defer filehandle.Close()
	reqBody := redfish.InlineObject{SoftwareImage: filehandle}
	//FirmwareInventoryDownloadImageOpts{SoftwareImage: optional.NewInterface(filehandle)}
	headerInfo := make(map[string]string)
	headerInfo["if-match"] = etag
	redfishApi := createAPIClient(headerInfo, hostIPV4addr)

	sl, response, err := redfishApi.FirmwareInventoryDownloadImage(ctx, reqBody)
	fmt.Printf("%+v %+v %+v", prettyPrint(sl), response, err)
	location, _ := response.Location()
	return string(location.RequestURI()), err

}

func GetETagHttpURI(ctx context.Context, hostIPV4addr string) string {
	redfishApi := createAPIClient(make(map[string]string), hostIPV4addr)
	sl, response, err := redfishApi.FirmwareInventory(ctx)
	fmt.Printf("%+v %+v %+v", prettyPrint(sl), response, err)
	etag := response.Header["Etag"]
	fmt.Printf("%v", etag[0])
	return etag[0]
}

func getJobID(response *_nethttp.Response) string {
	jobID_location := response.Header["Location"]
	re := regexp.MustCompile(`(JID_.*)`)
	jobID := re.FindStringSubmatch(jobID_location[0])[1]
	return jobID
}

func SimpleUpdateRequest(ctx context.Context, hostIPV4addr string, imageURI string) string {
	headerInfo := make(map[string]string)
	redfishApi := createAPIClient(headerInfo, hostIPV4addr)
	reqBody := new(redfish.SimpleUpdateRequestBody)
	localUriImage := imageURI
	reqBody.ImageURI = localUriImage
	sl, response, err := redfishApi.UpdateServiceSimpleUpdate(ctx, *reqBody)
	fmt.Printf("%+v %+v %+v", prettyPrint(sl), response, err)
	return getJobID(response)
}

func ResetServer(ctx context.Context, hostIPV4addr string, systemId string, resetRequestBody redfish.ResetRequestBody) bool {

	headerInfo := make(map[string]string)
	redfishApi := createAPIClient(headerInfo, hostIPV4addr)

	//	resetRequestBody := redfish.ResetRequestBody{ResetType: redfish.RESETTYPE_FORCE_RESTART}

	sl, response, err := redfishApi.ResetSystem(ctx, systemId, resetRequestBody)
	fmt.Printf("%+v %+v %+v", prettyPrint(sl), response, err)
	if err != nil || (checkStatusCodeforGet(response.StatusCode) != true) {
		return false
	}
	return true
}

func SetSystem(ctx context.Context, hostIPV4addr string, systemId string, computerSystem redfish.ComputerSystem) bool {
	headerInfo := make(map[string]string)
	redfishApi := createAPIClient(headerInfo, hostIPV4addr)

	sl, response, err := redfishApi.SetSystem(ctx, systemId, computerSystem)
	fmt.Printf("%+v %+v %+v", prettyPrint(sl), response, err)
	if err != nil || (checkStatusCodeforGet(response.StatusCode) != true) {
		return false
	}
	return true
}

func GetSystem(ctx context.Context, hostIPV4addr string, systemID string) (*redfish.ComputerSystem, bool) {
	headerInfo := make(map[string]string)
	redfishApi := createAPIClient(headerInfo, hostIPV4addr)

	sl, response, err := redfishApi.GetSystem(ctx, systemID)
	fmt.Printf("%+v %+v %+v", prettyPrint(sl), response, err)

	if err != nil || (checkStatusCodeforGet(response.StatusCode) != true) {
		return nil, false
	}

	return &sl, true

}

func EjectVirtualMedia(ctx context.Context, hostIPV4addr string, managerID string, media string) bool {
	headerInfo := make(map[string]string)
	redfishApi := createAPIClient(headerInfo, hostIPV4addr)

	body := make(map[string]interface{})

	sl, response, err := redfishApi.EjectVirtualMedia(ctx, managerID, media, body)
	fmt.Printf("%+v %+v %+v", prettyPrint(sl), response, err)
	if err != nil || (checkStatusCodeforGet(response.StatusCode) != true) {
		return false
	}

	return true

}

func InsertVirtualMedia(ctx context.Context, hostIPV4addr string, managerID string, mediaID string, insertMediaReqBody redfish.InsertMediaRequestBody) bool {
	headerInfo := make(map[string]string)
	redfishApi := createAPIClient(headerInfo, hostIPV4addr)

	sl, response, err := redfishApi.InsertVirtualMedia(ctx, managerID, mediaID, insertMediaReqBody)

	fmt.Printf("%+v %+v %+v", prettyPrint(sl), response, err)
	if err != nil || (checkStatusCodeforGet(response.StatusCode) != true) {
		return false
	}

	return true

}

func GetVolumes(ctx context.Context, hostIPV4addr string, systemID string, controllerID string) []redfish.IdRef {
	headerInfo := make(map[string]string)
	redfishApi := createAPIClient(headerInfo, hostIPV4addr)

	sl, response, err := redfishApi.GetVolumes(ctx, systemID, controllerID)

	fmt.Printf("%+v %+v %+v", prettyPrint(sl), response, err)
	if err != nil || (checkStatusCodeforGet(response.StatusCode) != true) {
		return nil
	}
	return sl.Members
}

func DeleteVirtualDisk(ctx context.Context, hostIPV4addr string, systemID string, storageID string) string {
	headerInfo := make(map[string]string)
	redfishApi := createAPIClient(headerInfo, hostIPV4addr)

	response, err := redfishApi.DeleteVirtualdisk(ctx, systemID, storageID)

	if err != nil || (checkStatusCodeforGet(response.StatusCode) != true) {
		return ""
	}

	fmt.Printf("\n%v\n", response.Request)

	fmt.Printf("\n%+v\n %+v\n", response, err)
	var jobid string = ""
	jobid = getJobID(response)

	return jobid

}

func CreateVirtualDisk(ctx context.Context, hostIPV4addr string, systemID string, controllerID string, createVirtualDiskRequestBody redfish.CreateVirtualDiskRequestBody) string {
	headerInfo := make(map[string]string)
	redfishApi := createAPIClient(headerInfo, hostIPV4addr)
	sl, response, err := redfishApi.CreateVirtualDisk(ctx, systemID, controllerID, createVirtualDiskRequestBody)
	if err != nil || (checkStatusCodeforGet(response.StatusCode) != true) {
		return ""
	}
	fmt.Printf("\n%v\n", response.Request)
	fmt.Printf("\n%+v\n %+v\n %+v\n", prettyPrint(sl), response, err)
	var jobid string = ""
	jobid = getJobID(response)
	return jobid
}

func ListManagers(ctx context.Context, hostIPV4addr string) []string {
	headerInfo := make(map[string]string)
	redfishApi := createAPIClient(headerInfo, hostIPV4addr)
	sl, response, err := redfishApi.ListManagers(ctx)
	if err != nil || (checkStatusCodeforGet(response.StatusCode) != true) {
		fmt.Printf("%+v", err)
		return nil
	}
	fmt.Printf("\n%v\n", response.Request)
	fmt.Printf("\n%+v\n %+v\n %+v\n", prettyPrint(sl), response, err)

	idrefs := sl.Members

	if idrefs == nil {
		fmt.Printf("Failed to retrieve Manager ID")
		return nil
	}
	return retrieveStringsFromIdrefList(idrefs)

}

func retrieveStringsFromIdrefList(idrefs []redfish.IdRef) []string {
	idList := []string{}
	for _, id := range idrefs {
		fmt.Printf("Idref ID %v\n", id.OdataId)
		idInfo := strings.Split(id.OdataId, "/")
		if idInfo != nil {
			lenOflist := len(idInfo)
			managerId := idInfo[lenOflist-1]
			if (managerId == "") && (lenOflist >= 3) {
				managerId = idInfo[lenOflist-2]
			}
			idList = append(idList, managerId)
		}
	}
	return idList

}

func ListSystems(ctx context.Context, hostIPV4addr string) []string {
	headerInfo := make(map[string]string)
	redfishApi := createAPIClient(headerInfo, hostIPV4addr)
	sl, response, err := redfishApi.ListSystems(ctx)
	if err != nil || (checkStatusCodeforGet(response.StatusCode) != true) {
		fmt.Printf("%+v", err)
		return nil
	}
	fmt.Printf("\n%v\n", response.Request)
	fmt.Printf("\n%+v\n %+v\n %+v\n", prettyPrint(sl), response, err)

	idrefs := sl.Members

	if idrefs == nil {
		fmt.Printf("Failed to retrieve System ID")
		return nil
	}
	return retrieveStringsFromIdrefList(idrefs)

}

func GetRoot(ctx context.Context, hostIPV4addr string) *redfish.Root {
	headerInfo := make(map[string]string)
	redfishApi := createAPIClient(headerInfo, hostIPV4addr)
	sl, response, err := redfishApi.GetRoot(ctx)
	if err != nil || (checkStatusCodeforGet(response.StatusCode) != true) {
		fmt.Printf("%+v", err)
		return nil
	}
	fmt.Printf("\n%v\n", response.Request)
	fmt.Printf("\n%+v\n %+v\n %+v\n", prettyPrint(sl), response, err)
	return &sl

}
