package utils

import (
	"errors"
	"math"
	"testing"
)

func TestConvertToUint[T int | uint8 | uint16 | uint32 | uint64](t *testing.T) {
	tests := []struct {
		name    string
		val     interface{}
		want    interface{}
		wantErr error
	}{
		{
			name:    "Convert int to uint8",
			val:     42,
			want:    uint8(42),
			wantErr: nil,
		},
		{
			name:    "Convert int to uint16",
			val:     260,
			want:    uint16(math.MaxUint8 + 1),
			wantErr: nil,
		},
		{
			name:    "Convert int to uint32",
			val:     math.MaxUint16 + 1,
			want:    uint32(math.MaxUint16 + 1),
			wantErr: nil,
		},
		{
			name:    "Convert int to uint64",
			val:     math.MaxUint32 + 1,
			want:    uint64(math.MaxUint32 + 1),
			wantErr: nil,
		},
		{
			name:    "Convert uint8 to uint8",
			val:     uint8(123),
			want:    uint8(123),
			wantErr: nil,
		},
		{
			name:    "Convert uint16 to uint16",
			val:     uint16(12345),
			want:    uint16(12345),
			wantErr: nil,
		},
		{
			name:    "Convert uint32 to uint32",
			val:     uint32(123456789),
			want:    uint32(123456789),
			wantErr: nil,
		},
		{
			name:    "Convert uint64 to uint64",
			val:     uint64(123456789123456789),
			want:    uint64(123456789123456789),
			wantErr: nil,
		},
		{
			name:    "Out of range for uint8 to uint64",
			val:     1<<64 - 1,
			want:    nil,
			wantErr: errors.New("value out of range for uint8 to uint64"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertToUint(tt.val.(T))

			if err != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("ConvertToUint() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("ConvertToUint() = %v, want %v", got, tt.want)
			}
		})
	}
}
