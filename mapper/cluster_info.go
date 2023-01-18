package mapper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ClusterInfo struct {
	Name        string `json:"name"`
	ClusterName string `json:"cluster_name"`
	ClusterUUID string `json:"cluster_uuid"`
	Version     struct {
		Number                           string    `json:"number"`
		BuildFlavor                      string    `json:"build_flavor"`
		BuildType                        string    `json:"build_type"`
		BuildHash                        string    `json:"build_hash"`
		BuildDate                        time.Time `json:"build_date"`
		BuildSnapshot                    bool      `json:"build_snapshot"`
		LuceneVersion                    string    `json:"lucene_version"`
		MinimumWireCompatibilityVersion  string    `json:"minimum_wire_compatibility_version"`
		MinimumIndexCompatibilityVersion string    `json:"minimum_index_compatibility_version"`
	} `json:"version"`
}

func GetVersion(es_url string) (int, error) {
	req := fmt.Sprintf("http://%v", es_url)

	data, err := http.Get(req)
	if err != nil {
		return -1, err
	}

	var response ClusterInfo
	a, err := ioutil.ReadAll(data.Body)
	if err != nil {
		log.Print(err)
		return -1, err
	}
	err = json.Unmarshal(a, &response)

	if err != nil {
		log.Printf("Error %v", err)
		return -1, err
	}

	version := strings.Split(response.Version.Number, ".")[0]

	v, e := strconv.ParseInt(version, 10, 32)
	if e != nil {
		log.Print(e)
		return -1, e
	}

	return int(v), nil
}
