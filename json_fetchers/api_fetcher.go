package json_fetchers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

/* This exposes an API for anyone to push ES queries to */

type Request struct {
	URL  string
	Body string
}

var RequestBuffer chan Request
var ServerInitialized bool

type APIFetcher struct {
}

func NewAPIFetcher(port int) *APIFetcher {
	go StartServer(port)
	time.Sleep(200 * time.Millisecond)
	k := APIFetcher{}
	return &k
}

func (f *APIFetcher) GetNextJsonQueries() (map[string][]map[string]interface{}, error) {
	QSegregator := NewQuerySegregator()

	i := 0
	batchBreaker := 0
	for i < 2000 {
		select {
		case a := <-RequestBuffer:
			batchBreaker = 0
			if strings.Contains(a.URL, "_msearch") {
				QSegregator.AppendQueryMap(ProcessMSearchRequest(a.Body))
			} else {
				var d map[string]interface{}
				index := GetIndexFromURL(a.URL)
				json.Unmarshal([]byte(a.Body), &d)
				QSegregator.Append(index, d)
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

	return QSegregator.GetRawResponse(), nil
}

func StartServer(port int) {
	if ServerInitialized {
		log.Print("Server already listening")
		return
	}
	if RequestBuffer == nil {
		RequestBuffer = make(chan Request, 10000)
	}

	log.Printf("Initializing Server for API Fetcher on port %v\n", port)
	r := mux.NewRouter()
	r.PathPrefix("/").HandlerFunc(AcceptRequest)
	ServerInitialized = true
	err := http.ListenAndServe(fmt.Sprintf(":%v", port), r)
	if err != nil {
		log.Fatalf("WebApp failed with error %v", err.Error())
	}
}

func AcceptRequest(rw http.ResponseWriter, r *http.Request) {
	if RequestBuffer == nil {
		RequestBuffer = make(chan Request, 10000)
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
	}

	b := string(strings.Replace(string(body), "\n", "", -1))
	url := fmt.Sprintf("%v", r.URL)
	req := Request{url, string(b)}

	RequestBuffer <- req
	fmt.Fprintf(rw, "Done")
}

func (f *APIFetcher) Close() {

}
