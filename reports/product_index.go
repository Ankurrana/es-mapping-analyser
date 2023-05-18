package reports

/*
	Read queries gathered by Dody's Query Analyser agent into VV's MemSql Database running on 149 server
*/

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"sync"
	"time"

	"github.com/ankur-toko/es-mapping-analyser/json_fetchers"
	"github.com/ankur-toko/es-mapping-analyser/mapper"
	"github.com/ankur-toko/es-mapping-analyser/optimization_engine"
	"github.com/ankur-toko/es-mapping-analyser/query_analyser"
)

var ClusterAnalysis *ClusterAnalyzer

type ClusterAnalyzer struct {
	UsageMapMap      map[string]query_analyser.UsageMap
	PropertiesMap    map[string]mapper.Properties
	OptimizationsMap map[string]optimization_engine.OptimizationSet
	Recommendations  map[string][]string
	QReportMap       map[string]*QMReport
	AliasesMap       map[string]string
	ESURL            string
	mu               sync.Mutex
}

func NewClusterAnalyzer(mappings map[string]mapper.Mapping, raw_aliases map[string]string, esUrl string) *ClusterAnalyzer {
	az := ClusterAnalyzer{}
	ClusterAnalysis = &az
	az.UsageMapMap = make(map[string]query_analyser.UsageMap)
	az.PropertiesMap = make(map[string]mapper.Properties)
	az.QReportMap = make(map[string]*QMReport)
	az.AliasesMap = raw_aliases
	az.ESURL = esUrl
	az.PopulateRecommendations(mappings)
	for index, mapping := range mappings {
		az.PropertiesMap[index] = mapping.Mappings.Properties
		az.UsageMapMap[index] = query_analyser.NewUsageMap()
		az.QReportMap[index] = NewQueryReport()
	}
	return &az
}

func (az *ClusterAnalyzer) Analyze(input map[string][]map[string]interface{}) {
	for index, queries := range input {
		requestIndex := index
		index = az.AliasesMap[index]
		usageMap := az.UsageMapMap[index]
		if index == "" {
			log.Printf("Index not found: %v\n", requestIndex)
			continue
		}
		for _, query := range queries {
			az.QReportMap[index].AddQuery(fmt.Sprintf("%v", query))
			query_analyser.ReadJson(query, &[]string{}, &usageMap)
		}
	}
}

func (az *ClusterAnalyzer) PopulateOptimizations(indices string) {

	for index := range az.PropertiesMap {
		index = az.AliasesMap[index]
		qm := az.QReportMap[index]
		usageMap := az.UsageMapMap[index]
		optimizations := optimization_engine.FindOptimizations(&usageMap, az.PropertiesMap[index])
		qm.Initialize(usageMap, optimizations)
		qm.AddRecommendations(az.Recommendations[index])
		qm.Print()
	}
}

func RefreshClusterAnalyserState() {
	allIndexesMap, err := mapper.GetAllMappings(ClusterAnalysis.ESURL)
	if err != nil {
		log.Printf("Unable to update mapping %v\n", err.Error())
	}
	all_aliases := mapper.GetAliases(ClusterAnalysis.ESURL)

	ClusterAnalysis.AliasesMap = all_aliases
	ClusterAnalysis.PopulateRecommendations(allIndexesMap)
	ClusterAnalysis.mu.Lock()
	for index, mapping := range allIndexesMap {
		ClusterAnalysis.PropertiesMap[index] = mapping.Mappings.Properties
	}
	ClusterAnalysis.mu.Unlock()
}

func (az *ClusterAnalyzer) PopulateRecommendations(mappings map[string]mapper.Mapping) {
	az.Recommendations = map[string][]string{}
	SET_DYNAMIC_MAPPING_FALSE := "set dynamic mapping false"

	for index, mappings := range mappings {
		if mappings.Mappings.Dynamic == "" || mappings.Mappings.Dynamic == "true" {
			az.Recommendations[index] = append(az.Recommendations[index], SET_DYNAMIC_MAPPING_FALSE)
		}
	}
}

func RunAnalysis(esUrl string, port int) error {
	allIndexesMap, err := mapper.GetAllMappings(esUrl)
	all_aliases := mapper.GetAliases(esUrl)
	no_data_count := 0
	az := NewClusterAnalyzer(allIndexesMap, all_aliases, esUrl)

	fetcher := GetFetcher_Product(port)
	defer fetcher.Close()

	if err != nil {
		log.Print(err)
	}
	for j := 0; j < MAX_FETCHES; j++ {
		data, _ := fetcher.GetNextJsonQueries()
		if len(data) > 0 {
			no_data_count = 0
			az.mu.Lock()
			az.Analyze(data)
			az.mu.Unlock()
		} else {
			no_data_count++
			if no_data_count%100 == 0 {
				log.Println("no queries logged")
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
	return nil
}

func GetFetcher_Product(port int) json_fetchers.Fetcher {
	return GetAPIFetcher(port)
}

func GetTailFetcher_Product() json_fetchers.Fetcher {
	//filename := "/Users/ankurrana-mbp/Documents/data/data.log"
	filename := "/Users/ankurrana-mbp/Documents/data/es_queries.log"
	dbfetcher, error := json_fetchers.CreateTailFileFetcher(filename, 2000)
	if error != nil {
		log.Printf("error while creating DB fetcher %v", error)
	}
	return dbfetcher
}

func GetAPIFetcher(port int) json_fetchers.Fetcher {
	fetcher := json_fetchers.NewAPIFetcher(port)
	return fetcher
}

func GetQUMMapReportFor(index_regex string, update_mapping bool) map[string]QMJSONReport {
	m := map[string]QMJSONReport{}
	if ClusterAnalysis == nil {
		return nil
	}
	if update_mapping {
		RefreshClusterAnalyserState()
	}
	ClusterAnalysis.PopulateOptimizations(index_regex)
	re := regexp.MustCompile(index_regex)

	for index, QMreport := range ClusterAnalysis.QReportMap {
		if len(index_regex) == 0 || re.Match([]byte(index)) {
			if QMreport.QueriesAnalyzedCount > 0 {
				m[index] = QMreport.JSONReport(index)
			}
		}
	}

	return m
}

func GetReportFor(index_regex string, update_mapping bool) string {
	m := GetQUMMapReportFor(index_regex, update_mapping)
	jsonString, err := json.MarshalIndent(m, "", "\t")
	if err != nil {
		log.Print(err)
	}
	return string(jsonString)
}
