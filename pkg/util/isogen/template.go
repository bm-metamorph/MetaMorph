package isogen

import (
	"fmt"
	"html/template"
	"os"

	"path"
	"regexp"
	"strconv"
        "io/ioutil"

	config "github.com/bm-metamorph/MetaMorph/pkg/config"
	"github.com/bm-metamorph/MetaMorph/pkg/db/models/node"
	"encoding/base64"
)

type BMHNode struct {
	*node.Node
}

func (bmhnode *BMHNode) CreateNetplanFileFromString(outputdir string, modulename string) error {
	var networkConfig = bmhnode.NetworkConfig
	if networkConfig == "" {
		return  fmt.Errorf("NetworkConfig is empty")
	}
	decodedStringInBytes, result := IsBase64(networkConfig)
	if result != true{
		return fmt.Errorf("NetworkConfig is not Base64 encoded")
	}

	filepath := config.Get("templates." + modulename + ".filepath").(string)
	outputfilepathAbsolute := path.Join(outputdir, filepath)

	//Write stirng to file
	err := ioutil.WriteFile(outputfilepathAbsolute,decodedStringInBytes,0644)
	return err

}

func (bmhnode *BMHNode) CreateNetplanFileFromTemplate(outputdir string, modulename string) error {
	var err error
	interfacelist, err := node.GetBondInterfaces(bmhnode.NodeUUID.String())
	if err != nil {
		return err
	}
	nameserverlist, err := node.GetNameServers(bmhnode.NodeUUID.String())
	if err != nil {
		return err
	}

	bondParameters,err := node.GetBondParameters(bmhnode.NodeUUID.String())
	if err != nil {
		return err
	}


	bmhnode.BondInterfaces = interfacelist
	bmhnode.NameServers = nameserverlist
	bmhnode.BondParameters = bondParameters


	err = bmhnode.CreateFileFromTemplate(outputdir, modulename)

	return err

}

func (bmhnode *BMHNode) CreatePressedFileFromTemplate(outputdir string, modulename string) error {

	var err error

	partitionlist, err := node.GetPartitions(bmhnode.NodeUUID.String())
	var filesystem *node.Filesystem
	if err == nil {
		for index,part := range partitionlist{
			filesystem, err = node.GetFilesystem(part.ID)
			if err != nil{
				return err
			}
			partitionlist[index].Filesystem = *filesystem
		} 
		bmhnode.Partitions = partitionlist

		err = bmhnode.CreateFileFromTemplate(outputdir, modulename)
	}
	return err
}
func (bmhnode *BMHNode) CreateFileFromTemplate(outputdir string, modulename string) error {

	var err error

	fmt.Println("Creating " + modulename + " from Template")

	template_rootpath := config.Get("templates.rootdir").(string)

	templatepath := config.Get("templates." + modulename + ".config").(string)
	filepath := config.Get("templates." + modulename + ".filepath").(string)

	templatepathAbsolute := path.Join(template_rootpath, templatepath)
	outputfilepathAbsolute := path.Join(outputdir, filepath)

	if _, err = os.Stat(templatepathAbsolute); os.IsNotExist(err) {
		fmt.Printf("Template file for "+modulename+"does not exist : %v\n", err)
		return err
	}
	if _, err = os.Stat(path.Dir(outputfilepathAbsolute)); os.IsNotExist(err) {
		fmt.Printf("Output file directory for "+modulename+"does not exist : %v\n", err)
		return err
	}
	tmpl, err := template.ParseFiles(templatepathAbsolute)

	if err != nil {
		return err
	}

	f, err := os.Create(outputfilepathAbsolute)
	if err != nil {
		fmt.Printf("Failed to create file: %v\n", err)
		return err
	}

	err = tmpl.Execute(f, bmhnode)

	return err

}

func (bmhnode *BMHNode) GetDiskSizeMB(diskspace string) (string, error) {
	disksizeMB, _, err := getDiskSpaceinMB(diskspace)
	return disksizeMB, err
}
func (bmhnode *BMHNode) GetMaxDiskSizeMB(diskspace string) (string, error) {
	_, maxdiskSizeinMB, err := getDiskSpaceinMB(diskspace)
	return maxdiskSizeinMB, err
}

func getDiskSpaceinMB(diskspace string) (diskspaceinMB string, maxdiskSizeinMB string, err error) {
	//check if there is diskspace listed with a numeber followed by g
	// if there is > in front of the
	re := regexp.MustCompile(`(>*)(\d+)([a-z])`)
	t := re.FindSubmatch([]byte(diskspace))
	if len(t) == 4 {
		disksizeGB, err := strconv.Atoi(string(t[2]))
		disksizeMB := strconv.Itoa(disksizeGB * 1024)
		var maxdiskSizeinMB string
		if string(t[1]) == ">" {
			maxdiskSizeinMB = "-1"
		} else {
			maxdiskSizeinMB = disksizeMB
		}
		return disksizeMB, maxdiskSizeinMB, err
	}
	return "", "", err

}

func IsBase64(s string) ([]byte, bool) {
	decodedString , err := base64.StdEncoding.DecodeString(s)
	return decodedString , err == nil
}
