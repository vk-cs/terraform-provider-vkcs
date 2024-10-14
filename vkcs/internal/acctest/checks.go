package acctest

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestCheckResourceListAttr(name, key string, value []string) resource.TestCheckFunc {
	checkFuncs := make([]resource.TestCheckFunc, 0, len(value)+1)
	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr(name, key+".#", strconv.Itoa(len(value))))
	for i, el := range value {
		idxKey := key + "." + strconv.Itoa(i)
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr(name, idxKey, el))
	}
	return resource.ComposeAggregateTestCheckFunc(checkFuncs...)
}

func TestCheckResourceMapAttr(name, key string, value map[string]string) resource.TestCheckFunc {
	checkFuncs := make([]resource.TestCheckFunc, 0, len(value)+1)
	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr(name, key+".%", strconv.Itoa(len(value))))
	for k, v := range value {
		mKey := key + "." + k
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr(name, mKey, v))
	}
	return resource.ComposeAggregateTestCheckFunc(checkFuncs...)
}
