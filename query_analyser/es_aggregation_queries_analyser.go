package query_analyser

import (
	"fmt"
	"log"

	funk "github.com/thoas/go-funk"
)

func ReadJsonForAggsQueries(m map[string]interface{}, pre_keys *[]string, usageMap *UsageMap) {
	for k, val := range m {
		if IsAggLeafField(k) {
			context := "aggregations"
			ProcessAggLeaf(k, val, usageMap, context)
		} else if v, ok := val.(map[string]interface{}); ok {
			ReadJsonForAggsQueries(v, pre_keys, usageMap)
		}
	}
}

func AggregationFields() []string {
	return []string{"sum", "stats", "rate", "percentiles", "max", "min", "range", "exist", "match", "term", "terms", "auto_date_histogram", "date_histogram", "date_range", "avg", "histogram", "missing", "multi_terms"}
}

func AggregationFieldsWithPreix() []string {
	list := []string{}
	for _, i := range AggregationFields() {
		list = append(list, getAggKey(i))
	}
	return list
}

func IsAggLeafField(k string) bool {
	return funk.Contains(AggregationFields(), k)
}

func getAggKey(a string) string {
	return fmt.Sprintf("agg_%v", a)
}

func ProcessAggLeaf(k string, v interface{}, usageMap *UsageMap, context string) {
	if !IsAggLeafField(k) {
		log.Printf("This is not a leaf field! : %v\n", k)
		return
	}
	if k == "source" {
		fs := GetFieldsFromSource(v)
		if len(fs) > 0 {
			for _, f := range fs {
				usageMap.AddFieldUsage(f, getAggKey(k))
			}
		}
		return
	}

	fields := []string{}
	// A leaf field could be a string, []string or map[string]interface{}
	m, okMap := v.(map[string]interface{})
	arr, okArray := v.([]interface{})
	str, okString := v.(string)

	if okMap {
		fields = fetchFieldsFromMap(m)
	} else if okArray {
		for _, item := range arr {
			// Exploding array into individual items
			ProcessLeaf(k, item, usageMap, context)
		}
	} else if okString {
		// Script is already handled
		log.Printf("%v : %v", v, str)

	} else {
		log.Printf("Are you Crazy? : %v\n", v)
	}

	// Update the usage map
	for _, f := range fields {
		usageMap.AddFieldUsage(f, getAggKey(k))
	}

}
