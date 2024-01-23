package mapping_definer

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/ankur-toko/es-mapping-analyser/optimization_engine"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func RecommendMappingDeprecated(rawMapping json.RawMessage, optimisationSet optimization_engine.OptimizationSet) string {
	rawJsonStrT, _ := json.Marshal(&rawMapping)
	rawJsonStr := string(rawJsonStrT)
	// Read the raw mapping and update the mapping with the optimisation from the optimisation set
	var err error
	for field, optimisation := range optimisationSet {
		for _, code := range optimisation.Code {
			if code == "-DOCVALUES" {
				rawJsonStr, err = sjson.Set(rawJsonStr, fmt.Sprintf("mappings.properties.%v.doc_values", field), "false")
				if err != nil {
					log.Print("Error setting doc values false")
				}
			}

			if code == "+OBJECT" {
				rawJsonStr, err = sjson.Set(rawJsonStr, fmt.Sprintf("mappings.properties.%v.type", field), "object")
				if err != nil {
					log.Printf("Error setting type for this field %v \n", field)
				}

				rawJsonStr, err = sjson.Set(rawJsonStr, fmt.Sprintf("mappings.properties.%v.enabled", field), "false")
				if err != nil {
					log.Print("Error setting enabled false")
				}
			}

			if code == "+KEYWORD" {
				rawJsonStr, err = sjson.Set(rawJsonStr, fmt.Sprintf("mappings.properties.%v.type", field), "keyword")
				if err != nil {
					log.Printf("Error setting type for this field %v \n", field)
				}
			}

			if code == "-INDEX" {
				rawJsonStr, err = sjson.Set(rawJsonStr, fmt.Sprintf("mappings.properties.%v.index", field), "false")
				if err != nil {
					log.Printf("Error setting type for this field %v \n", field)
				}
			}

		}

	}
	log.Print(rawMapping)
	return rawJsonStr
}

func RecommendMapping(rawMapping json.RawMessage, optimisationSet optimization_engine.OptimizationSet) string {
	rawJsonStrT, _ := json.Marshal(&rawMapping)
	rawJsonStr := string(rawJsonStrT)
	// Read the raw mapping and update the mapping with the optimisation from the optimisation set
	var err error
	for field, optimisation := range optimisationSet {
		for _, code := range optimisation.Code {
			fullpath := GetFullPath(rawJsonStr, field)
			if code == "+OBJECT" {
				rawJsonStr, err = sjson.Set(rawJsonStr, fmt.Sprintf("%v", fullpath), map[string]interface{}{"type": "object", "enabled": "false"})
				if err != nil {
					log.Printf("Error setting type for this field %v \n", field)
				}

				// rawJsonStr, err = sjson.Set(rawJsonStr, fmt.Sprintf("%v.type", fullpath), "object")
				// if err != nil {
				// 	log.Printf("Error setting type for this field %v \n", field)
				// }

				// rawJsonStr, err = sjson.Set(rawJsonStr, fmt.Sprintf("%v.enabled", fullpath), "false")
				// if err != nil {
				// 	log.Print("Error setting enabled false")
				// }
			} else {

				if code == "-DOCVALUES" {
					rawJsonStr, err = sjson.Set(rawJsonStr, fmt.Sprintf("%v.doc_values", fullpath), "false")
					if err != nil {
						log.Print("Error setting doc values false")
					}
				}
				if code == "+KEYWORD" {
					rawJsonStr, err = sjson.Set(rawJsonStr, fmt.Sprintf("%v.type", fullpath), "keyword")
					if err != nil {
						log.Printf("Error setting type for this field %v \n", field)
					}
				}

				if code == "-INDEX" {
					rawJsonStr, err = sjson.Set(rawJsonStr, fmt.Sprintf("%v.index", fullpath), "false")
					if err != nil {
						log.Printf("Error setting type for this field %v \n", field)
					}
				}
			}

		}
	}
	log.Print(string(rawJsonStr))
	return rawJsonStr
}

func GetFullPath(rawJsonStr string, field string) string {
	arrSplit := strings.Split(field, ".")
	initialPath := fmt.Sprintf("%v.%v", "mappings.properties", arrSplit[0])
	finalPath := initialPath

	previous_path := initialPath
	for i := 1; i < len(arrSplit); i++ {
		next_path1 := fmt.Sprintf("%v.%v.%v", previous_path, "properties", arrSplit[i])
		next_path2 := fmt.Sprintf("%v.%v.%v", previous_path, "fields", arrSplit[i])
		result1 := gjson.Get(rawJsonStr, next_path1)
		result2 := gjson.Get(rawJsonStr, next_path2)
		if result1.Exists() {
			finalPath = next_path1
		}
		if result2.Exists() {
			finalPath = next_path2
		}

		previous_path = finalPath
	}

	return finalPath
}
