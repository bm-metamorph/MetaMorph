package redfish

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/bm-metamorph/MetaMorph/pkg/db/models/node"
	"github.com/manojkva/metamorph-plugin/pkg/util"
	"github.com/manojkva/metamorph-plugin/pkg/logger"
	version "github.com/hashicorp/go-version"
	"go.uber.org/zap"
)

func (bmhnode *BMHNode) UpgradeEachFirmware(filepath string) bool {
	logger.Log.Info("UpgradeEachFirmware()")
	redfishClient := getRedfishClient(bmhnode)
	return redfishClient.UpgradeFirmware(filepath)
}

func IsVersionHigher(providedVersion string, versionfromNode string) bool {
	logger.Log.Info("IsVersionHigher()")
	var err error
	vprovided, err := version.NewVersion(providedVersion)
	if err != nil {
		logger.Log.Error(fmt.Sprintf("Provided version could not be parsed"))
		return false
	}
	vfromNode, err := version.NewVersion(versionfromNode)
	if err != nil {
		logger.Log.Error(fmt.Sprintf("Firmware version from node could not be parsed"))
		return false
	}
	if vfromNode.Equal(vprovided) {
		logger.Log.Debug(fmt.Sprintf("Version provided is equal to one in the node. Proceeding with installation.."))
		return  true
	}
	if vfromNode.LessThan(vprovided) {
		logger.Log.Error(fmt.Sprintf("Version provided is lower than one in the node"))
		return false
	}
	return true
}

func (bmhnode *BMHNode) CheckUpgradeAllowed(providedName string, providedVersion string) bool {
	redfishClient := getRedfishClient(bmhnode)
	name, version, updateavailable := redfishClient.GetFirmwareDetails(providedName)
	if (name == "") || (version == "") {
		logger.Log.Error(fmt.Sprintf("Failed to retrieve firmware details for %v\n", providedName))
		return false
	}
	if updateavailable == false {
		logger.Log.Error(fmt.Sprintf("Firmware not updateable"),zap.String("Firmware", providedName), zap.String("Provided Version", providedVersion))
		return false
	}
	if IsVersionHigher(providedVersion, version) {
		logger.Log.Debug("Provided Version is Higher than the present one", zap.String("Provided Version", providedVersion),zap.String("Current Version", version))
		return true
	}
	return true
}

func (bmhnode *BMHNode) UpdateFirmware() error {
	logger.Log.Info("UpdateFirmware()")
	//iterate through the list of firmwares
	firwares, err := node.GetFirmwares(bmhnode.NodeUUID.String())
	if err != nil {
		logger.Log.Error("Failed to retreive firmwareURL", zap.Error(err))
		return fmt.Errorf("Failed to retreive firmwareURL. %v", err)
	}
	for _, firmware := range firwares {

		if bmhnode.CheckUpgradeAllowed(firmware.Name, firmware.Version) {

			filename := path.Base(firmware.URL)

			tempdir, err := ioutil.TempDir("/tmp", "firmware")
			if err != nil {
				logger.Log.Error("Failed to create temporary directory.", zap.Error(err))
				return fmt.Errorf("Failed to create temporary directory. %v", err)
			}
			defer os.RemoveAll(tempdir)
			firmwarefilepath := path.Join(tempdir, filename)
			err = util.DownloadUrl(firmwarefilepath, firmware.URL)

			if err != nil {
				logger.Log.Error("Failed to Download URL", zap.Error(err),zap.String("FirmwareURL", firmware.URL))
				return fmt.Errorf("Failed to Download URL, %v",err)
			}
			res := bmhnode.UpgradeEachFirmware(firmwarefilepath)
			if res == false {
				logger.Log.Error("Failed to Upgrede Firmware")
				return  fmt.Errorf("Failed to Upgrade Firmware")
			}
		} else{
			logger.Log.Error("Upgrade Validation of firmware failed", zap.String("FirmwareName",firmware.Name), zap.String("Version",firmware.Version))
			return fmt.Errorf("Check for Upgrade version info failed for %v and version %v", firmware.Name, firmware.Version)
		}

	}

	return nil

}
