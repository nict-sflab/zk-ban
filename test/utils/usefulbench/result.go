package usefulbench

import (
	"encoding/json"
	"log"
	"os"
)

func (b *UsefulBenchmaker) ResultJson() string {
	res, err := json.Marshal(b.Result)

	if err != nil {
		log.Fatal(err)
	}

	return string(res)
}

func (b *UsefulBenchmaker) SaveJson(filename string) error {
	result := b.ResultJson()

	err := os.WriteFile(filename, []byte(result), 0o644)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
