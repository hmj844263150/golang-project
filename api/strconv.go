package api

import (
	"espressif.com/chip/factory/dal"
	"strconv"
	"time"
)

var ptypeMap = map[string]uint8{"float64": 1, "float32": 2, "int": 3, "int64": 4, "int32": 5, "int16": 6, "int8": 7, "uint64": 8, "uint32": 9, "uint16": 10, "uint8": 11, "time.Time": 12, "string": 13, "*": 14}
var ptypeFunc = map[uint8]func(v string) (interface{}, error){1: toFloat64, 2: toFloat32, 3: toInt, 4: toInt64, 5: toInt32, 6: toInt16, 7: toInt8, 8: toUint64, 9: toUint32, 10: toUint16, 11: toUint8, 12: toTime, 13: toString, 14: toString}

func toFloat64(v string) (interface{}, error) {
	value, err := strconv.ParseFloat(v, 32)
	if err != nil {
		return nil, err
	}
	return value, nil
}
func toFloat32(v string) (interface{}, error) {
	value, err := strconv.ParseFloat(v, 32)
	if err != nil {
		return nil, err
	}
	return float32(value), nil
}
func toInt(v string) (interface{}, error) {
	value, err := strconv.ParseInt(v, 10, 0)
	if err != nil {
		return nil, err
	}
	return int(value), nil
}
func toInt64(v string) (interface{}, error) {
	value, err := strconv.ParseInt(v, 10, 0)
	if err != nil {
		return nil, err
	}
	return value, nil
}
func toInt32(v string) (interface{}, error) {
	value, err := strconv.ParseInt(v, 10, 0)
	if err != nil {
		return nil, err
	}
	return int32(value), nil
}
func toInt16(v string) (interface{}, error) {
	value, err := strconv.ParseInt(v, 10, 0)
	if err != nil {
		return nil, err
	}
	return int16(value), nil
}
func toInt8(v string) (interface{}, error) {
	value, err := strconv.ParseInt(v, 10, 0)
	if err != nil {
		return nil, err
	}
	return int8(value), nil
}
func toUint64(v string) (interface{}, error) {
	value, err := strconv.ParseInt(v, 10, 0)
	if err != nil {
		return nil, err
	}
	return uint64(value), nil
}
func toUint32(v string) (interface{}, error) {
	value, err := strconv.ParseInt(v, 10, 0)
	if err != nil {
		return nil, err
	}
	return uint32(value), nil
}
func toUint16(v string) (interface{}, error) {
	value, err := strconv.ParseInt(v, 10, 0)
	if err != nil {
		return nil, err
	}
	return uint16(value), nil
}
func toUint8(v string) (interface{}, error) {
	value, err := strconv.ParseInt(v, 10, 0)
	if err != nil {
		return nil, err
	}
	return uint8(value), nil
}

func toTime(v string) (interface{}, error) {
	value, err := dal.Time(v, time.UTC)
	if err != nil {
		return nil, err
	}
	return value, nil
}
func toString(v string) (interface{}, error) {
	return v, nil
}
