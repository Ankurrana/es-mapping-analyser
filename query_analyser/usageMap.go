package query_analyser

import (
	"fmt"
	"sort"

	"github.com/thoas/go-funk"
)

type UsageMap struct {
	FrequencyMap map[string]map[string]int
}

func NewUsageMap() UsageMap {
	u := UsageMap{}
	u.FrequencyMap = make(map[string]map[string]int)
	return u
}

/*
	FrequencyMap {
		"shop.id" : {
			"term" : 100,
			"match" : 30
		},
		"city_id" : {
			"exists" : 300
		}
	}
*/

func (u *UsageMap) AddFieldUsage(f string, use string) {
	// Adds one more usage of the field field and usecase use
	if u.FrequencyMap[f] == nil {
		u.FrequencyMap[f] = make(map[string]int)
	}
	u.FrequencyMap[f][use]++
}

func (u *UsageMap) GetUsageOf(f string) []string {
	usage := []string{}
	// Adds one more usage of the field field and usecase use

	if u.FrequencyMap[f] == nil {
		return usage
	}

	for usecase, frequency := range u.FrequencyMap[f] {
		if frequency > 0 {
			usage = append(usage, usecase)
		}
	}
	return usage

}

func (u *UsageMap) Print() {
	// Print the most used fields first,
	// All Keys
	keys := []string{}
	funk.ForEach(u.FrequencyMap, func(key string, _ interface{}) { keys = append(keys, key) })

	sort.Slice(keys, func(i, j int) bool {
		maxA := 0
		funk.ForEach(u.FrequencyMap[keys[i]], func(_ string, item int) {
			if item > maxA {
				maxA = item
			}
		})

		maxB := 0
		funk.ForEach(u.FrequencyMap[keys[j]], func(_ string, item int) {
			if item > maxB {
				maxB = item
			}
		})

		return maxA > maxB
	})

	for _, k := range keys {
		fmt.Printf("%v \n", k)
		for usecase, frequency := range u.FrequencyMap[k] {
			fmt.Printf("\t%v: %v\n", usecase, frequency)
		}
	}

	// for f, u := range u.FrequencyMap {
	// 	fmt.Printf("%v \n", f)
	// 	for usecase, frequency := range u {
	// 		fmt.Printf("\t%v: %v\n", usecase, frequency)
	// 	}
	// }
}
