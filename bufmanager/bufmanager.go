package bufmanager

import (
	"sync"
)

type Bufmanager struct {
	buf      [bufsize]byte //  缓存数组
	bufsc    []byte        // 缓存区
	writepos uint16        // 写入位置
	readpos  uint16        // 读取位置

	mu *sync.RWMutex
}

const (
	bufsize uint16 = 1 << 15
)

func NewBufManager() *Bufmanager {
	b := &Bufmanager{
		writepos: 0,
		readpos:  0,
		mu:       new(sync.RWMutex),
	}

	b.bufsc = b.buf[:]
	return b
}

// 写入数据
func (b *Bufmanager) WriteData(data []byte) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	datalen := len(data)
	if b.writepos >= b.readpos {
		if int(bufsize-b.writepos+b.readpos) < datalen {
			// panic("data to write is too large")
			return false
		}

		if int(bufsize-b.writepos) < datalen {
			predata := data[:(bufsize - b.writepos)]
			copy(b.bufsc[b.writepos:], predata[:])

			hinddata := data[(bufsize - b.writepos):]
			copy(b.bufsc[:len(hinddata)], hinddata[:])

			b.writepos = uint16(len(hinddata))
		} else {
			copy(b.bufsc[b.writepos:], data[:])
			b.writepos += uint16(datalen)
		}
	} else {
		if int(b.readpos-b.writepos) < datalen {
			// panic("data to write is too large")
			return false
		}

		copy(b.bufsc[b.writepos:(b.writepos+uint16(datalen))], data[:])
		b.writepos += uint16(datalen)
	}
	return true
}

// 读取数据
func (b *Bufmanager) ReadData(size uint16) ([]byte, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	var data []byte = make([]byte, size)

	if b.writepos >= b.readpos {
		if size > (b.writepos - b.readpos) {
			//panic("size want to read is overflow")
			return nil, false
		}

		copy(data[:size], b.bufsc[b.readpos:(b.readpos+size)])
		b.readpos += size
		return data, true
	} else {
		if size > (bufsize - b.readpos + b.writepos) {
			//panic("size want to read is overflow")
			return nil, false
		}

		if size <= (bufsize - b.readpos) {

			copy(data[:size], b.bufsc[b.readpos:(b.readpos+size)])
			b.readpos += size
			return data, true
		} else {
			copy(data[:(bufsize-b.readpos)], b.bufsc[b.readpos:])
			b.readpos = size + b.readpos - bufsize
			copy(data[(bufsize-b.readpos):size], b.bufsc[:b.readpos])
			return data, true
		}
	}
}

// 预读取数据
func (b *Bufmanager) ReadDataPrep(size uint16) ([]byte, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	var data []byte = make([]byte, size)

	if b.writepos >= b.readpos {
		if size > (b.writepos - b.readpos) {
			//panic("size want to read is overflow")
			return nil, false
		}

		copy(data[:size], b.bufsc[b.readpos:(b.readpos+size)])
		// b.readpos += size
		return data, true
	} else {
		if size > (bufsize - b.readpos + b.writepos) {
			//panic("size want to read is overflow")
			return nil, false
		}

		if size <= (bufsize - b.readpos) {

			copy(data[:size], b.bufsc[b.readpos:(b.readpos+size)])
			// b.readpos += size
			return data, true
		} else {
			copy(data[:(bufsize-b.readpos)], b.bufsc[b.readpos:])
			// b.readpos = size + b.readpos - bufsize
			copy(data[(bufsize-b.readpos):size], b.bufsc[:b.readpos])
			return data, true
		}
	}
}

// 清数据
func (b *Bufmanager) ClearData(size uint16) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.writepos >= b.readpos {
		if size > (b.writepos - b.readpos) {
			b.readpos = b.writepos
		} else {
			b.readpos += size
		}
	} else {
		if size > (bufsize - b.readpos + b.writepos) {
			b.readpos = b.writepos
		} else {
			if size <= (bufsize - b.readpos) {
				b.readpos += size
			} else {
				b.readpos = size + b.readpos - bufsize
			}
		}
	}
	return false
}

// 读取数据
func (b *Bufmanager) Read(p []byte) (int, error) {
	return 0, nil
}

// 可用数据量
func (b *Bufmanager) DataCount() uint16 {
	b.mu.Lock()
	defer b.mu.Unlock()

	var cnt uint16 = 0
	if b.writepos >= b.readpos {
		cnt = b.writepos - b.readpos
	} else {
		cnt = bufsize - b.readpos + b.writepos
	}
	return cnt
}

// for test
// 当前写入位置
func (b *Bufmanager) WritePos() uint16 {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.writepos
}

// for test
// 当前读取位置
func (b *Bufmanager) ReadPos() uint16 {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.readpos
}
