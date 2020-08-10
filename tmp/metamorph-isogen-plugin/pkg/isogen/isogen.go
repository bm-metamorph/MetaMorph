package isogen

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"crypto/md5"
	"io"
	"io/ioutil"

	"path"
	"strings"

	config "github.com/manojkva/metamorph-plugin/pkg/config"
	"github.com/manojkva/metamorph-plugin/pkg/util"
	"github.com/bm-metamorph/MetaMorph/pkg/db/models/node"
	"github.com/manojkva/metamorph-plugin/pkg/logger"
	"go.uber.org/zap"
)

func CreateDirectory(directoryPath string) error {
	logger.Log.Info("CreateDirectory()", zap.String("dirPath", directoryPath))
	pathErr := os.MkdirAll(directoryPath, 0777)
	if pathErr != nil {
		logger.Log.Error("Failed to create directory",
			zap.String("dirpath", directoryPath),
			zap.Error(pathErr))
		return pathErr
	}
	return nil
}


func ExtractIso(iso, target string) error {
	logger.Log.Info("ExtractingIso()", zap.String("isofile", iso), zap.String("TargetPath", target))
	cmd := exec.Command("mount", "-r", "-o", "loop", iso, target)
	err := cmd.Run()
	if err != nil {
		errMessage := "Failed to run mount command"
		logger.Log.Error(errMessage, zap.Error(err))
		return fmt.Errorf(errMessage+" %+v", err)
	}
	return err
}

func Checksum(file string) string {
	logger.Log.Info("Checksum()", zap.String("filename", file))
	f, err := os.Open(file)
	if err != nil {
		logger.Log.Error("Failed to openfile ", zap.String("filename", file), zap.Error(err))
	}
	defer f.Close()
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		logger.Log.Error("Failed to copy hash from MD5sum ", zap.Error(err))
	}
	checksum := fmt.Sprintf("%x", h.Sum(nil))
	return checksum
}

func CopyfileToDestination(sourcefilepath string, destinationfilepath string) error {

	logger.Log.Info("CopyfileToDestination()", zap.String("sourcefilepath", sourcefilepath),
		zap.String("destinationfilepath", destinationfilepath))

	input, err := ioutil.ReadFile(sourcefilepath)
	if err != nil {
		logger.Log.Error("Failed to read file with error", zap.String("sourcefilepath", sourcefilepath), zap.Error(err))
		return fmt.Errorf("Failed to read file %v with error : %+v", sourcefilepath, err)
	}
	err = ioutil.WriteFile(destinationfilepath, input, 0644)
	if err != nil {
		logger.Log.Error("Failed to write file to destination", zap.String("sourcefilepath", sourcefilepath),
			zap.String("destinationpath", destinationfilepath),
			zap.Error(err))
		return fmt.Errorf("Failed to write %v to destination %+v with error %+v", sourcefilepath, destinationfilepath, err)
	}
	return nil
}

