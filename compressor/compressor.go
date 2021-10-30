package compressor

import (
	"encoding/binary"

	lz4 "github.com/pierrec/lz4/v4"
)

const (
	bufSize = 2628288
)

// Compressor denotes a simplified LZ4 compressor, skipping all header and checksum
// logic and allowing use with multiple io.Writers
type Compressor struct {
	buf     []byte
	lz4Comp lz4.Compressor
}

// New instantiates a new compressor
func New() *Compressor {
	return &Compressor{
		buf:     make([]byte, bufSize+4),
		lz4Comp: lz4.Compressor{},
	}
}

// Compress compresses a block, returning the compressed size and data
func (c *Compressor) Compress(p []byte) (int, []byte, error) {

	l, err := c.lz4Comp.CompressBlock(p, c.buf[4:])
	if err != nil {
		return 0, nil, err
	}

	binary.LittleEndian.PutUint32(c.buf[:4], uint32(l))

	return l + 4, c.buf, nil
}
