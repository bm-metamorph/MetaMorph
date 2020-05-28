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
        ISOURL               string
        ISOChecksum          string
	ImageURL             string
	ChecksumURL          string
    DisableCertVerification bool
	ImageReadilyAvailable bool
	OamIP                string
	OamGateway           string
	NameServers          []NameServer `json:"NameServers"`
	OsDisk               string
	Partitions           []Partition
	GrubConfig           string
	KvmPolicy            KvmPolicy
	SSHPubKeys           []SSHPubKey
	BondInterfaces       []BondInterface
	BondParameters       []BondParameter
	IPMIIP               string
	IPMIUser             string
	IPMIPassword         string
	Vendor               string
	ServerModel          string
	biosVersion          string
	CPLDFirmwareVersion  string
	RAIDFirmwareVersion  string
	FirmwareVersion      string
	VirtualDisks         []VirtualDisk
	State                string `gorm:"DEFAULT:new"`
	ProvisioningIP       string
	ProvisionerPort      int
	HTTPPort             int
	BootActions          []BootAction
	NetworkConfig        string
        RAID_reset           bool  `gorm:"DEFAULT:true"`
  }

  type BootAction struct {
	  gorm.Model
	  NodeID     uint
	  Name       string
	  Location   string
	  Priority   uint
	  Control    string
	  Args       string
	  Status     string  `gorm:"DEFAULT:new"`
  }
  
  type NameServer struct {
	gorm.Model
	NodeID      uint
	NameServer  string `json:"NameServer"`
  }
  
  type Partition struct {
	gorm.Model
	NodeID      uint
	Name        string
	Size        string
	Bootable    bool
	Primary     bool
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
  
  type BondParameter struct {
	gorm.Model
	NodeID         uint
	Key           string
	Value       string
  }
  
  type VirtualDisk struct {
	gorm.Model
	NodeID         uint
	DiskName       string
	RaidType        int
	RaidController string
	PhysicalDisks []PhysicalDisk
  }
  
  type PhysicalDisk struct {
	gorm.Model
	VirtualDiskID uint
	PhysicalDisk  string
  }
