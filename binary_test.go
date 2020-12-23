package binary_test

import (
	"bytes"
	"testing"

	"github.com/frizinak/binary"
)

func TestEncDec(t *testing.T) {
	d := bytes.NewBuffer(nil)

	w := binary.NewWriter(d)

	long := "longsser than 256 abclonger than 256 abclonger than 256 abclonger than 256 abclonger than 256 abclonger than 256 abclonger than 256 abclonger than 256 abclonger than 256 abclonger than 256 abclonger than 256 abclonger than 256 abclonger than 256 abclonger than 256 abclonger than 256 abclonger than 256 abclonger than 256 abclonger than 256 abclonger than 256 abclonger than 256 abclonger than 256 abclonger than 256 abclonger than 256 abclonger than 256 abclonger than 256 abclonger than 256 abclonger than 256 abclonger than 256 abclonger than 256 abclonger than 256 abclonger than 256 abclonger than 256 abclonger than 256 abclonger than 256 abclonger than 256 abclonger than 256 abclonger than 256 abclonger than 256 abclonger than 256 abclonger than 256 abclonger than 256 abclonger than 256 abclonger than 256 abclonger than 256 abclonger than 256 abclonger than 256 abclonger than 256 abclonger than 256 abclonger than 256 abclonger than 256 abclonger than 256 abclonger than 256 abc longer than 256 abc "

	data := []interface{}{
		"hey ", "there", "split string",
		long,
		0, 1, 2, 32, 255,
		[]byte{1, 2, 3},
		[]byte{255, 240, 60},
	}

	bitSizes := []byte{8, 16, 32, 64}
	for _, bs := range bitSizes {
		for i := range data {
			switch v := data[i].(type) {
			case string:
				w.WriteString(v, bs)
			case int:
				w.WriteUint(uint64(v), bs)
			case []byte:
				w.WriteBytes(v, bs)
			default:
				panic(v)
			}
		}
	}

	if err := w.Err(); err != nil {
		t.Error(err)
		t.FailNow()
	}

	r := binary.NewReader(d)

	for _, bs := range bitSizes {
		for i := range data {
			switch v := data[i].(type) {
			case string:
				nv := r.ReadString(bs)
				if v == long && bs == 8 {
					if nv != v[:255] {
						t.Errorf("invalid string %s != %s", nv, v[:255])
					}
					continue
				}
				if nv != v {
					t.Errorf("invalid string %s != %s", nv, v)
				}
			case int:
				nv := r.ReadUint(bs)
				if int(nv) != v {
					t.Errorf("invalid int %d != %d", nv, v)
				}
				if nv != uint64(int(nv)) {
					t.Errorf("truncated int %d != %d", uint64(int(nv)), nv)
				}
			case []byte:
				nv := r.ReadBytes(bs)
				if !bytes.Equal(nv, v) {
					t.Errorf("invalid byte slice %v != %v", nv, v)
				}
			default:
				panic(v)
			}
		}
	}

	if err := r.Err(); err != nil {
		t.Error(err)
	}
}
