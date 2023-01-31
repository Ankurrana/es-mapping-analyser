package mapper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Properties map[string]struct {
	Type        string `json:"type"`
	DocValues   *bool  `json:"doc_values"`
	IgnoreAbove int    `json:"ignore_above"`
	Properties  Properties
	Fields      Properties `json:"fields"`
	Enabled     *bool      `json:"enabled"`
	Index       *bool      `json:"index"`
}

type Mapping struct {
	Mappings struct {
		Source     interface{} `json:"_source"`
		Properties Properties  `json:"properties"`
		Routing    interface{} `json:"_routing"`
		Dynamic    string      `json:"dynamic"`
	} `json:"mappings"`
}

type Mapping6 struct {
	Mappings map[string]struct {
		Source     interface{} `json:"_source"`
		Properties Properties  `json:"properties"`
		Routing    interface{} `json:"_routing"`
		Dynamic    string      `json:"dynamic"`
	} `json:"mappings"`
}

func (properties *Properties) Print() {
	properties.RecursivePrint(0)
}

func (properties *Properties) RecursivePrint(tabs int) {
	for f, p := range *properties {
		if p.Properties != nil && len(p.Properties) > 0 {
			p.Properties.RecursivePrint(tabs + 1)
			continue
		}
		for i := 0; i < tabs; i++ {
			fmt.Printf("\t")
		}
		fmt.Printf("\t%v:(%v,%v) \n", f, p.Type, *p.DocValues)

		if p.Fields != nil && len(p.Fields) > 0 {
			p.Fields.RecursivePrint(tabs + 1)
			continue
		}
	}
}

func (mapping *Mapping) Print() {
	mapping.Mappings.Properties.Print()
}

// Given an index it reads all its fields as well as fields in inner object types flattened
// are returned
func GetAllMappings(es_url string) (map[string]Mapping, error) {
	req := fmt.Sprintf("http://%v/_mapping", es_url)

	data, err := http.Get(req)
	if err != nil {
		return nil, err
	}

	var response map[string]Mapping

	a, err := ioutil.ReadAll(data.Body)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	v, err := GetVersion(es_url)
	if err != nil {
		fmt.Print("Unable to retrieve version number for ES")
		return nil, err
	}
	if v > 6 {
		err = json.Unmarshal(a, &response)

		if err != nil {
			log.Printf("Error %v", err)
			return nil, err
		}
	} else {
		response = make(map[string]Mapping)

		var response6 map[string]Mapping6
		err = json.Unmarshal(a, &response6)
		if err != nil {
			log.Printf("Err: %v\n", err)
		}

		for index, mapping := range response6 {
			override := Mapping{}

			keys := make([]string, 0, len(mapping.Mappings))
			for k := range mapping.Mappings {
				keys = append(keys, k)
			}

			override.Mappings.Properties = mapping.Mappings[keys[0]].Properties
			override.Mappings.Source = mapping.Mappings[keys[0]].Source
			override.Mappings.Dynamic = mapping.Mappings[keys[0]].Dynamic
			response[index] = override
			log.Println(index)

		}

	}
	return ExplodeAllIndices(response), nil

}

func GetAliases(es_url string) map[string]string {
	result := make(map[string]string)
	// Returns the actual index for any aliases
	req := fmt.Sprintf("http://%v/_aliases", es_url)
	data, err := http.Get(req)
	if err != nil {
		log.Print("Error while reading aliases")
		return nil
	}
	var aliasMap map[string]map[string]map[string]string
	a, err := ioutil.ReadAll(data.Body)
	if err != nil {
		log.Print(err)
	}

	if err != nil {
		log.Print(err)
	}
	json.Unmarshal(a, &aliasMap)

	for index_name, aliases_body := range aliasMap {
		result[index_name] = index_name
		for alias, _ := range aliases_body["aliases"] {
			result[alias] = index_name
		}

	}

	return result
}

func FlattenProperties(properties Properties, prefix string) map[string][]string {
	/*
		Flattens the fields and their properties in a map[field][keyword, doc_values]
	*/
	result := make(map[string][]string)

	for field, m := range properties {
		if prefix != "" {
			field = fmt.Sprintf("%v.%v", prefix, field)
		}
		if m.Type != "" {
			if result[field] == nil {
				result[field] = []string{}
			}
			if m.DocValues == nil {
				if (m.Type == "text" || m.Type == "object") || (m.Enabled != nil && !*m.Enabled) {
					var k bool = false
					m.DocValues = &k
				} else {
					var k bool = true
					m.DocValues = &k
				}
			}
			result[field] = append(result[field], m.Type)
			if *m.DocValues {
				result[field] = append(result[field], "doc_values")
			}
		}

		if m.Properties != nil {
			// inner properties
			inner_properties := FlattenProperties(m.Properties, field)
			// merging inner result with the outer result
			for a, b := range inner_properties {
				result[a] = b
			}
		}

		if m.Fields != nil {
			inner_properties := FlattenProperties(m.Fields, field)
			// merging inner result with the outer result
			for a, b := range inner_properties {
				result[a] = b
			}
		}

	}
	return result
}

func ExplodeAllIndices(all_mappings map[string]Mapping) map[string]Mapping {
	result := map[string]Mapping{}
	for indexName, mapObject := range all_mappings {
		_mapping := mapObject
		p := ExplodeInnerFields(mapObject.Mappings.Properties)
		_mapping.Mappings.Properties = p
		result[indexName] = _mapping
	}
	return result
}

// 142, 141, 284, 245
//Input
// {
//	"addr" : {///},
// 	"shop" : {
// 		"id" : {
// 			"type":"keyword";
// 			"enabled" : false
// 		},
// 		"name" : {
// 			"type" : "object",
// 			"enabled" : false
// 		}
// 	 }
// }

//Output
// {
// 	"addr" : {///},
// 	//"shop.id" : {type:}
// 	//},
// 	"shop.name" : {///}
// }

func ExplodeInnerFields(properties Properties) Properties {
	// Flattens the inner fields
	p := Properties{}

	for field, obj := range properties {
		var inner_properties Properties
		multi_fields := false
		if obj.Properties != nil {

			inner_properties = ExplodeInnerFields(obj.Properties)
		} else if obj.Fields != nil {
			multi_fields = true
			p[field] = obj
			inner_properties = ExplodeInnerFields(obj.Fields)
		} else {
			p[field] = obj
		}

		if multi_fields || inner_properties == nil {
			// Multi Fields support, mapping need to be added for both outer as well as inner fields
			// Inner properties will be processed inside the recursive call to ExplodeInnerFields()
			if (p[field].Type == "text" || p[field].Type == "object") || (p[field].Enabled != nil && !*p[field].Enabled) {
				var k bool = false
				dv := p[field]
				dv.DocValues = &k
				p[field] = dv
			} else {
				// What if doc_values is already set to false?
				var k bool = true
				if obj.DocValues != nil && !*obj.DocValues {
					k = false
				}
				dv := p[field]
				dv.DocValues = &k
				p[field] = dv
			}
		}
		// Adding inner properties and fields to the outer properties
		for f, pros := range inner_properties {
			prefix := fmt.Sprintf("%v.%v", field, f)
			p[prefix] = pros
		}
	}

	return p
}

// a
// 	 b
// 		e
// 	 c
// 	 d
// f

// a.b.e
// a.c
// a.d
// f
