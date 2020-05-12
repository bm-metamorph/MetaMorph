package isogen

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	// "gopkg.in/yaml.v2"
	"log"
	//"strings"
	"crypto/md5"
	"io"
	"io/ioutil"

	//"html/template"
	"path"
	"strings"

	config "bitbucket.com/metamorph/pkg/config"
	"bitbucket.com/metamorph/pkg/db/models/node"
)

func CreateDirectory(directoryPath string) error {
	pathErr := os.MkdirAll(directoryPath, 0777)
	if pathErr != nil {
		return pathErr
	}
	return nil
}

func DownloadUrl(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func ExtractIso(iso, target string) error {
	fmt.Println("Extracting ISO")
	cmd := exec.Command("mount", "-r", "-o", "loop", iso, target)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Failed to run mount command %v", err)
	}
	return err
}

func Checksum(file string) string {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}
	checksum := fmt.Sprintf("%x", h.Sum(nil))
	return checksum
}

func CopyfileToDestination(sourcefilepath string, destinationfilepath string) error {

	input, err := ioutil.ReadFile(sourcefilepath)
	if err != nil {
		return fmt.Errorf("Failed to open file %v with error : %v", sourcefilepath, destinationfilepath)
	}
	err = ioutil.WriteFile(destinationfilepath, input, 0644)
	if err != nil {
		return fmt.Errorf("Failed to write %v to destination %v with error %v", sourcefilepath, destinationfilepath, err)
	}
	return nil
}

func (bmhnode *BMHNode) PrepareISO() error {
	//check if ISO generation if required at all

	var err error
	if bmhnode.ImageReadilyAvailable {

		return fmt.Errorf("Customised ISO already available. No need to prepare custom ISO")
	}

	iso_rootpath := config.Get("iso.rootpath").(string)
	HTTPRootPath := config.Get("http.rootpath").(string)

	iso_urlpath := config.Get("image.url").(string)
	iso_checksum := config.Get("image.checksum").(string)

	iso_name_parts := strings.Split(iso_urlpath, "/")
	iso_name := iso_name_parts[len(iso_name_parts)-1]

	if _, err := os.Stat(iso_rootpath); os.IsNotExist(err) {
		return fmt.Errorf("ISO root directory not found : %v\n", err)
	}

	if _, err := os.Stat(HTTPRootPath); os.IsNotExist(err) {
		return fmt.Errorf("HTTP root directory not found : %v\n", err)
	}
	//Get temporary directory for copying ISO

	iso_tempdir := config.Get("iso.tempdir").(string)

	iso_DownloadFullpath := path.Join(iso_tempdir, iso_name)
	iso_DestinationFullpath := path.Join(iso_rootpath, iso_name)

	if _, err := os.Stat(iso_DownloadFullpath); os.IsNotExist(err) {
		err = os.MkdirAll(iso_tempdir,os.ModePerm)
		if err != nil {
			return fmt.Errorf("Failed to create dir %v with error  %v", iso_tempdir, err)
		}
		err = DownloadUrl(iso_DownloadFullpath, iso_urlpath)
		if err != nil {
			return fmt.Errorf("Failed to download ISO Image : %v", err)
		}

	} else {
		fmt.Printf("ISO vanilla Image already downloaded \n")
	}

	err = ValidateChecksum(iso_checksum, iso_DownloadFullpath)

	if err != nil {
		return fmt.Errorf("Failed to validate checksum. Error : %v", err)
	}

	err = ExtractAndCopyISO(iso_DownloadFullpath, iso_DestinationFullpath)

	if err != nil {
		return fmt.Errorf("Failed to Extract and Copy ISO file. Error %v", err)
	}

	//Add Preseed
	err = bmhnode.CreatePressedFileFromTemplate(iso_DestinationFullpath, "preseed")

	if err != nil {
		return fmt.Errorf("Failed to create Preseed file with error %v", err)
	}

	//Add Grub Config

	iso_custom_scripts_path := path.Join(iso_DestinationFullpath, "/setup/install/")

	err = os.MkdirAll(iso_custom_scripts_path, os.ModePerm)

	if err != nil {
		return fmt.Errorf("Failed to create directory %v with error %v", iso_custom_scripts_path, err)
	}

	err = bmhnode.CreateFileFromTemplate(iso_custom_scripts_path, "grub")

	if err != nil {
		return fmt.Errorf("Failed to create Grub file with error %v", err)
	}

	isolinux_cfg_destpath := iso_DestinationFullpath + "/isolinux/txt.cfg"

	metamorph_root := config.Get("templates.rootdir").(string)

	isolinuxtemplatepath := config.Get("templates.isolinux.config").(string)


	isolinux_cfg_sourcepath := path.Join(metamorph_root,isolinuxtemplatepath) 

	err = CopyfileToDestination(isolinux_cfg_sourcepath, isolinux_cfg_destpath)

	if err != nil {
		return fmt.Errorf("Failed to copy isolinux cfg file with error %v", err)
	}
    //Netplan
	err = bmhnode.CreateNetplanFileFromTemplate(iso_custom_scripts_path, "netplan")

	if err != nil {
		return fmt.Errorf("Failed to create Netplan file with error %v", err)
	}

	//MetaMorph Agent
	metamorph_assets_root := config.Get("assets.rootdir").(string)
	metamorph_agent_file_src := config.Get("assets.agent_binary.src").(string)
	metamorph_agent_file_src_abs := path.Join(metamorph_assets_root,metamorph_agent_file_src )

	metamorph_agent_file_dest := config.Get("assets.agent_binary.dest").(string)
	metamorph_agent_file_dest_abs := path.Join(iso_custom_scripts_path,metamorph_agent_file_dest )
	err = CopyfileToDestination(metamorph_agent_file_src_abs, metamorph_agent_file_dest_abs)

	if err != nil {
		return fmt.Errorf("Failed to copy metamorph agent file with error %v", err)
	}

	//MetaMorph Agent config
	bmhnode.ProvisioningIP = config.Get("provisioning.ip").(string)
	err = bmhnode.CreateFileFromTemplate(iso_custom_scripts_path, "agent_config")

	if err != nil {
		return fmt.Errorf("Failed to create Metamorph Agent config file with error %v", err)
	}

	//metamorph-client.service
	metamorph_servicetemplatepath := config.Get("templates.service.config").(string)
	
	metamorph_servicefilesourepath := path.Join(metamorph_root,metamorph_servicetemplatepath ) 

	metamorph_service_filename := config.Get("templates.service.filepath").(string)
	metamorph_servicefileDestpath   := path.Join(iso_custom_scripts_path,metamorph_service_filename)

	err = CopyfileToDestination(metamorph_servicefilesourepath, metamorph_servicefileDestpath)

	if err != nil {
		return fmt.Errorf("Failed to copy metamorph service file with error %v", err)
	}


	//init.sh

	bmhnode.ProvisioningIP = config.Get("provisioning.ip").(string)
	bmhnode.ProvisionerPort = config.Get("provisioning.port").(int)
	bmhnode.HTTPPort = config.Get("provisioning.httpport").(int)

	err = bmhnode.CreateFileFromTemplate(iso_custom_scripts_path, "init")

	if err != nil {
		return fmt.Errorf("Failed to create init file with error %v", err)
	}

	err = bmhnode.RepackageISO(iso_DestinationFullpath)

	return err

}

