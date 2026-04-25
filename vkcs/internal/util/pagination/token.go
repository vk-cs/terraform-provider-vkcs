package pagination

import (
	"github.com/gophercloud/gophercloud/pagination"
)

type TokenPageBase struct {
	pagination.PageResult
}

type pageResponse struct {
	Items     []any  `json:"items"`
	NextToken string `json:"nextToken"`
}

func (current TokenPageBase) IsEmpty() (bool, error) {
	var resp pageResponse
	if err := current.ExtractInto(&resp); err != nil {
		return false, err
	}
	return len(resp.Items) == 0, nil
}

func (current TokenPageBase) NextPageURL() (string, error) {
	var resp pageResponse
	if err := current.ExtractInto(&resp); err != nil {
		return "", err
	}

	if resp.NextToken == "" {
		return "", nil
	}

	curURL := current.URL
	q := curURL.Query()

	q.Set("nextToken", resp.NextToken)
	curURL.RawQuery = q.Encode()

	return curURL.String(), nil
}

func (current TokenPageBase) GetBody() any {
	return current.Body
}
