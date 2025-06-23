package common

import (
	"github.com/google/uuid"
	"github.com/gophercloud/gophercloud"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/randutil"

	"strings"
)

func SetHeaders(client *gophercloud.ServiceClient) {
	client.MoreHeaders = map[string]string{
		"x-req-id":         randutil.RandomRequestID(8),
		"X-MCS-Request-Id": strings.ReplaceAll(uuid.New().String(), "-", ""),
	}
}
