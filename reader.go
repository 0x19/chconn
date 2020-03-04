package chconn

import (
	"encoding/binary"
	"io"
	"math"
)

func NewReader(input io.Reader) *Reader {
	return &Reader{
		input: input,
	}
}

type Reader struct {
	compress bool
	input    io.Reader
	scratch  [binary.MaxVarintLen64]byte
}

func (r *Reader) Bool() (bool, error) {
	v, err := r.ReadByte()
	if err != nil {
		return false, err
	}
	return v == 1, nil
}

func (r *Reader) Uvarint() (uint64, error) {
	return binary.ReadUvarint(r)
}

func (r *Reader) Int8() (int8, error) {
	v, err := r.ReadByte()
	if err != nil {
		return 0, err
	}
	return int8(v), nil
}

func (r *Reader) Int16() (int16, error) {
	v, err := r.Uint16()
	if err != nil {
		return 0, err
	}
	return int16(v), nil
}

func (r *Reader) Int32() (int32, error) {
	v, err := r.Uint32()
	if err != nil {
		return 0, err
	}
	return int32(v), nil
}

func (r *Reader) Int64() (int64, error) {
	v, err := r.Uint64()
	if err != nil {
		return 0, err
	}
	return int64(v), nil
}

func (r *Reader) Uint8() (uint8, error) {
	v, err := r.ReadByte()
	if err != nil {
		return 0, err
	}
	return v, nil
}

func (r *Reader) Uint16() (uint16, error) {
	if _, err := r.input.Read(r.scratch[:2]); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint16(r.scratch[:2]), nil
}

func (r *Reader) Uint32() (uint32, error) {
	if _, err := r.input.Read(r.scratch[:4]); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(r.scratch[:4]), nil
}

func (r *Reader) Uint64() (uint64, error) {
	if _, err := r.input.Read(r.scratch[:8]); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(r.scratch[:8]), nil
}

func (r *Reader) Float32() (float32, error) {
	v, err := r.Uint32()
	if err != nil {
		return 0, err
	}
	return math.Float32frombits(v), nil
}

func (r *Reader) Float64() (float64, error) {
	v, err := r.Uint64()
	if err != nil {
		return 0, err
	}
	return math.Float64frombits(v), nil
}

func (r *Reader) FixedString(len int) ([]byte, error) {
	buf := make([]byte, len)
	if _, err := r.Read(buf); err != nil {
		return nil, err
	}
	return buf, nil
}

func (r *Reader) String() (string, error) {
	strlen, err := r.Uvarint()
	if err != nil {
		return "", err
	}
	str, err := r.FixedString(int(strlen))
	if err != nil {
		return "", err
	}
	return string(str), nil
}

func (r *Reader) ByteArray() ([]byte, error) {
	strlen, err := r.Uvarint()
	if err != nil {
		return nil, err
	}
	return r.FixedString(int(strlen))
}

func (r *Reader) ReadByte() (byte, error) {
	if _, err := r.input.Read(r.scratch[:1]); err != nil {
		return 0x0, err
	}
	return r.scratch[0], nil
}

func (r *Reader) Read(buf []byte) (int, error) {
	return r.input.Read(buf)
}

func (r *Reader) Date() (byte, error) {
	if _, err := r.input.Read(r.scratch[:1]); err != nil {
		return 0x0, err
	}
	return r.scratch[0], nil
}
