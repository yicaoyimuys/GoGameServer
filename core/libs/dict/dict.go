package dict

import (
	"reflect"
)

func getValue(data interface{}, key interface{}) interface{} {
	dataType := reflect.TypeOf(data).Kind()
	if dataType != reflect.Map {
		return nil
	}

	dataValue := reflect.ValueOf(data)
	keyValue := reflect.ValueOf(key)
	value := dataValue.MapIndex(keyValue)
	if value.IsValid() {
		return value.Interface()
	}
	return nil
}

func GetBool(data interface{}, key interface{}) bool {
	result := false
	value := getValue(data, key)
	if value != nil && reflect.TypeOf(value).Kind() == reflect.Bool {
		result = value.(bool)
	}
	return result
}

func GetString(data interface{}, key interface{}) string {
	result := ""
	value := getValue(data, key)
	if value != nil && reflect.TypeOf(value).Kind() == reflect.String {
		result = value.(string)
	}
	return result
}

func GetStringMap(data interface{}, key interface{}) map[string]interface{} {
	var result map[string]interface{}
	value := getValue(data, key)
	if value != nil && reflect.TypeOf(value).Kind() == reflect.TypeOf(result).Kind() {
		result = value.(map[string]interface{})
	}
	return result
}

func GetUint16(data interface{}, key interface{}) uint16 {
	var result uint16 = 0
	value := getValue(data, key)
	if value != nil {
		valueType := reflect.TypeOf(value).Kind()
		if valueType == reflect.Float64 {
			result = uint16(value.(float64))
		} else if valueType == reflect.Uint16 {
			result = value.(uint16)
		}
	}
	return result
}

func GetUint32(data interface{}, key interface{}) uint32 {
	var result uint32 = 0
	value := getValue(data, key)
	if value != nil {
		valueType := reflect.TypeOf(value).Kind()
		if valueType == reflect.Float64 {
			result = uint32(value.(float64))
		} else if valueType == reflect.Uint32 {
			result = value.(uint32)
		}
	}
	return result
}

func GetUint64(data interface{}, key interface{}) uint64 {
	var result uint64 = 0
	value := getValue(data, key)
	if value != nil {
		valueType := reflect.TypeOf(value).Kind()
		if valueType == reflect.Float64 {
			result = uint64(value.(float64))
		} else if valueType == reflect.Uint64 {
			result = value.(uint64)
		}
	}
	return result
}

func GetInt64(data interface{}, key interface{}) int64 {
	var result int64 = 0
	value := getValue(data, key)
	if value != nil {
		valueType := reflect.TypeOf(value).Kind()
		if valueType == reflect.Float64 {
			result = int64(value.(float64))
		} else if valueType == reflect.Int64 {
			result = value.(int64)
		}
	}
	return result
}

func GetUint8(data interface{}, key interface{}) uint8 {
	var result uint8 = 0
	value := getValue(data, key)
	if value != nil {
		valueType := reflect.TypeOf(value).Kind()
		if valueType == reflect.Float64 {
			result = uint8(value.(float64))
		} else if valueType == reflect.Uint8 {
			result = value.(uint8)
		}
	}
	return result
}

func GetInt(data interface{}, key interface{}) int {
	var result int = 0
	value := getValue(data, key)
	if value != nil {
		valueType := reflect.TypeOf(value).Kind()
		if valueType == reflect.Float64 {
			result = int(value.(float64))
		} else if valueType == reflect.Int {
			result = value.(int)
		}
	}
	return result
}
