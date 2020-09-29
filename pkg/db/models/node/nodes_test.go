package node

import (
	"fmt"
	"testing"
	 "github.com/jinzhu/gorm"
	"github.com/manojkva/metamorph-plugin/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestGetDB(t *testing.T){
	config.SetLoggerConfig("logger.apipath")
	db:=getDB()
	fmt.Println(db)
	fmt.Printf("%T",db)
	assert.NotEqual(t, db, nil) 
	dbPath := config.Get("database.path")
        db, err := gorm.Open("sqlite3", dbPath)
        if err != nil {
                fmt.Println("failed to connect database")
		t.Error("got",err , "want", nil)
        }
}

func TestGetNodes(t *testing.T){
	config.SetLoggerConfig("logger.apipath")
	node,err:=GetNodes()
	if len(node)<=0 || err!=nil {
		 t.Error("got",len(node) , "want", "greater then 0")
		 t.Error("got",err , "want", nil)
	}


}


func TestGetNameServers(t *testing.T){
	config.SetLoggerConfig("logger.apipath")
	node := CreateTestNode()
	nameserver,err:=GetNameServers(node.NodeUUID.String())

	if len(nameserver)<=0 || err!=nil {
		 t.Error("got",len(nameserver) , "want", "greater then 0")
		 t.Error("got",err , "want", nil)
	}


}

func TestGetPartitions(t *testing.T){
	config.SetLoggerConfig("logger.apipath")
	node := CreateTestNode()
	partitions,err:=GetPartitions(node.NodeUUID.String())

	if len(partitions)<=0 || err!=nil {
		 t.Error("got",len(partitions) , "want", "greater then 0")
		 t.Error("got",err , "want", nil)
	}


}

func TestGetSSHPubKeys(t *testing.T){
	config.SetLoggerConfig("logger.apipath")
	node := CreateTestNode()
	sshPubKeys,err:=GetSSHPubKeys(node.NodeUUID.String())

	if len(sshPubKeys)<=0 || err!=nil {
		 t.Error("got",len(sshPubKeys) , "want", "greater then 0")
		 t.Error("got",err , "want", nil)
	}


}

func TestGetFirmwares(t *testing.T){
	config.SetLoggerConfig("logger.apipath")
	node := CreateTestNode()
	firmwares,err:=GetFirmwares(node.NodeUUID.String())

	if len(firmwares)<=0 || err!=nil {
		 t.Error("got",len(firmwares) , "want", "greater then 0")
		 t.Error("got",err , "want", nil)
	}


}
func TestGetVirtualDisks(t *testing.T){
	config.SetLoggerConfig("logger.apipath")
	node := CreateTestNode()
	virtualdisks,err:=GetVirtualDisks(node.NodeUUID.String())
	fmt.Println("virtualdisks[0]")
	fmt.Println(virtualdisks[0].ID)
	fmt.Printf("%T",virtualdisks[0])

	if len(virtualdisks)<=0 || err!=nil {
		 t.Error("got",len(virtualdisks) , "want", "greater then 0")
		 t.Error("got",err , "want", nil)
	}


}
func TestGetBootActions(t *testing.T){
	config.SetLoggerConfig("logger.apipath")
	node := CreateTestNode()
	bootactions,err:=GetBootActions(node.NodeUUID.String())

	if len(bootactions)<=0 || err!=nil {
		 t.Error("got",len(bootactions) , "want", "greater then 0")
		 t.Error("got",err , "want", nil)
	}


}

func TestGetPhysicalDisks(t *testing.T){
	config.SetLoggerConfig("logger.apipath")
        node := CreateTestNode()
        virtualdisks,err:=GetVirtualDisks(node.NodeUUID.String())
        fmt.Println("virtualdisks[0]")
        fmt.Println(virtualdisks[0].ID)
        fmt.Printf("%T",virtualdisks[0])

	physcical_disks,err:=GetPhysicalDisks(virtualdisks[0].ID)

	if len(physcical_disks)<=0 || err!=nil {
		 t.Error("got",len(physcical_disks) , "want", "greater then 0")
		 t.Error("got",err , "want", nil)
	}


}

func TestDescribe(t *testing.T){
	config.SetLoggerConfig("logger.apipath")
	node := CreateTestNode()
	res,err:=Describe(node.NodeUUID.String())

	if len(res)<=0 || err!=nil {
		 t.Error("got",len(res) , "want", "greater then 0")
		 t.Error("got",err , "want", nil)
	}


}



func init() {
	config.SetLoggerConfig("logger.apipath")
}

func TestGetBondParameters(t *testing.T) {
	config.SetLoggerConfig("logger.apipath")
	node := CreateTestNode()
	bondParameters, _ := GetBondParameters(node.NodeUUID.String())
	fmt.Println("length",len(bondParameters))
	if node != nil {
		if  len(bondParameters) != 6  && len(bondParameters) != 0 {
		t.Error("got",len(bondParameters) , "want", "6 or 0")
	}
	}
}
func TestGetKvmPolicy(t *testing.T) {
	config.SetLoggerConfig("logger.apipath")
	node := CreateTestNode()
	kvmPolicy, err := GetKvmPolicy(node.NodeUUID.String())
	assert.Equal(t, err, nil)
	assert.Equal(t, kvmPolicy.CpuAllocation, "1:1")
	assert.Equal(t, kvmPolicy.CpuPinning, "enabled")
	assert.Equal(t, kvmPolicy.CpuHyperthreading, "enabled")
}
func TestGetFilesystem(t *testing.T) {
	config.SetLoggerConfig("logger.apipath")
	node := CreateTestNode()
	partitions, err := GetPartitions(node.NodeUUID.String())
	assert.Equal(t, err, nil)
	for index, part := range partitions {
		filesystem, _ := GetFilesystem(part.ID)
		partitions[index].Filesystem = *filesystem
	}
	fmt.Printf("%v", partitions[0].Filesystem)
}

func  TestGetPlugins(t * testing.T){
	config.SetLoggerConfig("logger.apipath")
	node := CreateTestNode()
	plugins, err := GetPlugins(node.NodeUUID.String())
	assert.Equal(t, err, nil)
	fmt.Println("PLUGINS")
	fmt.Printf("%+v", plugins)
	apis,err := GetPluginAPIs(plugins.ID)
	assert.Equal(t, err, nil)
	fmt.Println("APIS")
	fmt.Printf("%+v", apis)
}
func TestCreateTestNode(t *testing.T) {
	config.SetLoggerConfig("logger.apipath")
	node := CreateTestNode()
	assert.NotEqual(t, node, nil)
}
