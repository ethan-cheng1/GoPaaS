package from

import (
	"errors"
	"git.imooc.com/coding-535/common"
	"git.imooc.com/coding-535/podApi/proto/podApi"
	"strings"

	"reflect"
	"strconv"
	"time"
)

// Map data to struct based on name tags and convert types
func FromToPodStruct(data map[string]*podApi.Pair, obj interface{}) {
	objValue := reflect.ValueOf(obj).Elem()
	for i := 0; i < objValue.NumField(); i++ {
		// Get corresponding SQL value
		dataTag := strings.Replace(objValue.Type().Field(i).Tag.Get("json"),",omitempty","",-1)
		dataSlice, ok := data[dataTag]
		if !ok {
			continue
		}
		valueSlice := dataSlice.Values
		if len(valueSlice)<=0{
			continue
		}
		// Exclude port and env
		if dataTag == "pod_port" ||dataTag=="pod_env"{
			continue
		}
		value:=valueSlice[0]
		// Separate handling for ports and environment variables
		// Get corresponding field name
		name := objValue.Type().Field(i).Name
		// Get corresponding field type
		structFieldType := objValue.Field(i).Type()
		// Get variable type, can also directly write "string type"
		val := reflect.ValueOf(value)
		var err error
		if structFieldType != val.Type() {
			// Type conversion
			val, err = TypeConversion(value, structFieldType.Name()) // Type conversion
			if err != nil {
				common.Error(err)
			}
		}
		// Set type value
		objValue.FieldByName(name).Set(val)
	}
}

//类型转换
func TypeConversion(value string, ntype string) (reflect.Value, error) {
	if ntype == "string" {
		return reflect.ValueOf(value), nil
	} else if ntype == "time.Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "int" {
		i, err := strconv.Atoi(value)
		return reflect.ValueOf(i), err
	} else if ntype == "int32" {
		i, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return reflect.ValueOf(int32(i)), err
		}
		return reflect.ValueOf(int32(i)), err
	} else if ntype == "int64" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(i), err
	} else if ntype == "float32" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(float32(i)), err
	} else if ntype == "float64" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(i), err
	}

	//else if ....... add other type conversions

	return reflect.ValueOf(value), errors.New("Unknown type: " + ntype)
}

