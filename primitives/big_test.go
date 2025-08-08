package primitives_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/akakou/zk-ban/primitives"
)

type Json struct {
	A   primitives.BigInt `json:"a"`
	B   primitives.BigInt `json:"b"`
	Int int               `json:"short"`
}

func TestEncode(t *testing.T) {
	a := Json{
		A:   *primitives.NewBigInt(1900),
		B:   *primitives.NewBigInt(3923),
		Int: 2,
	}

	buf, err := json.Marshal(a)

	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", string(buf))

	b := Json{}
	err = json.Unmarshal(buf, &b)
	if err != nil {
		panic(err)
	}

	if a.A.Cmp(&b.A.Int) != 0 {
		t.Fatalf("large %v: %v\n", a.A.String(), b.A.String())
	}

	if a.B.Cmp(&b.B.Int) != 0 {
		t.Fatalf("large %v: %v\n", a.B.String(), b.B.String())
	}

	if a.Int != b.Int {
		t.Fatalf("short %v: %v\n", a.Int, b.Int)
	}

}
