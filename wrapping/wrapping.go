package wrapping

type WrappingInt32Interface interface {
	SetRawValue(uint32) *WrappingInt32
	RawValue() uint32
	Wrap(uint64, WrappingInt32) WrappingInt32
	UnWrap(WrappingInt32, WrappingInt32, uint64) uint64
}

type WrappingInt32 struct {
	rawValue uint32
}

var _ WrappingInt32Interface = (*WrappingInt32)(nil)

func NewWrrappingInt32() *WrappingInt32 {
	return &WrappingInt32{rawValue: 0}
}

func (w *WrappingInt32) SetRawValue(v uint32) *WrappingInt32 {
	w.rawValue = v
	return w
}
func (w *WrappingInt32) RawValue() uint32 { return w.rawValue }
func (w *WrappingInt32) Wrap(n uint64, isn WrappingInt32) WrappingInt32 {
	return WrappingInt32{rawValue: isn.rawValue + uint32(n)}
}
func (w *WrappingInt32) UnWrap(n, isn WrappingInt32, checkPoint uint64) uint64 {
	INT32Range := uint64(1 << 32)
	//log.Printf("INT32Range :%d, 1 << 32 :%d", INT32Range, 1<<32)
	offset := uint64(n.rawValue - isn.rawValue)
	//log.Printf("n.rawValue =%d,isn.rawValue =%d,offset = %d", n.rawValue, isn.rawValue, offset)
	if checkPoint > offset {
		realCheckPoint := (checkPoint - offset) + (INT32Range >> 1)
		var wrapNum = realCheckPoint / INT32Range
		//	log.Printf("reaCheckPoint =%d,INT32Range =%d,wrapNum =%d", realCheckPoint, INT32Range, wrapNum)
		return wrapNum*INT32Range + offset
	}
	return offset
}
