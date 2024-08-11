package utils

type BufferViewListInterface interface {
	RemovePrefix(int)
	Size() int
	AsIOVecs() [][]byte
}

type BufferViewList struct {
	views []string
}

var _ BufferViewListInterface = (*BufferViewList)(nil)

func NewBufferViewList(bufferList *BufferList) *BufferViewList {
	var views []string
	for _, buffer := range bufferList.Buffers() {
		views = append(views, buffer.Str())
	}
	return &BufferViewList{
		views: views,
	}
}

func (bvl *BufferViewList) RemovePrefix(n int) {
	for n > 0 && len(bvl.views) > 0 {
		view := bvl.views[0]
		viewLen := len(view)
		if n < viewLen {
			bvl.views[0] = view[n:]
			n = 0
		} else {
			n -= viewLen
			bvl.views = bvl.views[1:]
		}
	}
}

func (bvl *BufferViewList) Size() int {
	totalSize := 0
	for _, view := range bvl.views {
		totalSize += len(view)
	}
	return totalSize
}

func (bvl *BufferViewList) AsIOVecs() [][]byte {
	iovecs := make([][]byte, len(bvl.views))
	for i, view := range bvl.views {
		iovecs[i] = []byte(view)
	}
	return iovecs
}
