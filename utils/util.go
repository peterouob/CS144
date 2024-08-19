package utils

type InternetChecksumInterface interface {
	Add([]byte)
	Value() uint16
}

type InternetChecksum struct {
	sum    uint32
	parity bool
}

func (in *InternetChecksum) Add(data []byte) {
	for i := 0; i < len(data); i++ {
		val := data[i]
		if !in.parity {
			val = val << 8
		}
		in.sum += uint32(val)
		in.parity = !in.parity
	}
}

func (in *InternetChecksum) Value() uint16 {
	ret := in.sum
	for ret > 0xffff {
		ret = (ret >> 16) + (ret & 0xffff)
	}
	return ^uint16(ret)
}
