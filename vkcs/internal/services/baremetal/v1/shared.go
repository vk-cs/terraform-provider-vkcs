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
	NetworkInterfaces []*NetworkInterfaceConfig `json:"networkInterfaces,omitempty"`
	Bonds             []*BondConfig             `json:"bonds,omitempty"`
	RaidType          *string                   `json:"raidType,omitempty"`
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
