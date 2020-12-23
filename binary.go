package binary

import (
	"encoding/binary"
	"fmt"
	"io"
)

type Reader struct {
	r    io.Reader
	err  error
	bint []byte
}

func NewReader(r io.Reader) *Reader {
	return &Reader{r: r, bint: make([]byte, 8)}
}

func (b *Reader) Err() error        { return b.err }
func (b *Reader) Reader() io.Reader { return b.r }

func (b *Reader) read(n []byte) {
	if b.err != nil {
		for i := range n {
			n[i] = 0
		}
		return
	}

	d, err := io.ReadFull(b.r, n)
	if err == io.EOF && d == len(n) {
		return
	}
	if err != nil {
		for i := range n {
			n[i] = 0
		}
		b.err = err
	}
}

func (b *Reader) ReadUint8() uint8 {
	b.read(b.bint[:1])
	return b.bint[0]
}

func (b *Reader) ReadUint16() uint16 {
	b.read(b.bint[:2])
	return binary.LittleEndian.Uint16(b.bint)
}

func (b *Reader) ReadUint32() uint32 {
	b.read(b.bint[:4])
	return binary.LittleEndian.Uint32(b.bint)
}

func (b *Reader) ReadUint64() uint64 {
	b.read(b.bint[:8])
	return binary.LittleEndian.Uint64(b.bint)
}

func (b *Reader) ReadUint(bitSize byte) uint64 {
	switch bitSize {
	case 8:
		return uint64(b.ReadUint8())
	case 16:
		return uint64(b.ReadUint16())
	case 32:
		return uint64(b.ReadUint32())
	case 64:
		return b.ReadUint64()
	}
	panic(fmt.Sprintf("invalid bitsize %d", bitSize))
}

func (b *Reader) ReadBytes(bitSize byte) []byte {
	n := b.ReadUint(bitSize)
	d := make([]byte, n)
	b.read(d)
	return d
}

func (b *Reader) ReadString(bitSize byte) string {
	return string(b.ReadBytes(bitSize))
}

type Writer struct {
	w    io.Writer
	err  error
	bint []byte
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{w: w, bint: make([]byte, 8)}
}

func (b *Writer) Err() error        { return b.err }
func (b *Writer) Writer() io.Writer { return b.w }

func (b *Writer) write(n []byte) {
	if b.err != nil {
		return
	}
	_, err := b.w.Write(n)
	if err != nil {
		b.err = err
	}
}

func (b *Writer) WriteUint8(n uint8) {
	b.bint[0] = n
	b.write(b.bint[:1])
}

func (b *Writer) WriteUint16(n uint16) {
	binary.LittleEndian.PutUint16(b.bint, n)
	b.write(b.bint[:2])
}

func (b *Writer) WriteUint32(n uint32) {
	binary.LittleEndian.PutUint32(b.bint, n)
	b.write(b.bint[:4])
}

func (b *Writer) WriteUint64(n uint64) {
	binary.LittleEndian.PutUint64(b.bint, n)
	b.write(b.bint[:8])
}

func (b *Writer) WriteUint(n uint64, bitSize byte) {
	switch bitSize {
	case 8:
		b.WriteUint8(uint8(n))
	case 16:
		b.WriteUint16(uint16(n))
	case 32:
		b.WriteUint32(uint32(n))
	case 64:
		b.WriteUint64(uint64(n))
	default:
		panic(fmt.Sprintf("invalid bitsize %d", bitSize))
	}
}

func (b *Writer) WriteBytes(str []byte, bitSize byte) {
	l := uint64(len(str))
	var max uint64 = 1<<bitSize - 1
	if l > max {
		str = str[:max]
		l = max
	}
	b.WriteUint(l, bitSize)
	b.write(str)
}

func (b *Writer) WriteString(str string, bitSize byte) {
	b.WriteBytes([]byte(str), bitSize)
}
