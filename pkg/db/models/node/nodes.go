package node

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/manojkva/metamorph-plugin/pkg/config"
	"github.com/manojkva/metamorph-plugin/pkg/logger"
	//	"github.com/bm-metamorph/MetaMorph/pkg/drivers/redfish"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"go.uber.org/zap"
)

func getDB() *gorm.DB {
	dbPath := config.Get("database.path")
	db, err := gorm.Open("sqlite3", dbPath)
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(
		&Node{},
		&NameServer{},
		&Partition{},
		&Filesystem{},
		&KvmPolicy{},
		&SSHPubKey{},
		&BondInterface{},
		&BondParameter{},
		&VirtualDisk{},
		&PhysicalDisk{},
		&BootAction{},
		&Firmware{},
		&Plugins{},
		&API{},
	)
	return db
}

// Only for Controller, Internal
//Functions. Will return all Details
//About node. Including Credentials
func GetNodes() ([]Node, error) {

	logger.Log.Info("GetNodes()")

	nodes := []Node{}
	db := getDB()
	defer db.Close()
	//db.Find(&nodes)
	db.Not("state", []string{"failed", "in-transition", "deploying", "readywait", "setupreadywait"}).Find(&nodes)
	if len(nodes) > 0 {
		return nodes, nil
	} else {
		logger.Log.Error("Nodes not found")
		return nil, errors.New("Nodes not found")
	}
}

func GetNameServers(node_uuid string) ([]NameServer, error) {
	logger.Log.Info("GetNameServers()")
	node := Node{}
	nameservers := []NameServer{}
	db := getDB()
	defer db.Close()
	node_uuid1, _ := uuid.Parse(node_uuid)
	db.Where("node_uuid = ?", node_uuid1).First(&node)
	db.Model(&node).Related(&nameservers)
	if len(nameservers) > 0 {
		return nameservers, nil
	} else {
		logger.Log.Error("No record found")
		return nil, errors.New(" No record Found")
	}
}

func GetPartitions(node_uuid string) ([]Partition, error) {
	logger.Log.Info("GetPartitions()")
	node := Node{}
	partitions := []Partition{}
	db := getDB()
	defer db.Close()
	node_uuid1, _ := uuid.Parse(node_uuid)
	db.Where("node_uuid = ?", node_uuid1).First(&node)
	db.Model(&node).Related(&partitions)
	if len(partitions) > 0 {
		return partitions, nil
	} else {
		logger.Log.Error("No record found")
		return nil, errors.New(" No record Found")
	}
}

func GetSSHPubKeys(node_uuid string) ([]SSHPubKey, error) {
	logger.Log.Info("GetSSHPubKeys()")
	node := Node{}
	sshPubKeys := []SSHPubKey{}
	db := getDB()
	defer db.Close()
	node_uuid1, _ := uuid.Parse(node_uuid)
	db.Where("node_uuid = ?", node_uuid1).First(&node)
	db.Model(&node).Related(&sshPubKeys)
	if len(sshPubKeys) > 0 {
		return sshPubKeys, nil
	} else {
		logger.Log.Error("No record found")
		return nil, errors.New(" No record Found")
	}
}

func GetBondInterfaces(node_uuid string) ([]BondInterface, error) {
	logger.Log.Info("GetBondInterfaces()")
	node := Node{}
	bondInterfaces := []BondInterface{}
	db := getDB()
	defer db.Close()
	node_uuid1, _ := uuid.Parse(node_uuid)
	db.Where("node_uuid = ?", node_uuid1).First(&node)
	db.Model(&node).Related(&bondInterfaces)
	if len(bondInterfaces) > 0 {
		return bondInterfaces, nil
	} else {
		logger.Log.Error("No record found")
		return nil, errors.New("No record Found")
	}
}
func GetFirmwares(node_uuid string) ([]Firmware, error) {
	logger.Log.Info("GetFirmwares()")
	node := Node{}
	firmwares := []Firmware{}
	db := getDB()
	defer db.Close()
	node_uuid1, _ := uuid.Parse(node_uuid)
	db.Where("node_uuid = ?", node_uuid1).First(&node)
	db.Model(&node).Related(&firmwares)
	if len(firmwares) > 0 {
		return firmwares, nil
	} else {
		logger.Log.Error("No record found")
		return nil, errors.New("No record Found")
	}
}

