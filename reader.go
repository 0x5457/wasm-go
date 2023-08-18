package wasm_go

import (
	"io"
)

type leb128Reader struct {
	bytes []byte
	pos   int
}

func (r *leb128Reader) eatBytes(length uint32) ([]byte, error) {
	end := r.pos + int(length)
	if end > len(r.bytes) {
		return nil, io.EOF
	}
	bs := r.bytes[r.pos : r.pos+int(length)]
	r.pos += int(length)
	return bs, nil
}

func (r *leb128Reader) eatString(length uint32) (string, error) {
	b, err := r.eatBytes(length)
	if err != nil {
		return "", err
	}
	return string(b), err
}

func (r *leb128Reader) eatU8() (uint8, error) {
	if r.pos >= len(r.bytes) {
		return 0, io.EOF
	}
	r.pos += 1
	return r.bytes[r.pos-1], nil
}

func (r *leb128Reader) eatU64() (uint64, error) {
	v, shift := uint64(0), 0
	for {
		u8, err := r.eatU8()
		if err != nil {
			return 0, err
		}
		v |= (uint64(u8) & 0x7F) << shift
		shift += 7
		if u8&0x80>>7 == 0 {
			break
		}
	}
	return v, nil
}

func (r *leb128Reader) eatI64() (int64, error) {
	v, shift := int64(0), 0
	for {
		u8, err := r.eatU8()
		if err != nil {
			return 0, err
		}
		v |= (int64(u8) & 0x7F) << shift
		shift += 7
		if u8&0x80>>7 == 0 {
			if u8&0x40>>3 != 0 {
				// negative number
				v |= ^0 << shift
			}
			break
		}
	}
	return v, nil
}

func (r *leb128Reader) eatI32() (int32, error) {
	v, err := r.eatI64()
	return int32(v), err
}

func (r *leb128Reader) eatU32() (uint32, error) {
	v, err := r.eatU64()
	return uint32(v), err
}
