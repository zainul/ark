package find

import (
	"reflect"
)

func InArray(value interface{}, list interface{}) bool {
	refValue := reflect.ValueOf(value)
	if refValue.Kind() == reflect.Ptr {
		refValue = refValue.Elem()
	}

	refList := reflect.ValueOf(list)
	if refList.Kind() == reflect.Ptr {
		refList = refList.Elem()
	}

	// only support slice for list parameter
	if refList.Kind() != reflect.Slice {
		panic("util.InArray list only support slice type")
	}
	// return false if slice is empty
	if refList.Len() < 1 {
		return false
	}

	// get the first slice and check the variable type must be the same as request value
	refKind := refValue.Kind()
	firstVal := refList.Index(0)
	if refKind != firstVal.Kind() {
		panic("util.InArray must use same variable type between value and slice")
	}

	if refKind == reflect.Array || refKind == reflect.Slice || refKind == reflect.Map || refKind == reflect.Struct {
		for i, j := 0, refList.Len(); i < j; i++ {
			if eq := reflect.DeepEqual(refValue.Interface(), refList.Index(i).Interface()); eq {
				return true
			}
		}
	} else {
		for i, j := 0, refList.Len(); i < j; i++ {
			if refValue.Interface() == refList.Index(i).Interface() {
				return true
			}
		}
	}

	return false
}

func InArrayInt(value int, list []int) bool {
	for _, x := range list {
		if x == value {
			return true
		}
	}

	return false
}

func InArrayInt64(value int64, list []int64) bool {
	for _, x := range list {
		if x == value {
			return true
		}
	}

	return false
}

func InArrayString(value string, list []string) bool {
	for _, x := range list {
		if x == value {
			return true
		}
	}

	return false
}
