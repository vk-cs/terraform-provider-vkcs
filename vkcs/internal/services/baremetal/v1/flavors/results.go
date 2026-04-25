package flavors

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
	paginationutil "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/pagination"
)

type DiskType string

const (
	DiskTypeSSD DiskType = "SSD"
	DiskTypeHDD DiskType = "HDD"
)

type Disk struct {
	Size int64    `json:"size"`
	Type DiskType `json:"type"`
}

type NetworkInterface struct {
	Name string `json:"name"`
}

type Flavor struct {
	FlavorId           string              `json:"flavorId"`
	FlavorName         string              `json:"flavorName"`
	Status             string              `json:"status"`
	CpuModel           string              `json:"cpuModel"`
	CpuCores           int64               `json:"cpuCores"`
	RamGb              int64               `json:"ramGb"`
	Disks              []*Disk             `json:"disks"`
	NetworkInterfaces  []*NetworkInterface `json:"networkInterfaces"`
	ServersCountByAZ   map[string]*int64   `json:"serversCountByAZ"`
	ServersCount       int64               `json:"serversCount"`
	BondAndVlanCapable bool                `json:"bondAndVlanCapable"`
}

type commonFlavorResult struct {
	gophercloud.Result
}

func (r commonFlavorResult) Extract() (*Flavor, error) {
	var f Flavor
	err := r.ExtractInto(&f)
	return &f, err
}

// GetResult represents result of baremetal flavor get.
type GetResult struct {
	commonFlavorResult
}

// Page represents a page of baremetal flavors.
type Page struct {
	paginationutil.TokenPageBase
}

func ExtractFlavors(p pagination.Page) ([]Flavor, error) {
	var s struct {
		Items []Flavor `json:"items"`
	}
	err := p.(Page).ExtractInto(&s)
	return s.Items, err
}
