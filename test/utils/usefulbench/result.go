package usefulbench

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

func (b *UsefulBenchmaker) ResultJson() string {
	now := time.Now()
	b.Result["env"] = make(map[string]time.Duration)
	b.Result["env"]["end"] = time.Duration(now.UnixNano())
	b.Result["env"]["count"] = time.Duration(b.Max)

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
