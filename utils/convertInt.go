package utils

import (
	"errors"
	"math"
)

func ConvertToUint[T int | uint8 | uint16 | uint32 | uint64](val T) (interface{}, error) {
	switch {
	case val >= 0 && val <= math.MaxUint8:
		return uint8(val), nil
	case val >= 0 && math.MaxUint8 < val && val <= math.MaxUint16:
		return uint16(val), nil
	case val >= 0 && math.MaxUint16 < val && val <= math.MaxUint32:
		return uint32(val), nil
	case val >= 0 && math.MaxUint32 < val && val <= math.MaxUint64:
		return uint64(val), nil
	default:
		return nil, errors.New("value out of range for uint8 to uint64")
	}
}
