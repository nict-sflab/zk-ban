package encode

import (
	"bytes"
	"io"
)

type WithWriteTo interface {
	WriteTo(io.Writer) (int64, error)
}

type WithWriteDump interface {
	WriteDump(io.Writer) error
}

type WithReadFrom interface {
	ReadFrom(io.Reader) (int64, error)
}

func EncodeWithWriteTo[T WithWriteTo](proof T) ([]byte, error) {
	var buffer bytes.Buffer
	_, err := proof.WriteTo(&buffer)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func EncodeWithWriteDump[T WithWriteDump](proof T) ([]byte, error) {
	var buffer bytes.Buffer
	err := proof.WriteDump(&buffer)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func DecodeWithReadFrom[T WithReadFrom](buf []byte, t T) error {
	csReader := bytes.NewReader(buf)
	_, err := t.ReadFrom(csReader)
	return err
}
