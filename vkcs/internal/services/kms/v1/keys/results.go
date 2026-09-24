package keys

import "github.com/gophercloud/gophercloud"

// GetResult is the result of a get request. Call its Extract method
// to interpret a result as a Key.
type GetResult struct {
	gophercloud.Result
}

// // Extract interprets a get result as a ServiceUser.
// func (r GetResult) Extract() (*ServiceUser, error) {
// 	var s ServiceUser
// 	err := r.ExtractInto(&s)
// 	return &s, err
// }
