package node

import (
	"encoding/json"
	"fmt"
	"errors"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"bitbucket.com/metamorph/pkg/config"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
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
		&BondParameters{},
		&VirtualDisk{},
		&PhysicalDisk{},
		)
	return db
}


func Describe(node_uuid string) ([]byte, error) {
	node := Node{}
	db := getDB()
	defer db.Close()

	node_uuid1 ,_ := uuid.Parse(node_uuid)
	db.Where("node_uuid = ?", node_uuid1).First(&node)
	if node.NodeUUID.String() == node_uuid {
		fmt.Println(node)
		res, _ := json.Marshal(node)
		return res, nil
	} else {
		return  nil, errors.New("Node not found")
	}
}

func Create(data []byte) ( string, error) {
	db := getDB()
	defer db.Close()

	var node Node
	UUID, err := uuid.NewRandom()
	err = json.Unmarshal(data, &node)
	node.NodeUUID = UUID
	err = db.Create(&node).Error
	if err != nil {
		return "", err
	} else {
		return node.NodeUUID.String(), nil
	}
}
