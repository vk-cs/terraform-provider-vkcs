package keys

import "github.com/gophercloud/gophercloud"

// GetResult is the result of a get request.
// Call its Extract method to interpret a result as a Key.
type GetResult struct {
	gophercloud.Result
}

// // Extract interprets a get result as a ServiceUser.
// func (r GetResult) Extract() (*ServiceUser, error) {
// 	var s ServiceUser
// 	err := r.ExtractInto(&s)
// 	return &s, err
// }

// EncryptResult is the result of an encrypt request. Call its Extract method
// to interpret a result as a Key.
type EncryptResult struct {
	gophercloud.Result
}

type EncryptData struct {
	Data struct {
		Ciphertext string `json:"ciphertext"`
	} `json:"data"`
}

func (r EncryptResult) Extract() (string, error) {
	var d EncryptData
	err := r.ExtractInto(&d)
	return d.Data.Ciphertext, err
}
