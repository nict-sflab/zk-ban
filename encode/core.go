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

type WithReadDump interface {
	ReadDump(io.Reader) error
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
	reader := bytes.NewReader(buf)
	_, err := t.ReadFrom(reader)
	return err
}

func DecodeWithReadDump[T WithReadDump](buf []byte, t T) error {
	reader := bytes.NewReader(buf)
	err := t.ReadDump(reader)
	return err

}
