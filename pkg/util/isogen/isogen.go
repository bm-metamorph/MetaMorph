package isogen

import (
  "net/http"
  "fmt"
  "os"
  "os/exec"
 // "gopkg.in/yaml.v2"
  "log"
  //"strings"
  "io/ioutil"
  "io"
  "crypto/md5"
  "path/filepath"
  //"html/template"
)

func CalculateChecksum(file string) string {
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

func CreateDirectory(directoryPath string) error {
	pathErr := os.MkdirAll(directoryPath,0777)
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
		  log.Fatal(err)
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


func CreateISOLinuxConfig(isoinux_txt_cfg_path string) {

	isoinux_txt_cfg_file, _ := filepath.Abs("./pkg/provisioner/redfish/templates/hwe_kernel/isolinux_txt.cfg")
	input, err := ioutil.ReadFile(isoinux_txt_cfg_file)
	if err != nil {
			fmt.Println(err)
			//return false, err
	}

	err = ioutil.WriteFile(isoinux_txt_cfg_path, input, 0644)
	if err != nil {
			fmt.Println("Error creating", isoinux_txt_cfg_path )
			fmt.Println(err)
			//return false, err
	}
}

/*
// func PrepareISO (iso_url string, iso_checksum_url string, user_data string, node *node, storageConfig string, platformConfig string) (){
func PrepareISO (iso_url string, iso_checksum_url string, user_data string, node *node, hp string) (){

	//Check whether Download is needed
		// If needed download, Verify Checksum, Mount and copy
		   // Contents to new temp location
		iso_url_parts := strings.Split(iso_url,"/")
		iso_name := iso_url_parts[ len(iso_url_parts) -1 ]
		//Final Customized ISO
		iso_dest_path := ISORootPath + "/" + iso_name

		//Where to Download ISO
		iso_path :=  ISORootPath + "/isos/" + iso_name
		os.MkdirAll(ISORootPath + "/isos/", os.ModePerm)

	if _, err := os.Stat(iso_path); os.IsNotExist(err) {

      fmt.Println("ISO does not exists")
			DownloadUrl(iso_path, iso_url)
			resp, err := http.Get(iso_checksum_url)
			if err != nil {
				fmt.Println(err)
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			iso_check_sum := strings.TrimSuffix(string(body), "\n")
			calculated_iso_check_sum := Checksum(iso_path)
      if iso_check_sum == calculated_iso_check_sum {
            fmt.Println("Checksum verification Successfull")
      } else {
          fmt.Println("Checksum verification failed")
        }
      }
		  fmt.Println("Iso Exisis")
		//CreateDirectory(mount_path)

  		mount_path, err := ioutil.TempDir("/opt", "iso-")
  		if err != nil {
  			fmt.Println(err)
  		}
  		defer os.RemoveAll(mount_path)

  		CreateDirectory(iso_dest_path)
  		ExtractIso(iso_path, mount_path)
  		fmt.Println("Copying extracted ISO to: " + iso_dest_path)
  		cmd := exec.Command("rsync", "-a", mount_path + "/", iso_dest_path + "/")
  		cmd.Run()
  		cmd = exec.Command("umount",  mount_path)
  		err = cmd.Run()
  		if err != nil {
  			  log.Fatal(err)
  		}

  	// If not needed, Start editing the files in ISO temp Location

      hostProf := preseed.New(hp)

    	iso_custom_scripts_path := iso_dest_path + "/setup/install/"
    	os.MkdirAll(iso_custom_scripts_path, os.ModePerm)
    	preseed_file_path :=  iso_dest_path + "/preseed/hwe-ubuntu-server.seed"
    	preseed_template_path := "./templates/preseed.tmpl"
    	err = hostProf.CreatePreseedfile(node.Name, storageConfig,  preseed_template_path, preseed_file_path )
    	if err != nil {
    		fmt.Println("Error in Creating preseed file")
    		log.Fatal(err)
    	  }
    	  grub_file_path := iso_custom_scripts_path + "grub.conf"
    	  grub_template_path := "./pkg/provisioner/redfish/templates/grub.tmpl"
    	  err = hostProf.CreateGrubfile(platformConfig, grub_file_path, grub_template_path )
    	  if err != nil {
    		  fmt.Println("Error in Creating grub file")
    		  log.Fatal(err)
    		}

    	isoinux_txt_cfg_path := iso_dest_path + "/isolinux/txt.cfg"
    	CreateISOLinuxConfig(isoinux_txt_cfg_path)

    	ud := UserData{}
    	err = yaml.Unmarshal([]byte(user_data), &ud)
    	if err != nil {
    		fmt.Println("Error in Creating isolinux config")
    		log.Fatalf("error: %v", err)
    	}

    	network_file_name_parts := strings.Split(ud.Write_files[0].Path, "/")
    	network_file_name  := network_file_name_parts[len(network_file_name_parts)-1]
    	network_data := ud.Write_files[0].Content

    	err = ioutil.WriteFile(iso_custom_scripts_path + network_file_name , []byte(network_data), 0644)
    	if err != nil {
    		//return false, err
    		fmt.Println("Error in Creatinng network config file")
    		fmt.Println(err)
    	}

    	metamorph_client_service_file, _ := filepath.Abs("./pkg/provisioner/redfish/templates/metamorph-client.service")
    	input, err := ioutil.ReadFile(metamorph_client_service_file)
    	if err != nil {
    			fmt.Println(err)
    			fmt.Println("Error in creating MetaMorph client service file")
    			//return false, err
    	}
    	err = ioutil.WriteFile(iso_custom_scripts_path + "metamorph-client.service", input, 0644)
    	if err != nil {
    			fmt.Println("Error creating", iso_custom_scripts_path + "metamorph-client.service")
    			fmt.Println(err)
    			//return false, err
    	}

    	//os.Getenv("PROVISIONING_IP") + ":" + os.Getenv("PROVISIONER_PORT")

    	config := map[string]string{
    		"NodeUuid"  :         node.UUID,
    		"ProvisioningIP" :    os.Getenv("PROVISIONING_IP"),
    		"ProvisionerPort" :   os.Getenv("PROVISIONER_PORT"),
    		"HttpPort"        :   os.Getenv("HTTP_PORT"),
    	}

    	init_sh_file, _ := filepath.Abs("./pkg/provisioner/redfish/templates/init.sh")
    	t, err := template.ParseFiles(init_sh_file)
    	if err != nil {
    			fmt.Println(err)
    			fmt.Println("Error in parsing init.sh  template file ")
    			//return false, err
    	}
    	dest_file, err := os.Create(iso_custom_scripts_path + "init.sh")
    	err = t.Execute(dest_file, config)
    	if err != nil {
    		log.Print("execute: ", err)
    		fmt.Println("Error in Creating init.sh file ")
    		//return false, err
    	}
    	dest_file.Close()

    	// Repack ISO
    	// Calculate Image Serve URL and update node details
    	image_name := node.UUID + "-ubuntu.iso"
    	cmd = exec.Command(
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
    	 HTTPRootPath + "/" + image_name,
         iso_dest_path,
    	)
    	fmt.Println(cmd)
    	if err := cmd.Run(); err != nil {
    		fmt.Println(err)
    		fmt.Println("Error in Creating ISO")
    		//return false, fmt.Errorf("error creating configdrive iso: %s", err.Error())
    	}

        node.ImageURL = "http://" + os.Getenv("PROVISIONING_IP") + ":" + os.Getenv("HTTP_PORT") + "/" + image_name
    	//fmt.Println(node)
    	UpdateNode(node)

	//return true, nil
}
*/