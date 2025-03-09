package highlevel

import (
	"encoding/json"
	"errors"
)

type Result struct {
	Out string
	Err string
}

func NewResult(out string, err error) *Result {
	errString := ""
	if err != nil {
		errString = err.Error()
	}

	return &Result{
		Out: out,
		Err: errString,
	}
}

func (result *Result) Bytes() []byte {
	bytes, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}

	return bytes
}

func (result *Result) FromBytes(bytes []byte) error {
	err := json.Unmarshal(bytes, result)
	if err != nil {
		return errors.New(result.Err)
	}

	return err
}