func (bmhnode *BMHNode) CreateISO() error {
	logger.Log.Info("CreateISO()..")
	//check if ISO generation if required at all

	var err error
	var errMessage string
	if bmhnode.ImageReadilyAvailable {
		errMessage = "Customised ISO already available. No need to prepare custom ISO"
		logger.Log.Error(errMessage)
		return fmt.Errorf(errMessage)
	}

	iso_rootpath := config.Get("iso.rootpath").(string)
	HTTPRootPath := config.Get("http.rootpath").(string)

	iso_urlpath := bmhnode.ISOURL
	iso_checksum := bmhnode.ISOChecksum

	iso_name_parts := strings.Split(iso_urlpath, "/")
	iso_name := iso_name_parts[len(iso_name_parts)-1]

	if _, err := os.Stat(iso_rootpath); os.IsNotExist(err) {
                errMessage = "iso root directory not found : "
                logger.Log.Error(errMessage, zap.String("ISO Root",iso_rootpath),zap.Error(err))
		return fmt.Errorf(errMessage + " %v\n", err)
	}

	if _, err := os.Stat(HTTPRootPath); os.IsNotExist(err) {
                errMessage = "HTTP root directory not found :"
                logger.Log.Error(errMessage,zap.String("HTTP Root",HTTPRootPath), zap.Error(err))
		return fmt.Errorf(errMessage +" %v\n", err)
	}
	//Get temporary directory for copying ISO

	iso_tempdir := config.Get("iso.tempdir").(string)

	iso_DownloadFullpath := path.Join(iso_tempdir, iso_name)
	iso_DestinationFullpath := path.Join(iso_rootpath, bmhnode.NodeUUID.String(), iso_name)

	if _, err := os.Stat(iso_DownloadFullpath); os.IsNotExist(err) {
		err = os.MkdirAll(iso_tempdir, os.ModePerm)
		if err != nil {
                        errMessage = "Failed to create dir %v with error " 
                        logger.Log.Error(errMessage, zap.String("iso temp dir", iso_tempdir), zap.Error(err))
			return fmt.Errorf(errMessage + " %+v", iso_tempdir, err)
		}
		err = util.DownloadUrl(iso_DownloadFullpath, iso_urlpath)
		if err != nil {
                        errMessage = "Failed to download ISO Image :" 
                        logger.Log.Error(errMessage, zap.Error(err))
			return fmt.Errorf(errMessage +" %+v", err)
		}

	} else {
                logger.Log.Info("ISO vanilla Image already downloaded ")
	}

	err = ValidateChecksum(iso_checksum, iso_DownloadFullpath)

	if err != nil {
                errMessage = "Failed to validate checksum. Error :"
                logger.Log.Error(errMessage,zap.Error(err))
		return fmt.Errorf(errMessage +" %+v", err)
	}

	err = ExtractAndCopyISO(iso_DownloadFullpath, iso_DestinationFullpath)

	if err != nil {
                errMessage = "Failed to Extract and Copy ISO file. Error "
                logger.Log.Error(errMessage, zap.Error(err))
		return fmt.Errorf(errMessage +" %+v", err)
	}

	//Add Preseed
	err = bmhnode.CreatePressedFileFromTemplate(iso_DestinationFullpath, "preseed")

	if err != nil {
                errMessage = "Failed to create Preseed file with error" 
                logger.Log.Error(errMessage, zap.Error(err))
		return fmt.Errorf(errMessage + "%+v", err)
	}

	//Add Grub Config

	iso_custom_scripts_path := path.Join(iso_DestinationFullpath, "/setup/install/")

	err = os.MkdirAll(iso_custom_scripts_path, os.ModePerm)

	if err != nil {
		logger.Log.Error("Failed to create directory with error", zap.String("ISO custom scripts path",iso_custom_scripts_path), zap.Error(err))
		return fmt.Errorf("Failed to create directory %v with error %v", iso_custom_scripts_path, err)
	}

	err = bmhnode.CreateFileFromTemplate(iso_custom_scripts_path, "grub")

	if err != nil {
                errMessage = "Failed to create Grub file with error"
		logger.Log.Error(errMessage, zap.Error(err))
		return fmt.Errorf(errMessage + "%+v", err)
	}

	isolinux_cfg_destpath := iso_DestinationFullpath + "/isolinux/txt.cfg"

	metamorph_root := config.Get("templates.rootdir").(string)

	isolinuxtemplatepath := config.Get("templates.isolinux.config").(string)

	isolinux_cfg_sourcepath := path.Join(metamorph_root, isolinuxtemplatepath)

	err = CopyfileToDestination(isolinux_cfg_sourcepath, isolinux_cfg_destpath)

	if err != nil {
                errMessage = "Failed to copy isolinux cfg file with error"
		logger.Log.Error(errMessage, zap.Error(err))
		return fmt.Errorf(errMessage + "%+v", err)
	}
	//Netplan
	// Try using the NetPlanCloudInit string
	err = bmhnode.CreateNetplanFileFromString(iso_custom_scripts_path, "netplan")
	if err !=  nil {
		// Try using the standalone network config data structure
		err = bmhnode.CreateNetplanFileFromTemplate(iso_custom_scripts_path, "netplan")
	}

	if err != nil {
                errMessage = "Failed to create Netplan file with error"
		logger.Log.Error(errMessage , zap.Error(err))
		return fmt.Errorf(errMessage + "%+v", err)
	}

	//MetaMorph Agent
	metamorph_assets_root := config.Get("assets.rootdir").(string)
	metamorph_agent_file_src := config.Get("assets.agent_binary.src").(string)
	metamorph_agent_file_src_abs := path.Join(metamorph_assets_root, metamorph_agent_file_src)

	metamorph_agent_file_dest := config.Get("assets.agent_binary.dest").(string)
	metamorph_agent_file_dest_abs := path.Join(iso_custom_scripts_path, metamorph_agent_file_dest)
	err = CopyfileToDestination(metamorph_agent_file_src_abs, metamorph_agent_file_dest_abs)

	if err != nil {
                errMessage = "Failed to copy metamorph agent file with error" 
		logger.Log.Error(errMessage, zap.Error(err))
		return fmt.Errorf(errMessage + "%+v", err)
	}

	//MetaMorph Agent config
	bmhnode.ProvisioningIP = config.Get("provisioning.ip").(string)
	err = bmhnode.CreateFileFromTemplate(iso_custom_scripts_path, "agent_config")

	if err != nil {
                errMessage = "Failed to create Metamorph Agent config file with error"
		logger.Log.Error( errMessage, zap.Error(err))
		return fmt.Errorf(errMessage + " %+v", err)
	}

	//metamorph-client.service
	metamorph_servicetemplatepath := config.Get("templates.service.config").(string)

	metamorph_servicefilesourepath := path.Join(metamorph_root, metamorph_servicetemplatepath)

	metamorph_service_filename := config.Get("templates.service.filepath").(string)
	metamorph_servicefileDestpath := path.Join(iso_custom_scripts_path, metamorph_service_filename)

	err = CopyfileToDestination(metamorph_servicefilesourepath, metamorph_servicefileDestpath)

	if err != nil {
                errMessage = "Failed to copy metamorph service file with error"
		logger.Log.Error(errMessage, zap.Error(err))
		return fmt.Errorf(errMessage + "%+v", err)
	}

	//init.sh

	bmhnode.ProvisioningIP = config.Get("provisioning.ip").(string)
	bmhnode.ProvisionerPort = config.Get("provisioning.port").(int)
	bmhnode.HTTPPort = config.Get("provisioning.httpport").(int)

	err = bmhnode.CreateFileFromTemplate(iso_custom_scripts_path, "init")

	if err != nil {
                errMessage = "Failed to create init file with error"
		logger.Log.Error(errMessage, zap.Error(err))
		return fmt.Errorf(errMessage +  "%+v", err)
	}

	err = bmhnode.RepackageISO(iso_DestinationFullpath)

	return err

}

