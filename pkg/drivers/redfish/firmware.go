package redfish

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"bitbucket.com/metamorph/pkg/db/models/node"
	"bitbucket.com/metamorph/pkg/util/isogen"
)

func (bmhnode *BMHNode) UpgradeFirmware(filepath string) bool {
	redfishClient := getRedfishClient(bmhnode)
	return redfishClient.UpgradeFirmware(filepath)
}

func (bmhnode *BMHNode) UpgradeFirmwareList() bool {
	//iterate through the list of firmwares
	var res bool = false
	firwareURLs, err := node.GetFirmwareURLs(bmhnode.NodeUUID.String())
	if err != nil {
		fmt.Printf("Failed to retreive firmwareURL")
		return false
	}
	for _, firmwareURL := range firwareURLs {

		filename := path.Base(firmwareURL.FirmwareURL)

		tempdir, err := ioutil.TempDir("/tmp", "firmware")
		if err != nil {
			fmt.Printf("Failed to create temporary directory")
			return false
		}
		defer os.RemoveAll(tempdir)
		firmwarefilepath := path.Join(tempdir, filename)
		err = isogen.DownloadUrl(firmwarefilepath, firmwareURL.FirmwareURL)

		if err != nil {
			fmt.Printf("Failed to Download URL")
			return false
		}
		res = bmhnode.UpgradeFirmware(firmwarefilepath)

	}
	return res

}
