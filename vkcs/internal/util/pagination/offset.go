package pagination

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// OffsetPageBase a page in a collection that's paginated by "offset" and "limit" query parameters.
type OffsetPageBase struct {
	pagination.PageResult

	// Body key under which the page's results are nested.
	Label string
}

// NextPageURL generates the URL for the page of results after this one.
func (current OffsetPageBase) NextPageURL() (string, error) {
	items, err := current.extractItems()
	if err != nil {
		return "", err
	}

	if len(items) == 0 {
		return "", nil
	}

	curURL := current.URL
	q := curURL.Query()

	var curOffset int
	if o := q.Get("offset"); o != "" {
		var err error
		curOffset, err = strconv.Atoi(o)
		if err != nil {
			return "", fmt.Errorf("error parsing offset: %w", err)
		}
	}

	offset := curOffset + len(items)
	q.Set("offset", strconv.Itoa(offset))
	curURL.RawQuery = q.Encode()

	return curURL.String(), nil
}

// IsEmpty satisifies the IsEmpty method of the Page interface.
func (current OffsetPageBase) IsEmpty() (bool, error) {
	items, err := current.extractItems()
	if err != nil {
		return false, err
	}

	return len(items) == 0, nil
}

// GetBody returns the linked page's body. This method is needed to satisfy the
// Page interface.
func (current OffsetPageBase) GetBody() any {
	return current.Body
}

func (current OffsetPageBase) extractItems() ([]any, error) {
	var items []any

	if current.Label == "" {
		var ok bool
		items, ok = current.Body.([]any)
		if !ok {
			return nil, newErrUnexpectedType("[]interface{}", current.Body)
		}
	} else {
		subMap, ok := current.Body.(map[string]any)
		if !ok {
			return nil, newErrUnexpectedType("map[string]interface{}", current.Body)
		}

		value, ok := subMap[current.Label]
		if !ok {
			return nil, nil
		}

		items, ok = value.([]any)
		if !ok {
			return nil, newErrUnexpectedType("[]interface{}", subMap[current.Label])
		}
	}

	return items, nil
}

func newErrUnexpectedType(expType string, body any) error {
	err := gophercloud.ErrUnexpectedType{}
	err.Expected = expType
	err.Actual = fmt.Sprintf("%v", reflect.TypeOf(body))
	return err
}