func ValidateChecksum(checksumURL string, iso_path string) error {
	resp, err := http.Get(checksumURL)
	if err != nil {
                logger.Log.Error("Failed to retrieve checksum URL", zap.String("CheckSum URL", checksumURL), zap.Error(err))
		return fmt.Errorf("Failed to retrieve checksum URL %v. Failed with error %+v", checksumURL, err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	iso_checksum := strings.TrimSuffix(string(body), "\n")

	calculated_iso_checksum := Checksum(iso_path)

	if iso_checksum != calculated_iso_checksum {
                logger.Log.Error("Checksum validation failed. Expected checksum does not match  Calculated checksum",
                zap.String("Expected checksum", iso_checksum), zap.String("Calculated Checksum", calculated_iso_checksum))
		return fmt.Errorf("Checksum validation failed. Expected checksum : %+v , Calculated checksum : %v",
			iso_checksum, calculated_iso_checksum)
	}
        logger.Log.Info("Checksum validation successful")

	return nil

}

func ExtractAndCopyISO(iso_DownloadFullpath string, iso_DestinationFullpath string) error {

	//TODO : Where should the temp directories be created ?
	mount_path, err := ioutil.TempDir("/tmp", "iso-")
        var errMessage string 

	if err != nil {
                errMessage = "Failed to create temp directory with error :"
                logger.Log.Error(errMessage, zap.Error(err))
		return fmt.Errorf(errMessage + " %+v", err)
	}
	defer os.RemoveAll(mount_path)

	err = CreateDirectory(iso_DestinationFullpath)

	if err != nil {
                logger.Log.Error("Failed to create destination directory", zap.String("ISO Destination Full path", iso_DestinationFullpath), zap.Error(err))
		return fmt.Errorf("Failed to create destination directory %v with error : %+v", iso_DestinationFullpath, err)
	}

	err = ExtractIso(iso_DownloadFullpath, mount_path)
	if err != nil {
                logger.Log.Error("Failed to extract with error", zap.String("Download  Full Path", iso_DownloadFullpath), zap.Error(err))
		return fmt.Errorf("Failed to extract %v with error : %+v", iso_DestinationFullpath, err)
	}

	cmd := exec.Command("rsync", "-a", mount_path+"/", iso_DestinationFullpath+"/")
	err = cmd.Run()
	if err != nil {
                logger.Log.Error("Faild to copy iso ", zap.String("Destination Full path", iso_DestinationFullpath), zap.Error(err))
		return fmt.Errorf("Faild to copy iso to %v. Error : %+v", iso_DestinationFullpath, err)
	}

	cmd = exec.Command("umount", mount_path)
	err = cmd.Run()
	if err != nil {
		//Ignore this error as it is inconsequential
                logger.Log.Warn("Failed to Unmount. Ignoring error", zap.String("Mount Path", mount_path))
		fmt.Errorf("Failed to Unmount %v. Ignoring error", mount_path)
	}
	return nil
}

func (bmhnode *BMHNode) RepackageISO(iso_DestinationFullpath string) error {

	logger.Log.Info("RepackageISO()", zap.String("isoDestinationFull Path", iso_DestinationFullpath))

	image_name := bmhnode.NodeUUID.String() + "-ubuntu.iso"

	HTTPRootPath := config.Get("http.rootpath").(string)
	cmd := exec.Command(
		"mkisofs",
		"-r",
		"-V",
		"Custom Ubuntu Install CD",
		"-cache-inodes",
		"-J",
		"-l",
		"-b",
		"isolinux/isolinux.bin",
		"-c",
		"isolinux/boot.cat",
		"-no-emul-boot",
		"-boot-load-size",
		"4",
		"-boot-info-table",
		"-o",
		HTTPRootPath+"/"+image_name,
		iso_DestinationFullpath,
	)
	fmt.Println(cmd)
	if err := cmd.Run(); err != nil {
		logger.Log.Error("Error in Creating ISO",zap.Error(err))
		return fmt.Errorf("Failed creating iso image with error : %v", err)
	}

	imageURL := "http://" + config.Get("provisioning.ip").(string) + ":" +
		strconv.Itoa(config.Get("provisioning.httpport").(int)) + "/" + image_name

	// Update the DB with ImageURL
	node.Update(&node.Node{ImageURL: imageURL})

	// TODO : PROVISIONING IP etc moved to config ?
	return nil
}
