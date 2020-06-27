package redfish

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"bitbucket.com/metamorph/pkg/db/models/node"
	"bitbucket.com/metamorph/pkg/util/isogen"
	version "github.com/hashicorp/go-version"
)

func (bmhnode *BMHNode) UpgradeFirmware(filepath string) bool {
	redfishClient := getRedfishClient(bmhnode)
	return redfishClient.UpgradeFirmware(filepath)
}

func IsVersionHigher(providedVersion string, versionfromNode string) bool {
	var err error
	vprovided, err := version.NewVersion(providedVersion)
	if err != nil {
		fmt.Printf("Provided version could not be parsed")
		return false
	}
	vfromNode, err := version.NewVersion(versionfromNode)
	if err != nil {
		fmt.Printf("Firmware version from node could not be parsed")
		return false
	}
	if vfromNode.Equal(vprovided) {
		fmt.Printf("Version provided is equal to one in the node. Proceeding with installation..")
		return  true
	}
	if vfromNode.LessThan(vprovided) {
		fmt.Printf("Version provided is lower than one in the node")
		return false
	}
	return true
}

func (bmhnode *BMHNode) CheckUpgradeAllowed(providedName string, providedVersion string) bool {
	redfishClient := getRedfishClient(bmhnode)
	name, version, updateavailable := redfishClient.GetFirmwareDetails(providedName)
	if (name == "") || (version == "") {
		fmt.Printf("Failed to retrieve firmware details for %v\n", providedName)
		return false
	}
	if updateavailable == false {
		fmt.Printf("Firmware not updateable")
		return false
	}
	if IsVersionHigher(providedVersion, version) {
		return true
	}
	return true
}

func (bmhnode *BMHNode) UpgradeFirmwareList() bool {
	//iterate through the list of firmwares
	var res bool = false
	firwares, err := node.GetFirmwares(bmhnode.NodeUUID.String())
	if err != nil {
		fmt.Printf("Failed to retreive firmwareURL")
		return false
	}
	for _, firmware := range firwares {

		if bmhnode.CheckUpgradeAllowed(firmware.Name, firmware.Version) {

			filename := path.Base(firmware.URL)

			tempdir, err := ioutil.TempDir("/tmp", "firmware")
			if err != nil {
				fmt.Printf("Failed to create temporary directory")
				return false
			}
			defer os.RemoveAll(tempdir)
			firmwarefilepath := path.Join(tempdir, filename)
			err = isogen.DownloadUrl(firmwarefilepath, firmware.URL)

			if err != nil {
				fmt.Printf("Failed to Download URL")
				return false
			}
			res = bmhnode.UpgradeFirmware(firmwarefilepath)
			if res == false {
				return false
			}
		}
	}

	return res

}
