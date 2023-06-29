package utils

import (
	"errors"
	"strconv"
)

func ToFloat64(x any) (float64, error) {
	if x == nil {
		return 0, nil
	} 

	switch x.(type) {
	case int:
		return float64(x.(int)), nil
	case int8:
		return float64(x.(int8)), nil
	case int16:
		return float64(x.(int16)), nil
	case int32:
		return float64(x.(int32)), nil
	case int64:
		return float64(x.(int64)), nil
	case float32:
		return float64(x.(float32)), nil
	case float64:
		return float64(x.(float64)), nil
	case string:
		return strconv.ParseFloat(x.(string), 64)
	default:
		return 0, errors.New("Error: from ToFloat64")
	}
}

func ToInt64(x any) (int64, error) {
	if x == nil {
		return 0, nil
	} 

	switch x.(type) {
	case int:
		return int64(x.(int)), nil
	case int8:
		return int64(x.(int8)), nil
	case int16:
		return int64(x.(int16)), nil
	case int32:
		return int64(x.(int32)), nil
	case int64:
		return int64(x.(int64)), nil
	case float32:
		return int64(x.(float32)), nil
	case float64:
		return int64(x.(float64)), nil
	case string:
		i, err := strconv.Atoi(x.(string))
		if err != nil {
			return 0, err
		}
		return int64(i), nil
	default:
		return 0, errors.New("Error: from ToFloat64")
	}
}

func ToString(x any) string {
	if x == nil {
		return ""
	} 

	switch x.(type) {
	case int:
		return strconv.Itoa(int(x.(int)))
	case int8:
		return strconv.Itoa(int(x.(int8)))
	case int16:
		return strconv.Itoa(int(x.(int16)))
	case int32:
		return strconv.Itoa(int(x.(int32)))
	case int64:
		return strconv.Itoa(int(x.(int64)))
	case float32:
		return strconv.Itoa(int(x.(float32)))
	case float64:
		return strconv.Itoa(int(x.(float64)))
	case string:
		return string(x.(string))
	default:
		return ""
	}
}