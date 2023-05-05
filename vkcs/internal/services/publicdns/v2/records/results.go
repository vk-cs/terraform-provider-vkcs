package records

import "github.com/gophercloud/gophercloud"

type commonResult struct {
	gophercloud.Result
	Type string
}

// ExtractA interprets a GetResult, CreateResult or UpdateResult as a RecordA.
// An error is returned if the original call or the extraction failed.
func (r commonResult) ExtractA() (*RecordA, error) {
	var s RecordA
	err := r.ExtractInto(&s)
	return &s, err
}

// ExtractAAAA interprets a GetResult, CreateResult or UpdateResult as a RecordAAAA.
// An error is returned if the original call or the extraction failed.
func (r commonResult) ExtractAAAA() (*RecordAAAA, error) {
	var s RecordAAAA
	err := r.ExtractInto(&s)
	return &s, err
}

// ExtractA interprets a GetResult, CreateResult or UpdateResult as a RecordA.
// An error is returned if the original call or the extraction failed.
func (r commonResult) ExtractCNAME() (*RecordCNAME, error) {
	var s RecordCNAME
	err := r.ExtractInto(&s)
	return &s, err
}

// ExtractMX interprets a GetResult, CreateResult or UpdateResult as a RecordMX.
// An error is returned if the original call or the extraction failed.
func (r commonResult) ExtractMX() (*RecordMX, error) {
	var s RecordMX
	err := r.ExtractInto(&s)
	return &s, err
}

// ExtractNS interprets a GetResult, CreateResult or UpdateResult as a RecordNS.
// An error is returned if the original call or the extraction failed.
func (r commonResult) ExtractNS() (*RecordNS, error) {
	var s RecordNS
	err := r.ExtractInto(&s)
	return &s, err
}

// ExtractSRV interprets a GetResult, CreateResult or UpdateResult as a RecordSRV.
// An error is returned if the original call or the extraction failed.
func (r commonResult) ExtractSRV() (*RecordSRV, error) {
	var s RecordSRV
	err := r.ExtractInto(&s)
	return &s, err
}

// ExtractTXT interprets a GetResult, CreateResult or UpdateResult as a RecordTXT.
// An error is returned if the original call or the extraction failed.
func (r commonResult) ExtractTXT() (*RecordTXT, error) {
	var s RecordTXT
	err := r.ExtractInto(&s)
	return &s, err
}

// CreateResult is the result of a Create request. Call its Extract method
// to interpret the result as a record.
type CreateResult struct {
	commonResult
}

// GetResult is the result of a Get request. Call its Extract method
// to interpret the result as a record.
type GetResult struct {
	commonResult
}

// UpdateResult is the result of a Update request. Call its Extract method
// to interpret the result as a record.
type UpdateResult struct {
	commonResult
}

// DeleteResult is the result of a Delete request. Call its ExtractErr method
// to determine if the request succeeded or failed.
type DeleteResult struct {
	gophercloud.ErrResult
}