//Plugins

func GetPlugins(node_uuid string) (*Plugins, error) {
	logger.Log.Info("GetPlugins()")
	node := Node{}
	plugins := Plugins{}
	db := getDB()
	defer db.Close()
	node_uuid1, _ := uuid.Parse(node_uuid)
	db.Where("node_uuid = ?", node_uuid1).First(&node)
	err := db.Model(&node).Related(&plugins).Error
	if err == nil {
		return &plugins, err
	} else {
		logger.Log.Error("No record found", zap.Error(err))
		return nil, errors.New("No record Found")
	}
}

func GetPluginAPIs(pluginID uint) ([]API, error) {
	logger.Log.Info("GetPluginAPIs()")
	plugins := Plugins{}
	apis := []API{}
	db := getDB()
	defer db.Close()
	db.Where("id = ?", pluginID).First(&plugins)
	db.Model(&plugins).Related(&apis)
	if len(apis) > 0 {
		return apis, nil
	} else {
		logger.Log.Error("No record found")
		return nil, errors.New(" No record Found")
	}

}

func GetPluginForAPI(apis []API, apiName string)string{
	logger.Log.Info("GetPluginForAPI()")
	for _, api := range apis{
		if strings.ToLower(api.Name) == strings.ToLower(apiName) {
			return api.Plugin
		}
	}
	return ""

}

//VirtualDisk

func GetVirtualDisks(node_uuid string) ([]VirtualDisk, error) {
	logger.Log.Info("GetVirtualDisks()")
	node := Node{}
	virtualdisks := []VirtualDisk{}
	db := getDB()
	defer db.Close()
	node_uuid1, _ := uuid.Parse(node_uuid)
	db.Where("node_uuid = ?", node_uuid1).First(&node)
	db.Model(&node).Related(&virtualdisks)
	if len(virtualdisks) > 0 {
		return virtualdisks, nil
	} else {
		logger.Log.Error("No record found")
		return nil, errors.New(" No record Found")
	}
}

func GetBootActions(node_uuid string) ([]byte, error) {
	logger.Log.Info("GetBootActions()")
	node := Node{}
	bootactions := []BootAction{}

	db := getDB()
	defer db.Close()
	node_uuid1, _ := uuid.Parse(node_uuid)
	db.Where("node_uuid = ?", node_uuid1).First(&node)
	db.Model(&node).Where("status = ?", "new").Order("priority").Related(&bootactions)
	if len(bootactions) > 0 {
		res, _ := json.Marshal(bootactions)
		return res, nil
	} else {
		logger.Log.Error("No record found")
		return nil, errors.New("No Record Found")
	}
}

func GetPhysicalDisks(virtualDiskID uint) ([]PhysicalDisk, error) {
	logger.Log.Info("GetPhysicalDisks()")

	vdisk := VirtualDisk{}
	physcical_disks := []PhysicalDisk{}
	db := getDB()
	defer db.Close()
	db.Where("id = ?", virtualDiskID).First(&vdisk)
	db.Model(&vdisk).Related(&physcical_disks)
	if len(physcical_disks) > 0 {
		return physcical_disks, nil
	} else {
		logger.Log.Error("No record found")
		return nil, errors.New(" No record Found")
	}
}

func GetBondParameters(node_uuid string) ([]BondParameter, error) {
	logger.Log.Info("GetBondParameters()")
	node := Node{}
	bondParameters := []BondParameter{}
	db := getDB()
	defer db.Close()
	node_uuid1, _ := uuid.Parse(node_uuid)
	db.Where("node_uuid = ?", node_uuid1).First(&node)
	db.Model(&node).Related(&bondParameters)
	if len(bondParameters) == 0 {
		logger.Log.Error("No record found")
		return nil, errors.New(" No record Found")
	} else {
		return bondParameters, nil
	}
}

