package acctest

import (
	"reflect"
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

func TestCheckResourceAttrDeepEqual(name, key string, expected any) resource.TestCheckFunc {
	value := reflect.ValueOf(expected)
	switch value.Kind() {
	case reflect.Map:
		checkFuncs := make([]resource.TestCheckFunc, 0, value.Len()+1)
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr(name, key+".%", strconv.Itoa(value.Len())))
		iter := value.MapRange()
		for iter.Next() {
			checkFuncs = append(checkFuncs, TestCheckResourceAttrDeepEqual(name, key+"."+iter.Key().String(), iter.Value().Interface()))
		}

		return resource.ComposeAggregateTestCheckFunc(checkFuncs...)
	case reflect.Array, reflect.Slice:
		checkFuncs := make([]resource.TestCheckFunc, 0, value.Len()+1)
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr(name, key+".#", strconv.Itoa(value.Len())))
		for i := 0; i < value.Len(); i++ {
			checkFuncs = append(checkFuncs, TestCheckResourceAttrDeepEqual(name, key+"."+strconv.Itoa(i), value.Index(i).Interface()))
		}

		return resource.ComposeAggregateTestCheckFunc(checkFuncs...)
	case reflect.String:
		return resource.TestCheckResourceAttr(name, key, value.String())
	case reflect.Bool:
		return resource.TestCheckResourceAttr(name, key, strconv.FormatBool(value.Bool()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return resource.TestCheckResourceAttr(name, key, strconv.FormatInt(value.Int(), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return resource.TestCheckResourceAttr(name, key, strconv.FormatUint(value.Uint(), 10))
	case reflect.Float32, reflect.Float64:
		return resource.TestCheckResourceAttr(name, key, strconv.FormatFloat(value.Float(), 'f', -1, 64))
	default:
		panic("Unexpected kind of value " + value.Type().String())
	}
}
