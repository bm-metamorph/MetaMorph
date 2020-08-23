package isogen

import (
	"fmt"
	"html/template"
	"os"

	"io/ioutil"
	"path"
	"regexp"
	"strconv"

	"encoding/base64"

	"github.com/bm-metamorph/MetaMorph/pkg/db/models/node"
	config "github.com/manojkva/metamorph-plugin/pkg/config"
	"github.com/manojkva/metamorph-plugin/pkg/logger"
	"go.uber.org/zap"
)

type BMHNode struct {
	*node.Node
}

func (bmhnode *BMHNode) CreateFileFromString(inputString string, outputdir string, modulename string) error {
	logger.Log.Info("CreateFileFromString()")
	if inputString == "" {
		logger.Log.Error("Input String empty")
		return fmt.Errorf("Input String is empty")
	}
	decodedStringInBytes, result := IsBase64(inputString)
	if result != true {
		logger.Log.Error("Input String is not Base64 encoded")
		return fmt.Errorf("Input String is not Base64 encoded")
	}

	filepath := config.Get("templates." + modulename + ".filepath").(string)
	outputfilepathAbsolute := path.Join(outputdir, filepath)

	//Write stirng to file
	err := ioutil.WriteFile(outputfilepathAbsolute, decodedStringInBytes, 0644)
	return err

}

func (bmhnode *BMHNode) CreateNetplanFileFromTemplate(outputdir string, modulename string) error {
	logger.Log.Info("CreateNetplanFileFromTemplate()")
	var err error
	interfacelist, err := node.GetBondInterfaces(bmhnode.NodeUUID.String())
	if err != nil {
		logger.Log.Error("Failed to get Bond Interfaces", zap.Error(err))
		return err
	}
	nameserverlist, err := node.GetNameServers(bmhnode.NodeUUID.String())
	if err != nil {
		logger.Log.Error("Feiled to get Name Server details", zap.Error(err))
		return err
	}

	bondParameters, err := node.GetBondParameters(bmhnode.NodeUUID.String())
	if err != nil {
		logger.Log.Error("Failed to get Bond Parameters", zap.Error(err))
		return err
	}

	bmhnode.BondInterfaces = interfacelist
	bmhnode.NameServers = nameserverlist
	bmhnode.BondParameters = bondParameters

	err = bmhnode.CreateFileFromTemplate(outputdir, modulename)

	return err

}

func (bmhnode *BMHNode) CreatePressedFileFromTemplate(outputdir string, modulename string) error {
	logger.Log.Info("CreatePressedFileFromTemplate()")

	var err error

	partitionlist, err := node.GetPartitions(bmhnode.NodeUUID.String())
	var filesystem *node.Filesystem
	if err == nil {
		for index, part := range partitionlist {
			filesystem, err = node.GetFilesystem(part.ID)
			if err != nil {
				logger.Log.Error("Failed to get FileSystem Info from preseed template", zap.Error(err))
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
	logger.Log.Info("CreateFileFromTemplate()")

	var err error

	logger.Log.Debug(fmt.Sprintf("Creating " + modulename + " from Template"))

	template_rootpath := config.Get("templates.rootdir").(string)

	templatepath := config.Get("templates." + modulename + ".config").(string)
	filepath := config.Get("templates." + modulename + ".filepath").(string)

	templatepathAbsolute := path.Join(template_rootpath, templatepath)
	outputfilepathAbsolute := path.Join(outputdir, filepath)

	if _, err = os.Stat(templatepathAbsolute); os.IsNotExist(err) {
		logger.Log.Error(fmt.Sprintf("Template file for "+modulename+"does not exist : %v\n", err))
		return err
	}
	if _, err = os.Stat(path.Dir(outputfilepathAbsolute)); os.IsNotExist(err) {
		logger.Log.Error(fmt.Sprintf("Output file directory for "+modulename+"does not exist : %v\n", err))
		return err
	}
	tmpl, err := template.ParseFiles(templatepathAbsolute)

	if err != nil {
		logger.Log.Error("Failed to Parse template file", zap.Error(err))
		return err
	}

	f, err := os.Create(outputfilepathAbsolute)
	if err != nil {
		logger.Log.Error(fmt.Sprintf("Failed to create file: %v\n", err))
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
	decodedString, err := base64.StdEncoding.DecodeString(s)
	return decodedString, err == nil
}
