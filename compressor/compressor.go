package compressor

import (
	"encoding/binary"

	"github.com/pierrec/lz4"
)

const (
	bufSize = 2628288
)

// Compressor denotes a simplified LZ4 compressor, skipping all header and checksum
// logic and allowing use with multiple io.Writers
type Compressor struct {
	hashtable [1 << 16]int
	buf       []byte
}

// New instantiates a new compressor
func New() *Compressor {
	return &Compressor{
		buf: make([]byte, bufSize+4),
	}
}

func (c *Compressor) Compress(p []byte) (int, []byte, error) {

	l, err := lz4.CompressBlock(p, c.buf[4:], c.hashtable[:])
	if err != nil {
		return 0, nil, err
	}

	binary.LittleEndian.PutUint32(c.buf[:4], uint32(l))

	return l + 4, c.buf, nil
}
