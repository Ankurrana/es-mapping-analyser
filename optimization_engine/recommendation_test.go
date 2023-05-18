package optimization_engine

import (
	"testing"

	"github.com/ankur-toko/es-mapping-analyser/mapper"
	query_analyser "github.com/ankur-toko/es-mapping-analyser/query_analyser"
)

func TestFindOptimizations(t *testing.T) {

	var usageMap *query_analyser.UsageMap
	var properties mapper.Properties

	usageMap = &query_analyser.UsageMap{}
	usageMap.FrequencyMap = map[string]map[string]int{}

	k := map[string]int{}
	k["match"] = 10
	usageMap.FrequencyMap["MyMatchField"] = k

	FindOptimizations(usageMap, properties)

}
