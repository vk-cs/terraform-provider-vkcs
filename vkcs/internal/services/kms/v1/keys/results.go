package keys

import "github.com/gophercloud/gophercloud"

// ListResult is the result of a list request.
// Call its Extract method to interpret a result as a Keys.
type ListResult struct {
	gophercloud.Result
}

type Keys struct {
	Data struct {
		Keys []string `json:"keys"`
	} `json:"data"`
}

// Extract interprets a get result as a list of keys.
func (r ListResult) Extract() ([]string, error) {
	var k Keys
	err := r.ExtractInto(&k)
	return k.Data.Keys, err
}

// GetResult is the result of a get request.
// Call its Extract method to interpret a result as a GetKey.
type GetResult struct {
	gophercloud.Result
}

type GetKey struct {
	Data struct {
		Name            string `json:"name"`
		Type            string `json:"type"`
		DeletionAllowed bool   `json:"deletion_allowed"`
	} `json:"data"`
}

// Extract interprets a get result as a GetKey.
func (r GetResult) Extract() (GetKey, error) {
	var k GetKey
	err := r.ExtractInto(&k)
	return k, err
}

// GetResult is the result of a create request.
// Call its Extract method to interpret a result as a CreateKey.
type CreateResult struct {
	gophercloud.Result
}

type CreateKey struct {
	Data struct {
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"data"`
}

// Extract interprets a create result as a CreateKey.
func (r CreateResult) Extract() (CreateKey, error) {
	var k CreateKey
	err := r.ExtractInto(&k)
	return k, err
}

// UpdateResult is the result of an update request.
// Call its Extract method to interpret a result as an UpdateKey.
type UpdateResult struct {
	gophercloud.Result
}

type UpdateKey struct {
	Data struct {
		DeletionAllowed bool `json:"deletion_allowed"`
	} `json:"data"`
}

// Extract interprets an update result as an UpdateKey.
func (r UpdateResult) Extract() (UpdateKey, error) {
	var k UpdateKey
	err := r.ExtractInto(&k)
	return k, err
}

// DeleteResult is the result of a delete request.
// Call its Extract method to interpret a result as an error.
type DeleteResult struct {
	gophercloud.Result
}

// Extract interprets a get result as an error.
func (r DeleteResult) Extract() error {
	return r.Err
}

// EncryptResult is the result of an encrypt request. Call its Extract method
// to interpret a result as a string.
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

// DecryptResult is the result of an decrypt request. Call its Extract method
// to interpret a result as a string.
type DecryptResult struct {
	gophercloud.Result
}

type DecryptData struct {
	Data struct {
		Plaintext string `json:"plaintext"`
	} `json:"data"`
}

func (r DecryptResult) Extract() (string, error) {
	var d DecryptData
	err := r.ExtractInto(&d)
	return d.Data.Plaintext, err
}
