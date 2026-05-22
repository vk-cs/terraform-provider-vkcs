package servers

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
	paginationutil "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/pagination"
)

type Status string

const (
	StatusDiscovered Status = "DISCOVERED"
	StatusInProgress Status = "IN_PROGRESS"
	StatusActive     Status = "ACTIVE"
)

type BootOrderListItem struct {
	BootDeviceType string `json:"BootDeviceType"`
}

type LocalDiskInfo struct {
	Path   string  `json:"path"`
	SizeGb int64   `json:"sizeGb"`
	Type   string  `json:"type"`
	Model  *string `json:"model"`
}

type Server struct {
	ServerId          string               `json:"serverId"`
	ServerName        string               `json:"serverName"`
	ImageId           *string              `json:"imageId"`
	ImageSource       *string              `json:"imageSource"`
	ImageName         *string              `json:"imageName"`
	OsType            *string              `json:"osType"`
	Status            Status               `json:"status"`
	PowerState        string               `json:"powerState"`
	CpuTypes          []string             `json:"cpuTypes"`
	CpuCores          []int64              `json:"cpuCores"`
	MemoryMegabytes   int64                `json:"memoryMegabytes"`
	LocalDiskSizes    []int64              `json:"localDiskSizes"`
	AvailabilityZone  string               `json:"availabilityZone"`
	Tags              []string             `json:"tags"`
	IsLocked          bool                 `json:"isLocked"`
	TargetBootOrder   []*BootOrderListItem `json:"targetBootOrder"`
	ImageUsername     *string              `json:"imageUsername"`
	RaidType          *string              `json:"raidType"`
	FlavorId          *string              `json:"flavorId"`
	ProvisionProgress *int                 `json:"provisionProgress"`
	LocalDisksInfo    []*LocalDiskInfo     `json:"localDisksInfo"`
}

type commonServerResult struct {
	gophercloud.Result
}

func (r commonServerResult) Extract() (*Server, error) {
	var s Server
	err := r.ExtractInto(&s)
	return &s, err
}

// GetResult represents result of baremetal server get
type GetResult struct {
	commonServerResult
}

type UpdateResult struct {
	gophercloud.ErrResult
}

type DeleteResult struct {
	gophercloud.ErrResult
}

type ProvisionResult struct {
	gophercloud.ErrResult
}

// Page represents a page of baremetal server.
type Page struct {
	paginationutil.TokenPageBase
}

func (p Page) ExtractInto(v interface{}) error {
	return p.PageResult.ExtractInto(v)
}

func ExtractServers(p pagination.Page) ([]Server, error) {
	var s struct {
		Items []Server `json:"items"`
	}
	err := p.(Page).ExtractInto(&s)
	return s.Items, err
}
