package node

import (
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/jinzhu/gorm"
	"github.com/google/uuid"
)


type Node struct {
	gorm.Model
	NodeUUID             uuid.UUID
	Name                 string
	ImageURL             string
	OamIP                string
	OamGateway           string
	NameServers          []NameServer
	OsDisk               string
	Partitions           []Partition
	GrubConfig           string
	KvmPolicy            KvmPolicy
	SSHPubKeys           []SSHPubKey
	BondInterfaces       []BondInterface
	BondParameters       BondParameters
	IPMIUser             string
	IPMIPassword         string
	Vendor               string
	ServerModel          string
	biosVersion          string
	CPLDFirmwareVersion  string
	RAIDFirmwareVersion  string
	FirmwareVersion      string
	VirtualDisks        []VirtualDisk
  
  }
  
  
  type NameServer struct {
	gorm.Model
	NodeID      uint
	NameServer  string
  }
  
  type Partition struct {
	gorm.Model
	NodeID      uint
	Name        string
	Size        string
	Filesystem Filesystem
  }
  
  type Filesystem struct {
	gorm.Model
	PartitionID   uint
	Mountpoint    string
	Fstype        string
	MountOptions  string
  }
  
  type KvmPolicy struct {
	gorm.Model
	NodeID              uint
	CpuAllocation       string
	CpuPinning          string
	CpuHyperthreading   string 
  }
  
  type SSHPubKey struct {
	gorm.Model
	NodeID     uint
	SSHPubKey  string
  }
  
  type BondInterface struct {
	gorm.Model
	NodeID         uint
	BondInterface  string
  }
  
  type BondParameters struct {
	gorm.Model
	NodeID         uint
	Mode           string
	LacpRate       string
  }
  
  type VirtualDisk struct {
	gorm.Model
	NodeID         uint
	DiskName       string
	raidType        int
	PhysicalDisks []PhysicalDisk
  }
  
  type PhysicalDisk struct {
	gorm.Model
	VirtualDiskID uint
	PhysicalDisk  string
  }