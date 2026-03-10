package json

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/jinzhu/copier"
	jsoniter "github.com/json-iterator/go"
	"github.com/mitchellh/mapstructure"
)

// 定义JSON操作
var (
	//json          = jsoniter.ConfigCompatibleWithStandardLibrary
	Marshal       = json.Marshal
	Unmarshal     = json.Unmarshal
	MarshalIndent = json.MarshalIndent
	NewDecoder    = json.NewDecoder
	NewEncoder    = json.NewEncoder
)

// MarshalToString JSON编码为字符串
func MarshalToString(v interface{}) string {
	s, err := jsoniter.MarshalToString(v)
	if err != nil {
		return ""
	}
	return s
}

// CopyByJson json 转换
func CopyByJson(dst, src interface{}) error {
	if dst == nil {
		return fmt.Errorf("dst cannot be nil")
	}
	if src == nil {
		return fmt.Errorf("src cannot be nil")
	}
	b, err := Marshal(src)
	if err != nil {
		return fmt.Errorf("unable to marshal src: %s", err.Error())
	}

	if err := Unmarshal(b, dst); err != nil {
		return fmt.Errorf("unable to unmarshal into dst: %s", err.Error())
	}
	return nil
}

// MapToStruct map 转 struct
func MapToStruct(dst, src interface{}) error {
	return mapstructure.Decode(src, dst)
}

// Copy struct复制
func Copy(dst, src interface{}) error {
	return copier.Copy(dst, src)
}

// MapCopy map 转 map
func MapCopy(src map[string]interface{}) (dst map[string]interface{}) {
	dst = make(map[string]interface{})
	for k, v := range src {
		dst[k] = v
	}
	return
}

// DeepCopyMap 递归深度复制 map
func DeepCopyMap(originalMap map[string]interface{}) map[string]interface{} {
	copiedMap := make(map[string]interface{})

	for key, value := range originalMap {
		switch reflect.TypeOf(value).Kind() {
		case reflect.Map:
			// 如果值是嵌套的 map，则递归复制
			copiedMap[key] = DeepCopyMap(value.(map[string]interface{}))
		case reflect.Slice:
			// 如果值是数组（slice），则创建新的数组并进行深拷贝
			copiedMap[key] = DeepCopySlice(value)
		case reflect.Ptr:
			// 如果值是指针，则创建新的指针并进行深拷贝
			copiedMap[key] = DeepCopyPointer(value)
		case reflect.Struct:
			// 如果值是结构体，则创建新的结构体并进行深拷贝
			copiedMap[key] = DeepCopyStruct(value)
		default:
			// 否则，直接赋值给复制的 map
			copiedMap[key] = value
		}
	}

	return copiedMap
}

// DeepCopySlice 深度复制 slice
func DeepCopySlice(slice interface{}) interface{} {
	sliceValue := reflect.ValueOf(slice)
	sliceType := reflect.TypeOf(slice)

	// 创建新的数组并进行深拷贝
	newSlice := reflect.MakeSlice(sliceType, sliceValue.Len(), sliceValue.Len())
	reflect.Copy(newSlice, sliceValue)

	return newSlice.Interface()
}

// DeepCopyPointer 深度复制指针
func DeepCopyPointer(pointer interface{}) interface{} {
	pointerValue := reflect.ValueOf(pointer)
	pointerType := reflect.TypeOf(pointer)

	// 创建新的指针
	newPointer := reflect.New(pointerType.Elem())
	// 复制指针指向的值
	newPointer.Elem().Set(pointerValue.Elem())

	return newPointer.Interface()
}

// DeepCopyStruct 深度复制结构体
func DeepCopyStruct(structure interface{}) interface{} {
	structureValue := reflect.ValueOf(structure)
	structureType := reflect.TypeOf(structure)

	// 创建新的结构体
	newStructure := reflect.New(structureType).Elem()
	// 复制结构体值
	newStructure.Set(structureValue)

	return newStructure.Interface()
}

func DeepCopyGob(src, dst interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}
