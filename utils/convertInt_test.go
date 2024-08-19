package utils

import (
	"math"
	"testing"
)

func TestConvertIntToUint8(t *testing.T) {
	val := 42
	want := uint8(42)
	got, err := ConvertToUint(val)
	if err != nil {
		t.Errorf("ConvertToUint() error = %v, wantErr = nil", err)
	}
	if got != want {
		t.Errorf("ConvertToUint() = %v, want %v", got, want)
	}
}

func TestConvertIntToUint16(t *testing.T) {
	val := 260
	want := uint16(260)
	got, err := ConvertToUint(val)
	if err != nil {
		t.Errorf("ConvertToUint() error = %v, wantErr = nil", err)
	}
	if got != want {
		t.Errorf("ConvertToUint() = %v, want %v", got, want)
	}
}

func TestConvertIntToUint32(t *testing.T) {
	val := math.MaxUint16 + 1
	want := uint32(val)
	got, err := ConvertToUint(val)
	if err != nil {
		t.Errorf("ConvertToUint() error = %v, wantErr = nil", err)
	}
	if got != want {
		t.Errorf("ConvertToUint() = %v, want %v", got, want)
	}
}

func TestConvertIntToUint64(t *testing.T) {
	val := math.MaxUint32 + 1
	want := uint64(val)
	got, err := ConvertToUint(val)
	if err != nil {
		t.Errorf("ConvertToUint() error = %v, wantErr = nil", err)
	}
	if got != want {
		t.Errorf("ConvertToUint() = %v, want %v", got, want)
	}
}

func TestConvertUint8ToUint8(t *testing.T) {
	val := uint8(123)
	want := uint8(123)
	got, err := ConvertToUint(val)
	if err != nil {
		t.Errorf("ConvertToUint() error = %v, wantErr = nil", err)
	}
	if got != want {
		t.Errorf("ConvertToUint() = %v, want %v", got, want)
	}
}

func TestConvertUint16ToUint16(t *testing.T) {
	val := uint16(12345)
	want := uint16(12345)
	got, err := ConvertToUint(val)
	if err != nil {
		t.Errorf("ConvertToUint() error = %v, wantErr = nil", err)
	}
	if got != want {
		t.Errorf("ConvertToUint() = %v, want %v", got, want)
	}
}
