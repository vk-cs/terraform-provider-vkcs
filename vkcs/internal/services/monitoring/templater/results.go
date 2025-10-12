package templater

import (
	"encoding/base64"

	"github.com/gophercloud/gophercloud"
)

type ListUsersResult struct {
	gophercloud.Result
}

func (r ListUsersResult) ExtractErr() error {
	return r.Err
}

type Settings struct {
	UserID string `json:"user_id"`
	Script string `json:"script"`
}

type CreateUserResult struct {
	gophercloud.Result
}

func (r CreateUserResult) Extract() (Settings, error) {
	var monitoring Settings
	err := r.ExtractInto(&monitoring)
	if err != nil {
		return monitoring, err
	}

	bytes, err := base64.StdEncoding.DecodeString(monitoring.Script)
	if err != nil {
		return monitoring, err
	}

	monitoring.Script = string(bytes)

	return monitoring, nil
}
