package json_fetchers

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"
)

// deprecated as of now because it doesn't support multiple indices

/* Reads a file with json strings and outputs as a list of []interface{} */
type FileJsonFetcher struct {
	// Filepath
	Filepath string
	s        *bufio.Scanner
	f        *os.File
	batch    int
}

func CreateFileFetcher(filepath string, batch int) (*FileJsonFetcher, error) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	var f FileJsonFetcher
	f.f = file
	f.batch = batch
	// Create a new scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	f.s = scanner
	// Read each line and decode it into a map[string]interface{}
	return &f, nil
}

func (f *FileJsonFetcher) GetNextJsonQueries() ([]map[string]interface{}, error) {
	batch := f.batch
	scanner := f.s
	var result []map[string]interface{}

	i := 0
	for i < batch {
		next := scanner.Scan()
		if !next {
			// Nothing left to read
			log.Print("Nothing more to read")
		}
		var data map[string]interface{}
		err := json.Unmarshal(scanner.Bytes(), &data)
		if err != nil {
			return nil, err
		}
		result = append(result, data)
		i++
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		if err != io.EOF {
			log.Println("Error", err)
		}
	}
	return result, nil
}

func (f *FileJsonFetcher) Close() {
	f.f.Close()
}
