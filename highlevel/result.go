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

func (result *Result) String() string {
	bytes, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}

	return string(bytes)
}

func (result *Result) FromString(str string) error {
	err := json.Unmarshal([]byte(str), result)
	if err != nil {
		return errors.New(result.Err)
	}

	return err
}
