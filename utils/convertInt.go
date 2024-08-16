package utils

import (
	"errors"
	"math"
)

func ConvertToUint[T int | uint8 | uint16 | uint32 | uint64](val T) (interface{}, error) {
	switch {
	case val >= 0 && val <= math.MaxUint8:
		return uint8(val), nil
	case val >= 0 && math.MaxUint8 < uint16(val) && uint16(val) <= math.MaxUint16:
		return uint16(val), nil
	case val >= 0 && math.MaxUint16 < uint32(val) && uint32(val) <= math.MaxUint32:
		return uint32(val), nil
	case val >= 0 && math.MaxUint32 < uint64(val) && uint64(val) <= math.MaxUint64:
		return uint64(val), nil
	default:
		return nil, errors.New("value out of range for uint8 to uint64")
	}
}
