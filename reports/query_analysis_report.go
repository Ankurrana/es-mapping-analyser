package reports

import (
	"fmt"
	"hash/fnv"
	"log"
	"math"
	"sort"

	"github.com/ankur-toko/es-mapping-analyser/optimization_engine"
	"github.com/ankur-toko/es-mapping-analyser/query_analyser"
)

// Maximum number of queries fetches per call to db
const QUERY_BATCH_SIZE int = 10

// Total Number of queries analysed
const MAX_FETCHES int = 100000000

// Total Number of queries to be retained in the Query Report..
const MAX_QUERIES_IN_REPORT = 10

type QMReport struct {
	QueriesAnalysed      []string                            `json:"queries"`
	QueriesAnalyzedCount int                                 `json:"queries_count"`
	Optimizations        optimization_engine.OptimizationSet `json:"optimization"`
	UsageMap             query_analyser.UsageMap             `json:"usage_map"`
}

type QMJSONReport struct {
	Name          string                    `json:"name"`
	QueriesCount  int                       `json:"queries_analysed"`
	Optimizations map[string][]string       `json:"optimizations"`
	UsageMap      map[string]map[string]int `json:"usage_frequency_map"`
	HashCode      string                    `json:"hash_code"`
}

func NewQueryReport() *QMReport {
	optimizationSet := optimization_engine.OptimizationSet{}
	usageMap := query_analyser.UsageMap{}
	return &QMReport{[]string{}, 0, optimizationSet, usageMap}
}

func (qmr *QMReport) AddQuery(str string) {
	qmr.QueriesAnalyzedCount++

	if qmr.QueriesAnalyzedCount < 10 {
		log.Printf("query recieved: %v\n", str)
		if qmr.QueriesAnalyzedCount == 10 {
			log.Printf("will stop logging queries now..\n")
		}
	} else {
		if qmr.QueriesAnalyzedCount%nearest10(qmr.QueriesAnalyzedCount) == 0 {
			log.Printf("Queries Analysed Count: %v\n", qmr.QueriesAnalyzedCount)
		}
	}

	if qmr.QueriesAnalyzedCount < MAX_QUERIES_IN_REPORT {
		qmr.QueriesAnalysed = append(qmr.QueriesAnalysed, str)
	}
}

func (qmr *QMReport) Initialize(usageMap query_analyser.UsageMap, os optimization_engine.OptimizationSet) {
	qmr.UsageMap = usageMap
	qmr.Optimizations = os
}

func (qm *QMReport) Print() {
	queries_count := qm.QueriesAnalyzedCount
	fmt.Printf("Queries Analysed:%v\n\n", queries_count)

	fmt.Printf("Usage Map:\n")
	qm.UsageMap.Print()

	fmt.Printf("Optimizations:\n")
	qm.Optimizations.InversePrint()
}

func (qm *QMReport) JSONReport(name string) QMJSONReport {
	j := QMJSONReport{}
	j.Name = name
	j.QueriesCount = qm.QueriesAnalyzedCount
	j.Optimizations = qm.Optimizations.InverseSet()
	j.HashCode = fmt.Sprintf("%v%v", (hash(fmt.Sprint(j.Optimizations))), getCountStringFromOptimization(j.Optimizations))
	j.UsageMap = qm.UsageMap.FrequencyMap

	return j
}

func getCountStringFromOptimization(opt map[string][]string) string {
	res := ""
	keys := make([]string, 0, len(opt))
	for key := range opt {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		res = fmt.Sprintf("%v-%v", res, len(opt[key]))
	}
	return res
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func nearest10(num int) int {
	if num < 10 {
		return 1
	}

	x := int(math.Log10(float64((num))))
	return int(math.Pow(10, float64(x)))
}
