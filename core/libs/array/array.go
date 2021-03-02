package array

import "reflect"

func IndexOf(array interface{}, value interface{}) int {
	arrType := reflect.TypeOf(array).Kind()
	if arrType != reflect.Slice && arrType != reflect.Array {
		return -1
	}

	arr := reflect.ValueOf(array)
	for i := 0; i < arr.Len(); i++ {
		if arr.Index(i).Interface() == value {
			return i
		}
	}
	return -1
}

func InArray(array interface{}, value interface{}) bool {
	return IndexOf(array, value) != -1
}
