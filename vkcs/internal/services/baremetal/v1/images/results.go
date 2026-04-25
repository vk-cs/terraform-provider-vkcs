package images

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
	paginationutil "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/pagination"
)

type Image struct {
	ImageId              string `json:"imageId"`
	ImageName            string `json:"imageName"`
	OsType               string `json:"osType"`
	OsVersion            string `json:"osVersion"`
	DiskFormat           string `json:"diskFormat"`
	MinDiskGb            int64  `json:"minDiskGb"`
	Username             string `json:"username"`
	OsLocalizationStatus string `json:"osLocalizationStatus"`
	RaidType             string `json:"raidType"`
}

type commonFlavorResult struct {
	gophercloud.Result
}

func (r commonFlavorResult) Extract() (*Image, error) {
	var i Image
	err := r.ExtractInto(&i)
	return &i, err
}

// GetResult represents result of baremetal image get.
type GetResult struct {
	commonFlavorResult
}

// Page represents a page of baremetal images.
type Page struct {
	paginationutil.TokenPageBase
}

func ExtractImages(p pagination.Page) ([]Image, error) {
	var s struct {
		Items []Image `json:"items"`
	}
	err := p.(Page).ExtractInto(&s)
	return s.Items, err
}
