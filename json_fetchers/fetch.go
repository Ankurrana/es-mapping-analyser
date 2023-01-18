package json_fetchers

type Fetcher interface {
	// A json fetcher may give me jsons from multi indices, where the key is the index name
	GetNextJsonQueries() (map[string][]map[string]interface{}, error)

	// This function is called at the end to close any open resources
	Close()
}

type QuerySegregator struct {
	queries map[string][]map[string]interface{}
}

func (qs *QuerySegregator) Append(index string, query map[string]interface{}) {
	if qs.queries[index] == nil {
		qs.queries[index] = make([]map[string]interface{}, 0)
	}
	qs.queries[index] = append(qs.queries[index], query)
}

func (qs *QuerySegregator) AppendQueryMap(queryMap map[string][]map[string]interface{}) {
	for index, queries := range queryMap {
		if qs.queries[index] == nil {
			qs.queries[index] = make([]map[string]interface{}, 0)
		}
		qs.queries[index] = append(qs.queries[index], queries...)
	}
}

func NewQuerySegregator() *QuerySegregator {
	q := QuerySegregator{}
	q.queries = make(map[string][]map[string]interface{})
	return &q
}

func (qs *QuerySegregator) GetRawResponse() map[string][]map[string]interface{} {
	return qs.queries
}
