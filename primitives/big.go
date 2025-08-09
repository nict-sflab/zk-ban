package primitives

import (
	"io"
	"math/big"

	serializable "github.com/akakou/gnark-serializable"
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
	return serializable.WriteTo(b)
}

func (b *BigInt) UnmarshalJSON(buf []byte) error {
	if b == nil {
		b = NewBigInt(0)
	}

	return serializable.ReadFrom(buf, b)
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
