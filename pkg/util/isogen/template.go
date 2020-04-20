package isogen

import (
	"fmt"
	"html/template"
	"os"

	"path"
	"regexp"
	"strconv"

	config "bitbucket.com/metamorph/pkg/config"
	"bitbucket.com/metamorph/pkg/db/models/node"
)

type BMHNode node.Node

func (bmhnode *BMHNode) CreateFileFromTemplate(modulename string) error {

	var err error

	fmt.Println("Creating " + modulename + " from Template")

	metamorph_root := config.Get("metamorph-root").(string)

	templatepath := config.Get("templates." + modulename + ".config").(string)
	outputfilepath := config.Get("templates." + modulename + ".filepath").(string)

	templatepathAbsolute := path.Join(metamorph_root, templatepath)

	if _, err = os.Stat(templatepathAbsolute); os.IsNotExist(err) {
		fmt.Printf("Template file for "+modulename+"does not exist : %v\n", err)
		return err
	}
	if _, err = os.Stat(path.Dir(outputfilepath)); os.IsNotExist(err) {
		fmt.Printf("Output file for "+modulename+"does not exist : %v\n", err)
		return err
	}
/*
	tmpl := template.New(modulename).Funcs(template.FuncMap{
		"GetDiskSizeMB":    GetDiskSizeMB,
		"GetMaxDiskSizeMB": GetMaxDiskSizeMB,
	})
*/	
	tmpl, err  := template.ParseFiles(templatepathAbsolute)

	if err != nil {
		return err
	}

	f, err := os.Create(outputfilepath)
	if err != nil {
		fmt.Printf("Failed to create file: %v\n", err)
		return err
	}

	err = tmpl.Execute(f, bmhnode)

	return err

}

/*

func (bmhnode *BMHNode) CreatePreseedfile(templatepath string, preeseedpath string) (err error) {

	fmt.Println("Inside CreateaPreseeedfile()")

	//TODO : Check file path vaiidity

	tmpl, err := template.ParseFiles(templatepath)
	if err != nil {
		return err
	}

	f, err := os.Create(preeseedpath)
	if err != nil {
		fmt.Println("Failed to create file: ")
		return err
	}
	err = tmpl.Execute(f, bmhnode)

	return err
}
*/
func (bmhnode *BMHNode)GetDiskSizeMB(diskspace string) (string, error) {
	disksizeMB, _, err := getDiskSpaceinMB(diskspace)
	return disksizeMB, err
}
func (bmhnode * BMHNode) GetMaxDiskSizeMB(diskspace string) (string, error) {
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
