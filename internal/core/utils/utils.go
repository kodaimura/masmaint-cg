package utils

import (
	"fmt"
	"reflect"
	"math/rand"
	"time"
	"strconv"
)


func AtoiSlice(sl []string) ([]int, error) {
	isl := make([]int, len(sl))
	for i, v := range sl {
		x , err := strconv.Atoi(v)
		isl[i] = x
		if err != nil {
			return []int{}, err
		}
	}
	return isl, nil
}


func ItoaSlice(sl []int) []string {
	asl := make([]string, len(sl))
	for i, v := range sl {
		asl[i] = strconv.Itoa(v)
	}

	return asl
} 


func RandomString(length int, options ...string) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	if len(options) > 0 {
		charset = options[0]
	}
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}


func GetFieldValue(obj interface{}, fieldName string) (interface{}, error) {
	val := reflect.ValueOf(obj).Elem()
	fieldVal := val.FieldByName(fieldName)

	if !fieldVal.IsValid() {
		return nil, fmt.Errorf("No such field: %s", fieldName)
	}

	return fieldVal.Interface(), nil
}


func SetFieldValue(obj interface{}, fieldName string, newValue interface{}) error {
	val := reflect.ValueOf(obj).Elem()
	fieldVal := val.FieldByName(fieldName)

	if !fieldVal.IsValid() {
		return fmt.Errorf("No such field: %s", fieldName)
	}

	if !fieldVal.CanSet() {
		return fmt.Errorf("Cannot set field: %s", fieldName)
	}

	newVal := reflect.ValueOf(newValue)
	if fieldVal.Type() != newVal.Type() {
		return fmt.Errorf("Provided value type didn't match obj field type")
	}

	fieldVal.Set(newVal)
	return nil
}


func MapFields(dst, src interface{}) error {
    srcVal := reflect.ValueOf(src)
    dstVal := reflect.ValueOf(dst).Elem()

    if srcVal.Kind() == reflect.Slice || srcVal.Kind() == reflect.Array {
        if dstVal.Kind() != reflect.Slice && dstVal.Kind() != reflect.Array {
            return fmt.Errorf("dst must be a slice or array if src is a slice or array")
        }

        newSlice := reflect.MakeSlice(dstVal.Type(), srcVal.Len(), srcVal.Len())
        for i := 0; i < srcVal.Len(); i++ {
            srcElem := srcVal.Index(i).Interface()
            dstElem := reflect.New(newSlice.Index(i).Type()).Interface()

            if err := MapFields(dstElem, srcElem); err != nil {
                return err
            }
            newSlice.Index(i).Set(reflect.ValueOf(dstElem).Elem())
        }

        dstVal.Set(newSlice)
        return nil
    }

    if srcVal.Kind() != reflect.Struct || dstVal.Kind() != reflect.Struct {
        return fmt.Errorf("src and dst must be structs or arrays/slices of structs")
    }

    for i := 0; i < srcVal.NumField(); i++ {
        srcField := srcVal.Type().Field(i)
        dstField := dstVal.FieldByName(srcField.Name)

        if dstField.IsValid() && dstField.CanSet() && dstField.Type() == srcVal.Field(i).Type() {
            dstField.Set(srcVal.Field(i))
        }
    }

    return nil
}


func IsZero(value interface{}) bool {
	if value == nil {
        return true
    }
    v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return true
		}
		v = v.Elem()
	}
	return v.IsZero()
}