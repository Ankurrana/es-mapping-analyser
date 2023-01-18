package json_fetchers

import (
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/nxadm/tail"
)

/* Reads a file with json strings and outputs as a list of []interface{} */
type TailFileFetcher struct {
	// Filepath
	Filepath string
	Tail     *tail.Tail
	batch    int
}

func CreateTailFileFetcher(filepath string, batch int) (*TailFileFetcher, error) {

	tail, err := tail.TailFile(filepath, tail.Config{Follow: true})
	if err != nil {
		log.Println(err)
	}
	var f TailFileFetcher
	f.Filepath = filepath
	f.batch = batch
	f.Tail = tail
	return &f, nil
}

func (f *TailFileFetcher) GetNextJsonQueries() (map[string][]map[string]interface{}, error) {
	batch := f.batch
	result := NewQuerySegregator()

	tail := f.Tail
	batchBreaker := 0
	i := 0
	for i < batch {
		select {
		case data, ok := <-tail.Lines:
			if tail.Err() != nil && !strings.Contains(tail.Err().Error(), "still alive") {
				log.Print(tail.Err().Error())
			}
			if !ok {
				return nil, errors.New("error reading channel")
			}
			if strings.Contains(data.Text, "_msearch") {
				result.AppendQueryMap(ReadSearchQueryFromMSearch(data.Text))
			} else {
				index, query, ok := ReadJsonFromText(data.Text)
				if ok {
					result.Append(index, query)
				}
			}
		default:
			time.Sleep(1 * time.Second)
			batchBreaker++
		}
		if batchBreaker > 5 {
			break
		}
		i++
	}
	return result.GetRawResponse(), nil
}

func (f *TailFileFetcher) Close() {
	//close(f.Tail.Lines)
}

func ReadJsonFromText(data string) (string, map[string]interface{}, bool) {
	var d map[string]interface{}
	k := strings.Split(data, "$$$")
	url := k[1]
	index := GetIndexFromURL(url)
	if len(k) > 2 {
		json.Unmarshal([]byte(k[2]), &d)
	} else {
		return "", nil, false
	}
	return index, d, true
}

func ReadSearchQueryFromMSearch(data string) map[string][]map[string]interface{} {
	// Assuming data is a set of concatenated set of json usually in the order
	//{"index":"<INDEX_NAME>"}{<SEARCH QUERY>}{"index":"<INDEX_NAME_2>"}{<SEARCH QUERY 2>}

	responses := map[string][]map[string]interface{}{}
	k := strings.Split(data, "$$$")
	data = k[2]
	ProcessMSearchRequest(data)

	return responses
}

func ProcessMSearchRequest(data string) map[string][]map[string]interface{} {
	responses := map[string][]map[string]interface{}{}
	jsonStrings := []string{}

	brackets := 0
	s := -1
	for i, c := range data {
		if c == '{' {
			if brackets == 0 {
				s = i
			}
			brackets++
		}
		if c == '}' {
			brackets--
		}

		if brackets == 0 {
			// We found a json
			json := data[s : i+1]
			jsonStrings = append(jsonStrings, json)
		}
	}

	m := map[string]interface{}{}
	indexName := ""
	for i, val := range jsonStrings {
		m = map[string]interface{}{}
		e := json.Unmarshal([]byte(val), &m)
		if e != nil {
			log.Print(e)
			return nil
		}
		if i%2 == 0 {
			indexName = m["index"].(string)
		} else {
			if responses[indexName] == nil {
				responses[indexName] = make([]map[string]interface{}, 0)
			}
			responses[indexName] = append(responses[indexName], m)
		}
	}

	return responses
}

func GetIndexFromURL(url string) string {
	// /product_ace/_search
	strs := strings.Split(url, "/")
	if len(strs) >= 2 {
		return strs[1]
	}
	return ""
}
