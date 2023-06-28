package images

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/stretchr/testify/assert"
)

func TestResourceImagesImageFile(t *testing.T) {
	image := ResourceImagesImage().TestResourceData()
	imgCachePath := fmt.Sprintf("%s/.terraform/image_cache", os.Getenv("HOME"))
	image.Set("image_source_url", "http://download.cirros-cloud.net/0.3.5/cirros-0.3.5-x86_64-disk.img")
	image.Set("image_cache_path", imgCachePath)

	client := new(gophercloud.ServiceClient)
	client.ProviderClient = new(gophercloud.ProviderClient)

	filename, err := resourceImagesImageFile(client, image)

	assert.Equal(t, nil, err)
	assert.Equal(t, filepath.Join(imgCachePath, "e18f6225aa529b597b260081c6ecb1da.img"), filename)
}
