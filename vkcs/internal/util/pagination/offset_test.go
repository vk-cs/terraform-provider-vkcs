package pagination

import (
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
	"github.com/gophercloud/gophercloud/testhelper"
)

type OffsetPageResult struct {
	OffsetPageBase
}

func extractPageValues(r pagination.Page) ([]int, error) {
	var s struct {
		Values []int `json:"values"`
	}
	err := (r.(OffsetPageResult)).ExtractInto(&s)
	return s.Values, err
}

func createPager(baseURL string, client *gophercloud.ServiceClient, limit int) pagination.Pager {
	createPage := func(r pagination.PageResult) pagination.Page {
		return OffsetPageResult{OffsetPageBase{PageResult: r, Label: "values"}}
	}

	url := testhelper.Server.URL + baseURL
	if limit > 0 {
		url += "?limit=" + strconv.Itoa(limit)
	}

	return pagination.NewPager(client, url, createPage)
}

func createHandler(t *testing.T) http.HandlerFunc {
	values := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		testhelper.AssertNoErr(t, err)

		offset := 0
		if v := r.Form.Get("offset"); v != "" {
			var err error
			offset, err = strconv.Atoi(v)
			testhelper.AssertNoErr(t, err)
		}

		limit := 3
		if v := r.Form.Get("limit"); v != "" {
			var err error
			limit, err = strconv.Atoi(v)
			testhelper.AssertNoErr(t, err)
		}

		pageValues := []int{}
		if offset < len(values) {
			if offset+limit > len(values) {
				pageValues = values[offset:]
			} else {
				pageValues = values[offset : offset+limit]
			}
		}

		page := map[string][]int{"values": pageValues}
		pageJSON, err := json.Marshal(page)
		testhelper.AssertNoErr(t, err)

		w.Header().Add("Content-Type", "application/json")
		_, err = w.Write(pageJSON)
		testhelper.AssertNoErr(t, err)
	}
}

func createClient() *gophercloud.ServiceClient {
	return &gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{TokenID: "abc123"},
		Endpoint:       testhelper.Endpoint(),
	}
}

func TestEnumerateOffset(t *testing.T) {
	testhelper.SetupHTTP()
	testhelper.Mux.HandleFunc("/page", createHandler(t))
	defer testhelper.TeardownHTTP()

	client := createClient()
	limit := 2
	pager := createPager("/page", client, limit)

	callCount := 0
	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		actual, err := extractPageValues(page)
		if err != nil {
			return false, err
		}

		t.Logf("Handler invoked with %v", actual)

		var expected []int
		switch callCount {
		case 0:
			expected = []int{1, 2}
		case 1:
			expected = []int{3, 4}
		case 2:
			expected = []int{5, 6}
		case 3:
			expected = []int{7, 8}
		case 4:
			expected = []int{9}
		default:
			t.Fatalf("Unexpected call count: %d", callCount)
			return false, nil
		}

		testhelper.CheckDeepEquals(t, expected, actual)

		callCount++
		return true, nil
	})
	testhelper.AssertNoErr(t, err)
	testhelper.AssertEquals(t, 5, callCount)
}

func TestAllPagesOffset(t *testing.T) {
	testhelper.SetupHTTP()
	testhelper.Mux.HandleFunc("/page", createHandler(t))
	defer testhelper.TeardownHTTP()

	client := createClient()
	expected := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	for _, limit := range []int{0, 1, 5, 10, 20} {
		pager := createPager("/page", client, limit)
		page, err := pager.AllPages()
		testhelper.AssertNoErr(t, err)
		actual, err := extractPageValues(page)
		testhelper.AssertNoErr(t, err)
		testhelper.CheckDeepEquals(t, expected, actual)
	}
}
