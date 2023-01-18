package query_analyser

import (
	"fmt"
	"log"
	"regexp"

	funk "github.com/thoas/go-funk"
)

/*
	ES Mapping Profiler
	Uses Queries and Mapping as input to find possible mapping optimisations
*/

// ReadJson recursively traverses the ES JSON query and recursively reads all nested objects
// Determines the usage type of every field used in the query
func ReadJson(m map[string]interface{}, pre_keys *[]string, usageMap *UsageMap) {
	// m : es query json
	// prekeys can be used to determine the context of the usage
	for k, val := range m {
		*pre_keys = append(*pre_keys, k)
		if IsSpecialField(k) {
			if k == "aggs" || k == "aggregations" {
				ReadJsonForAggsQueries(m, pre_keys, usageMap)
			}
			// Do not traverse further some particular sections of the query
		} else if IsLeafField(k) {
			context := "search"
			ProcessLeaf(k, val, usageMap, context)
		} else if v, ok := val.([]interface{}); ok {
			// Expanding array of objects to map[string]maps and processing again
			_m := make(map[string]interface{})
			for i, item := range v {
				_m[fmt.Sprintf("%v", i)] = item
			}
			ReadJson(_m, pre_keys, usageMap)
		} else if v, ok := val.(map[string]interface{}); ok {
			ReadJson(v, pre_keys, usageMap)
		} else {
			// Could be other types like int, float etc
			//log.Printf("Ignoring %v", v)
		}
		*pre_keys = (*pre_keys)[:len(*pre_keys)-1]
	}
}

func IsSpecialField(field string) bool {
	ignorableFields := []string{"aggs", "aggregations"}
	return funk.Contains(ignorableFields, field)
}

func GetLeafFields() []string {
	return []string{"term", "terms", "match", "exists", "sort", "query_string", "range", "source"}
}

func IsLeafField(field string) bool {
	return funk.Contains(GetLeafFields(), field)
}

func ProcessLeaf(k string, v interface{}, usageMap *UsageMap, context string) {
	if !IsLeafField(k) {
		log.Printf("This is not a leaf field! : %v\n", k)
		return
	}
	if k == "source" {
		fs := GetFieldsFromSource(v)
		if len(fs) > 0 {
			for _, f := range fs {
				usageMap.AddFieldUsage(f, k)
			}
		}
		return
	}

	// Ignoring this for now
	// if k == "_source" {
	// 	fs := GetFieldsFromUnderscoreSource(v)
	// 	if len(fs) > 0 {
	// 		for _, f := range fs {
	// 			usageMap.AddFieldUsage(f, k)
	// 		}
	// 	}
	// 	return
	// }

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
		usageMap.AddFieldUsage(f, k)
	}

}

func fetchFieldsFromMap(m map[string]interface{}) []string {
	fields := []string{}
	if funk.Contains(m, "field") {
		fields = append(fields, fmt.Sprintf("%v", m["field"]))
	} else if funk.Contains(m, "fields") {
		v := m["fields"].([]interface{})
		for _, b := range v {
			f := b.(string)
			fields = append(fields, f)
		}
	} else {
		for key, _ := range m {
			if key[0] != '_' {
				// Keys starting from "_" are ES defined keys which we can ignore , for ex. _name
				fields = append(fields, key)
			}
		}
	}
	return fields
}

func GetFieldsFromSource(v interface{}) []string {
	fields := []string{}
	regex := regexp.MustCompile(`doc\['(\w+)'\]`)

	str := v.(string)

	matches := regex.FindAllStringSubmatch(str, -1)
	for _, match := range matches {
		fields = append(fields, match[1])
	}
	return fields
}

func GetFieldsFromUnderscoreSource(v interface{}) []string {
	m, ok := v.(map[string]interface{})
	if !ok {
		log.Print("Unable to transpose _source field to a map")
	}

	fields := []string{}

	if m["includes"] != nil {
		fs, ok := m["includes"].([]interface{})
		if !ok {
			log.Print("Unable to transpose _source.includes field to a map")
		}
		for i := 0; i < len(fs); i++ {
			f := fs[i].(string)
			fields = append(fields, f)
		}
	}

	if m["excludes"] != nil {
		fs, ok := m["excludes"].([]interface{})
		if !ok {
			log.Print("Unable to transpose _source.excludes field to a map")
		}
		for i := 0; i < len(fs); i++ {
			f := fs[i].(string)
			fields = append(fields, f)
		}
	}

	return fields
}

// Deprecated Now
func GetMappingFromUsageMap(usageMap map[string][]string) map[string][]string {
	// Doesn't look that fruitful, maybe we need to take a different approach
	/*
		Given a usageMap of the fields, return the appropriate mapping for each of the field
		Input
		{
			"shop.id" : ["terms","sort","terms_aggs"],
			"city_name" : ["terms"]
		}

		Output
		{
			"shop.id" : ["keyword","doc_values:true"],
			"city_name" : ["keyword"]
		}
	*/

	invertedUsageMap := InvertUsageMap(usageMap)

	// Derive type from usecase
	// Anything which is set to false is considered false
	// Todo.. Improve this later

	result := map[string][]string{}

	for field, usecases := range invertedUsageMap {
		var keyword bool = false
		var doc_values bool = false
		var text bool = false
		var index bool = false
		var object bool = false
		var long bool = false // We can't determine range
		var trueOrFalse bool = false

		for _, usecase := range usecases {
			if usecase == "sort" || usecase == "script" || usecase == "source" {
				doc_values = true
			}

			if usecase == "match" || usecase == "query_string" {
				text = true
				index = true
			}

			if usecase == "terms" || usecase == "term" {
				keyword = true
				index = true
			}

			if usecase == "range" {
				long = true
				doc_values = true
			}

			if usecase == "exists" {
				trueOrFalse = true
			}

			// What about other data type?
			// date, geo_point, byte, short, integer, boolean?
		}
		if !keyword && !text && !long && !doc_values {
			object = true
			index = false
		}

		if result[field] == nil {
			result[field] = []string{}
		}

		if keyword {
			result[field] = append(result[field], "keyword")
		}
		if doc_values {
			result[field] = append(result[field], "doc_values")
		}

		if text {
			result[field] = append(result[field], "text")
		}

		if long {
			result[field] = append(result[field], "long")
		}

		if keyword && text {
			// This is not yet supported
			result[field] = append(result[field], "multi_field")
		}

		if trueOrFalse {
			result[field] = append(result[field], "bool")
		}

		if index {
			result[field] = append(result[field], "index")
		}

		if object {
			result[field] = append(result[field], "object")
		}

	}

	return result

}

// Deprecated Now
func InvertUsageMap(usageMap map[string][]string) map[string][]string {
	invertedMap := map[string][]string{}
	for usecase, fields := range usageMap {
		for _, f := range fields {
			if invertedMap[f] == nil {
				invertedMap[f] = []string{}
			}
			invertedMap[f] = append(invertedMap[f], usecase)
		}
	}
	return invertedMap
}

// Maybe we can get better data from this, think better
