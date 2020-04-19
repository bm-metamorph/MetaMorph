package preseed

import (
	"fmt"
	"io/ioutil"
	"testing"
)

var PathofTmplfile =  "../../../../../pkg/provisioner/redfish/templates/preseed.tmpl"
var PathofGrubTmplfile =  "../../../../../pkg/provisioner/redfish/templates/grub.tmpl"

func TestCreatePreseedfileLVM(t *testing.T) {
	storageYaml, _ := ioutil.ReadFile("hostProfileLVM.yaml")
	CreatePreseedfile("mtn52g001" , string(storageYaml), PathofTmplfile, "test.seed")
}

func TestCreatePreseedfile(t *testing.T) {
	storageYaml, _ := ioutil.ReadFile("hostProfile.yaml")
	CreatePreseedfile("mtn52g001" , string(storageYaml), PathofTmplfile, "test.seed")
}

func TestReadfromYaml(t *testing.T) {
	storageYaml, _ := ioutil.ReadFile("hostProfile.yaml")
	readFromYaml(string(storageYaml))
}

func TestgetDiskSpaceinMB(t *testing.T) {
	dspace, maxdspace, _ := getDiskSpaceinMB(">300g")
	fmt.Printf("%s, %s", dspace, maxdspace)
}

func TestCreateGrubfile(t *testing.T){
        platformYaml, _ := ioutil.ReadFile("platProfile.yaml")
	CreateGrubfile(string(platformYaml), "grub.cfg", PathofGrubTmplfile )
}
