package json_fetchers

import (
	"log"
	"testing"
)

func TestFileJsonFetcher_New(t *testing.T) {
	f, _ := CreateFileFetcher("/Users/ankurrana-mbp/Documents/codebase/titli/result.json", 5)
	for i := 0; i < 5; i++ {
		data, _ := f.GetNextJsonQueries()
		for j := 0; j < len(data); j++ {
			log.Println(data[j]["from"])
		}
	}
}
