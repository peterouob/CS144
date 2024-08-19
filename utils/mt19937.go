package utils

const (
	n = 312
	m = 156

	matrixA uint64 = 0xB5026F5AA96619E9

	hiMask uint64 = 0xFFFFFFFF80000000
	loMask uint64 = 0x00000007FFFFFFFF

	notSeeded = n + 1
)

type MT19937Interface interface {
	Seed(int64)
	InitBySlice([]uint64)
}
type MT19937 struct {
	state []uint64
	index int
}

var _ MT19937Interface = (*MT19937)(nil)

func New() *MT19937 {
	return &MT19937{
		state: make([]uint64, n),
		index: notSeeded,
	}
}

func (mt *MT19937) Seed(seed int64) {
	x := mt.state
	x[0] = uint64(seed)
	for i := uint64(1); i < n; i++ {
		x[i] = 6364136223846793005*(x[i-1]^(x[i-1]>>62)) + i
	}
	mt.index = n
}

func (mt *MT19937) InitBySlice(key []uint64) {
	mt.Seed(19650218)
	x := mt.state
	i := uint64(1)
	j := 0
	k := len(key)
	if n > k {
		k = n
	}
	for k > 0 {
		x[i] = x[i] ^ ((x[i-1] ^ (x[i-1] >> 62)) * 3935559000370003845) + key[j] + uint64(j)
		i++
		if i >= n {
			x[0] = x[n-1]
			i = 1
		}
		j++
		if j >= len(key) {
			j = 0
		}
		k--
	}

	for j := uint64(0); j < n-1; j++ {
		x[i] = x[i] ^ ((x[i-1] ^ (x[i-1] >> 62)) * 2862933555777941757) - i
		i++
		if i >= n {
			x[0] = x[n-1]
			i = 1
		}
	}

	x[0] = 1 << 63
}

func (mt *MT19937) Uint64() uint64 {
	x := mt.state
	if mt.index >= n {
		if mt.index == notSeeded {
			mt.Seed(5489)
		}
		for i := 0; i < n-m; i++ {
			y := (x[i] & hiMask) | (x[i+1] & loMask)
			x[i] = x[i+m] ^ (y >> 1) ^ ((y & 1) * matrixA)
		}
		for i := n - m; i < n-1; i++ {
			y := (x[i] & hiMask) | (x[i+1] & loMask)
			x[i] = x[i+(m-n)] ^ (y >> 1) ^ ((y & 1) * matrixA)
		}
		y := (x[n-1] & hiMask) | (x[0] & loMask)
		x[n-1] = x[m-1] ^ (y >> 1) ^ ((y & 1) * matrixA)
		mt.index = 0
	}
	y := x[mt.index]
	y ^= (y >> 29) & 0x5555555555555555
	y ^= (y << 17) & 0x71D67FFFEDA60000
	y ^= (y << 37) & 0xFFF7EEE000000000
	y ^= y >> 43
	mt.index++
	return y
}

func (mt *MT19937) GenrandInt64() int64 {
	x := mt.state
	if mt.index >= n {
		if mt.index == notSeeded {
			mt.Seed(5489)
		}
		for i := 0; i < n-m; i++ {
			y := (x[i] & hiMask) | (x[i+1] & loMask)
			x[i] = x[i+m] ^ (y >> 1) ^ ((y & 1) * matrixA)
		}
		for i := n - m; i < n-1; i++ {
			y := (x[i] & hiMask) | (x[i+1] & loMask)
			x[i] = x[i+(m-n)] ^ (y >> 1) ^ ((y & 1) * matrixA)
		}
		y := (x[n-1] & hiMask) | (x[0] & loMask)
		x[n-1] = x[m-1] ^ (y >> 1) ^ ((y & 1) * matrixA)
		mt.index = 0
	}
	y := x[mt.index]
	y ^= (y >> 29) & 0x5555555555555555
	y ^= (y << 17) & 0x71D67FFFEDA60000
	y ^= (y << 37) & 0xFFF7EEE000000000
	y ^= y >> 43
	mt.index++
	return int64(y & 0x7fffffffffffffff)
}

func (mt *MT19937) Read(p []byte) (n int, err error) {
	n = len(p)
	for len(p) >= 8 {
		val := mt.Uint64()
		p[0] = byte(val)
		p[1] = byte(val >> 8)
		p[2] = byte(val >> 16)
		p[3] = byte(val >> 24)
		p[4] = byte(val >> 32)
		p[5] = byte(val >> 40)
		p[6] = byte(val >> 48)
		p[7] = byte(val >> 56)
		p = p[8:]
	}
	if len(p) > 0 {
		val := mt.Uint64()
		for i := 0; i < len(p); i++ {
			p[i] = byte(val)
			val >>= 8
		}
	}
	return n, nil
}
