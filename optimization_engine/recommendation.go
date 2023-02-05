package optimization_engine

import (
	"fmt"
	"sort"

	"github.com/ankur-toko/es-mapping-analyser/mapper"
	query_analyser "github.com/ankur-toko/es-mapping-analyser/query_analyser"
	"github.com/ankur-toko/es-mapping-analyser/utils"
	"github.com/thoas/go-funk"
)

type OptimizationSet map[string]Optimization

type Optimization struct {
	RemoveDocValues bool
	ConvertToObject bool
	Explain         []string
	Code            []string
}

func (mapOptimization *OptimizationSet) Print() {
	for field, Optimization := range *mapOptimization {
		fmt.Printf("%v\n", field)
		for _, x := range Optimization.Explain {
			fmt.Printf("\t%v", x)
		}
		fmt.Printf("\n")
	}
}

func (mapOptimization *OptimizationSet) InverseSet() map[string][]string {
	p := map[string][]string{}
	for field, Optimization := range *mapOptimization {
		for _, x := range Optimization.Code {
			if p[x] == nil {
				p[x] = []string{}
			}
			p[x] = append(p[x], field)
		}
	}

	for key := range p {
		sort.Strings(p[key])
	}
	return p
}

func (mapOptimization *OptimizationSet) InversePrint() {
	p := mapOptimization.InverseSet()
	for code, fields := range p {
		fmt.Printf("%v\n", code)
		for _, f := range fields {
			fmt.Printf("\t%v\n", f)
		}
		fmt.Printf("\n\n")
	}
}

func FindOptimizations(usageMap *query_analyser.UsageMap, properties mapper.Properties) OptimizationSet {
	/*
		On the basis of usage pattern of fields in the usageMap, find possible
		optimizations in the field mapping of ES Index. In other words,Checks
		if all usages of a field can be covered by cheaper mapping.

		Which optimizations are covered.
		1. Checks if doc_values can be removed from the mapping
		2. Checks if some fields are never searched, then recommends type:object for them.
		3. Checks if Texts can be converted to keyword if usage is only exact matches.
		4. Checks if numeric fields like long/short/byte can be converted to keyword
		   if usage does not have range or sort queries.
		   Todo: Are keywords cheaper than long?

		Also checks if we can optimize ES query by tweaking mapping?
		1. If range queries are being used for numeric keyword fields, then try converting
			it to long/short/byte?

	*/

	result := map[string]Optimization{}

	for field, props := range properties {
		p := Optimization{}
		p.Explain = []string{}

		if len(usageMap.GetUsageOf(field)) > 0 {
			usecases := usageMap.GetUsageOf(field)
			p.Code = []string{}
			// If no range, no aggregation and no scripts, then remove doc_values
			if *props.DocValues && !utils.Some(usecases, GetDocValuesUseCases()) {
				p.Code = append(p.Code, "-DOCVALUES")
				p.Explain = append(p.Explain, "No sort/range/script/aggs found for this field")
				p.Explain = append(p.Explain, "Convert to set doc_values as false")
			}

			// If no text searching found, then try converting text to keyword
			if props.Type == "text" && !utils.Some(usecases, GetTextualSearchUsecases()) {
				p.Code = append(p.Code, "-TEXT")
				p.Explain = append(p.Explain, "No usecases found for text search")
				p.Explain = append(p.Explain, "Option 1: Consider changing text type to other type")
			}

			// If type is numeric and there are no range queries and no sort queries, then we can convert to keyword
			if funk.Contains(GetIndexableNumericTypes(), props.Type) && !utils.Contains("range", usecases) && !utils.Contains("sort", usecases) {
				p.Code = append(p.Code, "+KEYWORD")
				p.Explain = append(p.Explain, "No usecases found for range or sort")
				p.Explain = append(p.Explain, "Consider changing numertic type (long/short/integer) to keyword")
			}

			// For range queries, we do need index bro, but not for sort
			// We can set index false when we do not want any aggregation nor any search
			if (props.Index == nil || *props.Index) && CanSetIndexFalse(usecases) {
				p.Code = append(p.Code, "-INDEX")
				p.Explain = append(p.Explain, "No usecases found for non indexable queries")
				p.Explain = append(p.Explain, "Consider setting index:false for this field")
			}

		} else {
			// No usage patterns found for the field. Use Cheapest Mapping
			if props.Type != "object" {
				p.Code = append(p.Code, "+OBJECT")
				p.ConvertToObject = true
				p.RemoveDocValues = true
				p.Explain = append(p.Explain, "No usage patterns found for this field")
				p.Explain = append(p.Explain, "Convert to type:object,enabled:false or remove if not needed")
			}

		}
		if len(p.Explain) > 0 {
			result[field] = p
		}
	}

	return result

}

func GetTextualSearchUsecases() []string {
	return []string{"match", "query_string"}
}

func GetDocValuesUseCases() []string {
	// Exists query also may use doc_values LuceneClass: DocValuesFieldExistsQuery
	// Range Queries may also use doc_values LuceneClass: IndexOrDocValuesQuery
	doc_value_usecases := []string{"sort", "range", "source", "script", "exists"}
	doc_value_usecases = append(doc_value_usecases, query_analyser.AggregationFieldsWithPreix()...)
	return doc_value_usecases
}

func GetIndexableNumericTypes() []string {
	return []string{"long", "integer", "short"}
}

func GetInverseIndexTrueUseCases() []string {
	// When can we set index:false?
	return []string{"sort", "source", "script"}
}

func CanSetIndexFalse(usecases []string) bool {
	for _, usecase := range usecases {
		if !funk.Contains(GetInverseIndexTrueUseCases(), usecase) {
			return false
		}
	}
	return true
}
