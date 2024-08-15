package utils

import (
	"errors"
	"math"
)

func ConvertToUint[T int | uint8 | uint16 | uint32 | uint64](val interface{}) (interface{}, error) {
	switch {
	case val.(T) >= 0 && val.(T) <= math.MaxUint8:
		return uint8(val.(T)), nil
	case val.(T) >= 0 && math.MaxUint8 < uint16(val.(T)) && uint16(val.(T)) <= math.MaxUint16:
		return uint16(val.(T)), nil
	case val.(T) >= 0 && math.MaxUint16 < uint32(val.(T)) && uint32(val.(T)) <= math.MaxUint32:
		return uint32(val.(T)), nil
	case val.(T) >= 0 && math.MaxUint32 < uint64(val.(T)) && uint(val.(T)) <= math.MaxUint64:
		return uint64(val.(T)), nil
	default:
		return nil, errors.New("value out of range for uint8 to uint64")
	}
}
