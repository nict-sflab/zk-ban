package primitives

import (
	"encoding/base64"
	"fmt"
	"math/big"
	"strconv"
)

type BigInt struct {
	big.Int
}

func (b BigInt) MarshalJSON() ([]byte, error) {
	raw := b.Bytes()

	enc := base64.URLEncoding.EncodeToString(raw)
	res := strconv.Quote(enc)

	return []byte(res), nil
}

func (b *BigInt) UnmarshalJSON(p []byte) error {
	if string(p) == "null" {
		return nil
	}
	fmt.Printf("\n\nthis is msg \n-----\n%v\n-----\n\n", string(p))

	unq, err := strconv.Unquote(string(p))
	if err != nil {
		return err
	}

	raw, err := base64.URLEncoding.DecodeString(unq)
	if err != nil {
		return err
	}

	b.SetBytes(raw)
	return nil
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