func ValidateChecksum(checksumURL string, iso_path string) error {
	resp, err := http.Get(checksumURL)
	if err != nil {
		return fmt.Errorf("Failed to retrieve checksum URL %v. Failed with error %v", checksumURL, err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	iso_checksum := strings.TrimSuffix(string(body), "\n")

	calculated_iso_checksum := Checksum(iso_path)

	if iso_checksum != calculated_iso_checksum {
		return fmt.Errorf("Checksum validation failed. Expected checksum : %v , Calculated checksum : %v",
			iso_checksum, calculated_iso_checksum)
	}
	fmt.Println("Checksum validation successful")

	return nil

}

func ExtractAndCopyISO(iso_DownloadFullpath string, iso_DestinationFullpath string) error {

	//TODO : Where should the temp directories be created ?
	mount_path, err := ioutil.TempDir("/tmp", "iso-")

	if err != nil {
		return fmt.Errorf("Failed to create temp directory with error : %v", err)
	}
	defer os.RemoveAll(mount_path)

	err = CreateDirectory(iso_DestinationFullpath)

	if err != nil {
		return fmt.Errorf("Failed to create destination directory %v with error : %v", iso_DestinationFullpath, err)
	}

	err = ExtractIso(iso_DownloadFullpath, mount_path)
	if err != nil {
		return fmt.Errorf("Failed to extract %v with error : %v", iso_DestinationFullpath, err)
	}

	cmd := exec.Command("rsync", "-a", mount_path+"/", iso_DestinationFullpath+"/")
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("Faild to copy iso to %v. Error : %v", iso_DestinationFullpath, err)
	}

	cmd = exec.Command("umount", mount_path)
	err = cmd.Run()
	if err != nil {
		//Ignore this error as it is inconsequential
		fmt.Errorf("Failed to Unmount %v. Ignoring error", mount_path)
	}
	return nil
}

func (bmhnode *BMHNode) RepackageISO(iso_DestinationFullpath string) error {

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
		fmt.Println("Error in Creating ISO")
		return fmt.Errorf("Failed creating iso image with error : %v", err)
	}

	bmhnode.ImageURL = "http://" + config.Get("provisioning.ip").(string) + ":" +
	                               strconv.Itoa(config.Get("provisioning.httpport").(int)) + "/" + image_name

	// Update the DB with ImageURL
	node.Update(bmhnode.Node)

	// TODO : PROVISIONING IP etc moved to config ?
	return nil
}