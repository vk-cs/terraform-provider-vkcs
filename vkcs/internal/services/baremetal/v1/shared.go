package v1

type ProvisionType string

const (
	ProvisionTypeIMAGE ProvisionType = "IMAGE"
	ProvisionTypeCDROM ProvisionType = "CDROM"
	ProvisionTypeNOOS  ProvisionType = "NO_OS"
)

type ImageSource string

const (
	ImageSourceGLANCE ImageSource = "GLANCE"
	ImageSourcePUBLIC ImageSource = "PUBLIC"
)

type ProvisionFields struct {
	ProvisionType     ProvisionType             `json:"provisionType,omitempty"`
	ImageId           *string                   `json:"imageId,omitempty"`
	ImageSource       ImageSource               `json:"imageSource,omitempty"`
	KeypairName       string                    `json:"keypairName,omitempty"`
	UserData          *string                   `json:"userData,omitempty"`
	Monitoring        *bool                     `json:"monitoring,omitempty"`
	NetworkInterfaces []*NetworkInterfaceConfig `json:"networkInterfaces,omitempty"`
	Bonds             []*BondConfig             `json:"bonds,omitempty"`
	StorageLayout     *StorageLayout            `json:"storageLayout,omitempty"`
}

// StorageLayout is the declarative custom disk layout: disks carry their own
// partitions, raids are assembled from whole disks. A raid member disk must
// not declare partitions of its own.
type StorageLayout struct {
	Disks []*StorageDisk `json:"disks"`
	Raids []*StorageRaid `json:"raids,omitempty"`
}

type StorageDisk struct {
	Id         string              `json:"id"`
	Type       string              `json:"type"`
	SizeGib    int64               `json:"sizeGib"`
	Partitions []*StoragePartition `json:"partitions,omitempty"`
}

// StoragePartition carries its own mount and filesystem type. A nil or empty
// mount means the partition is created but not mounted; both representations
// are passed through to the API as supplied.
type StoragePartition struct {
	Mount  *string `json:"mount,omitempty"`
	Fstype string  `json:"fstype"`
	Size   string  `json:"size"`
}

type StorageRaid struct {
	Id         string              `json:"id"`
	Type       string              `json:"type"`
	Members    []string            `json:"members"`
	Partitions []*StoragePartition `json:"partitions"`
}

type NetworkInterfaceConfig struct {
	NicName string        `json:"nicName"`
	Vlans   []*VlanConfig `json:"vlans"`
}

type BondConfig struct {
	BondName       string        `json:"bondName"`
	InterfaceNames []string      `json:"interfaceNames"`
	Vlans          []*VlanConfig `json:"vlans"`
}

type VlanConfig struct {
	VlanId    *int64 `json:"vlanId"`
	IsNative  bool   `json:"isNative"`
	NetworkId string `json:"networkId"`
	SubnetId  string `json:"subnetId"`
}