func GetKvmPolicy(node_uuid string) (*KvmPolicy, error) {
	logger.Log.Info("GetKvmPolicy()")
	node := Node{}
	kvmPolicy := KvmPolicy{}
	db := getDB()
	defer db.Close()
	node_uuid1, _ := uuid.Parse(node_uuid)
	db.Where("node_uuid = ?", node_uuid1).First(&node)
	db.Model(&node).Related(&kvmPolicy)
	if kvmPolicy == (KvmPolicy{}) {
		logger.Log.Error("No record found")
		return nil, errors.New(" No record Found")
	} else {
		return &kvmPolicy, nil
	}
}

func GetFilesystem(partitionId uint) (*Filesystem, error) {
	logger.Log.Info("GetFilesystem()")
	partition := Partition{}
	filesystem := Filesystem{}
	db := getDB()
	defer db.Close()
	db.Where("id = ?", partitionId).First(&partition)
	db.Model(&partition).Related(&filesystem)
	if filesystem == (Filesystem{}) {
		logger.Log.Error("No record found")
		return nil, errors.New(" No record Found")
	} else {
		return &filesystem, nil
	}
}

func Describe(node_uuid string) ([]byte, error) {
	logger.Log.Info("Describe()")
	node := Node{}
	db := getDB()
	defer db.Close()

	node_uuid1, _ := uuid.Parse(node_uuid)
	db.Where("node_uuid = ?", node_uuid1).First(&node)
	if node.NodeUUID.String() == node_uuid {
		fmt.Println(node)
		res, _ := json.Marshal(node)
		return res, nil
	} else {
		logger.Log.Error("Node not found", zap.String("NodeUUID",node_uuid))
		return nil, errors.New("Node not found")
	}
}

func Delete(node_uuid string) error {
	logger.Log.Info("Delete()")
	node := Node{}
	db := getDB()
	defer db.Close()

	node_uuid1, _ := uuid.Parse(node_uuid)
	err := db.Where("node_uuid = ?", node_uuid1).First(&node).Error
	if err == nil {
		err = db.Delete(&node).Error
	}
	return err

}

/*
func Update(node *Node) error {
	db := getDB()
	defer db.Close()
	err := db.Save(node).Error
	return err
}
*/
/*
func UpdateRaw(node_uuid string, data []byte )error{
	var node Node
	err := json.Unmarshal(data, &node)
	if err == nil {
		err = Update(node_uuid, &node)
	}
	return  err
}

*/
func Update(updateNode *Node) error {
	logger.Log.Info("Update()")
	node := Node{}
	db := getDB()
	defer db.Close()

	db.Where("node_uuid = ?", updateNode.NodeUUID).First(&node)
	err := db.Model(&node).Updates(updateNode).Error
	return err

}

func UpdateTaskStatus(task *BootAction) error {
	logger.Log.Info("UpdateTaskStatus()")
	db := getDB()
	defer db.Close()
	err := db.Save(task).Error
	return err
}

func Create(data []byte) (string, error) {
	logger.Log.Info("Create()")
	db := getDB()
	defer db.Close()

	var node Node
	//var uuidString string

	UUID, err := uuid.NewRandom()
	err = json.Unmarshal(data, &node)
	/*
		//Get UUID using Redfish Library.
		err := json.Unmarshal(data, &node)
		if err == nil {
			uuidString, _ := redfish.GetUUID(node.IPMIIP, node.IPMIUser, node.IPMIPassword)
			if uuidString == "" {

				err = errors.New(fmt.Sprintf("Failed to retreive Node UUID for nodename : %v", node.Name))
			}
		}
		UUID, err := uuid.Parse(uuidString)
	*/

	if err == nil {
		node.NodeUUID = UUID
		node.State = "new"
		err = db.Create(&node).Error
	}
	if err != nil {
		logger.Log.Error("Failed to create Node", zap.Error(err))
		return "", err
	} else {
		return node.NodeUUID.String(), nil
	}
}

// For Testing purpose only
func CreateTestNode() *Node {
	data, err := ioutil.ReadFile(config.Get("testing.inputfile").(string))
        if err != nil {
          fmt.Println(err)
        }
	uuid, err := Create(data)
        if err != nil {
          fmt.Println(err)
        }
	nodelist, err := GetNodes()
	if err != nil {
		return nil
	}
	for _, node := range nodelist {
		if node.NodeUUID.String() == uuid {
			return &node
		}
	}

	return nil

}
