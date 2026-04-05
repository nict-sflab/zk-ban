package primitives

import (
	"bytes"
	"encoding/base64"
	"io"
	"math/big"
	"strconv"

	"github.com/consensys/gnark-crypto/ecc/bls12-377/fr"
)

type BigInt struct {
	big.Int
}

func (bigint *BigInt) IsNil() bool {
	return bigint == nil
}

func (bigint *BigInt) WriteTo(writer io.Writer) (int64, error) {
	buf := bigint.Int.Bytes()
	i, err := writer.Write(buf)
	return int64(i), err
}

func (bigint *BigInt) ReadFrom(reader io.Reader) (int64, error) {
	if bigint == nil {
		bigint = NewBigInt(0)
	}

	buf, err := io.ReadAll(reader)
	if err != nil {
		return 0, err
	}

	bigint.Int = *bigint.Int.SetBytes(buf)

	return 0, nil
}

func (b *BigInt) MarshalJSON() ([]byte, error) {
	return WriteTo(b)
}

func (b *BigInt) UnmarshalJSON(buf []byte) error {
	if b == nil {
		b = NewBigInt(0)
	}

	return ReadFrom(buf, b)
}

func NewBigInt(number int64) *BigInt {
	return &BigInt{
		*big.NewInt(number),
	}
}

func BigIntFromBytes(buf []byte) *BigInt {
	b := big.NewInt(0).SetBytes(buf)
	return &BigInt{*b}
}

func RandBigInt() *BigInt {
	var r fr.Element
	r.SetRandom()

	big := NewBigInt(0)
	res := r.BigInt(&big.Int)

	return &BigInt{*res}
}

type Writable interface {
	WriteTo(io.Writer) (int64, error)
	IsNil() bool
}

func WriteTo(data Writable) ([]byte, error) {
	if data.IsNil() {
		return []byte("null"), nil
	}
	var buffer bytes.Buffer
	_, err := data.WriteTo(&buffer)
	if err != nil {
		return nil, err
	}

	raw := buffer.Bytes()
	enc := base64.URLEncoding.EncodeToString(raw)
	res := strconv.Quote(enc)

	return []byte(res), nil
}

type Readable interface {
	ReadFrom(io.Reader) (int64, error)
}

func ReadFrom(buf []byte, data Readable) error {
	if string(buf) == "null" {
		return nil
	}

	unq, err := strconv.Unquote(string(buf))
	if err != nil {
		return err
	}

	raw, err := base64.URLEncoding.DecodeString(unq)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(raw)
	_, err = data.ReadFrom(reader)

	return err
}
