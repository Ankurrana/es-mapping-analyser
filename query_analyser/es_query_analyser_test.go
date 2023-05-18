package query_analyser

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ankur-toko/es-mapping-analyser/utils"
)

var myJson string = `{
    "from": 24,
    "query": {
        "function_score": {
            "boost_mode": "replace",
            "functions": [
                {
                    "script_score": {
                        "script": {
                            "source": "doc['bestmatch_score_controller3'].value + doc['bestmatch_score_controller23'].value"
                        }
                    },
                    "weight": 1
                },
                {
                    "filter": {
                        "bool": {
                            "_name": "7",
                            "filter": [
                                {
                                    "range": {
                                        "price": {
                                            "from": 93000,
                                            "include_lower": true,
                                            "include_upper": true,
                                            "to": 130000
                                        }
                                    }
                                },{
                                    "match" : {
                                        "name" : "baju"
                                    }
                                },
                                {
                                    "terms": {
                                        "child_category_id": [
                                            2292
                                        ]
                                    }
                                }
                            ]
                        }
                    },
                    "script_score": {
                        "script": {
                            "source": "6.000000"
                        }
                    }
                },
                {
                    "filter": {
                        "terms": {
                            "_name": "21",
                            "warehouse.city_id": [
                                165
                            ]
                        }
                    },
                    "script_score": {
                        "script": {
                            "source": "5"
                        }
                    }
                },
                {
                    "filter": {
                        "terms": {
                            "_name": "18",
                            "child_category_id": [
                                2502
                            ]
                        }
                    },
                    "script_score": {
                        "script": {
                            "source": "5.0"
                        }
                    }
                },
                {
                    "filter": {
                        "terms": {
                            "_name": "18",
                            "child_category_id": [
                                2163
                            ]
                        }
                    },
                    "script_score": {
                        "script": {
                            "source": "3.0"
                        }
                    }
                },
                {
                    "filter": {
                        "terms": {
                            "_name": "4",
                            "id": [
                                1464334978,
                                1464474671,
                                1464477029,
                                1464483747
                            ]
                        }
                    },
                    "script_score": {
                        "script": {
                            "source": "20.0"
                        }
                    }
                },
                {
                    "filter": {
                        "term": {
                            "warehouse.city_id": {
                                "_name": "5",
                                "value": 165
                            }
                        }
                    },
                    "script_score": {
                        "script": {
                            "source": "5.000000"
                        }
                    }
                },
                {
                    "filter": {
                        "term": {
                            "nearby_districts": {
                                "_name": "14",
                                "value": 2168
                            }
                        }
                    },
                    "script_score": {
                        "script": {
                            "source": "5.0"
                        }
                    }
                }
            ],
            "query": {
                "bool": {
                    "filter": {
                        "bool": {
                            "filter": [
                                {
                                    "terms": {
                                        "department_id": [
                                            1129,
                                            2163,
                                            2289,
                                            2292,
                                            2502,
                                            2505,
                                            2577
                                        ]
                                    }
                                },
                                {
                                    "bool": {
                                        "minimum_should_match": "1",
                                        "should": [
                                            {
                                                "bool": {
                                                    "must_not": {
                                                        "exists": {
                                                            "field": "is_ico"
                                                        }
                                                    }
                                                }
                                            },
                                            {
                                                "term": {
                                                    "nearby_districts": 2168
                                                }
                                            }
                                        ]
                                    }
                                },
                                {
                                    "bool": {
                                        "must_not": {
                                            "exists": {
                                                "field": "shop.is_tokonow"
                                            }
                                        }
                                    }
                                }
                            ],
                            "must": {
                                "query_string": {
                                    "default_operator": "AND",
                                    "fields": [
                                        "extended_name"
                                    ],
                                    "fuzziness": "0",
                                    "query": "betadine sore throat (semprot OR semprotan OR penyemprot OR spray)"
                                }
                            },
                            "must_not": [
                                {
                                    "exists": {
                                        "field": "parent_id"
                                    }
                                },
                                {
                                    "terms": {
                                        "department_id": [
                                            "2091",
                                            "2096",
                                            "2097",
                                            "3911",
                                            "4398",
                                            "4405"
                                        ]
                                    }
                                },
                                {
                                    "terms": {
                                        "keyword_type": [
                                            2,
                                            3
                                        ]
                                    }
                                }
                            ]
                        }
                    }
                }
            },
            "score_mode": "sum"
        }
    },
    "sort2": {
        "price" : "desc"
    },
    "sort": [{
        "price" : "desc"
    },{
        "bestmatch_score_default" : "asc"
    }],
    "timeout": "2s"
}`

// match:[name]
// terms:[child_category_id warehouse.city_id id department_id keyword_type]
// range:[price]
// term:[warehouse.city_id nearby_districts]
// query_string:[extended_name]
// exists:[parent_id is_ico shop.is_tokonow]
// sort:[price bestmatch_score_default]
// source:[bestmatch_score_controller3 bestmatch_score_controller23]

var myJson2 string = `{
    "query": {
      "bool": {
        "must": [
          {
            "match": {
              "title": "Elasticsearch tutorial"
            }
          },
          {
            "nested": {
              "path": "comments",
              "query": {
                "bool": {
                  "must": [
                    {
                      "range": {
                        "comments.rating": {
                          "gte": 5
                        }
                      }
                    }
                  ]
                }
              }
            }
          }
        ],
        "filter": {
          "script": {
            "script": {
              "source": "Math.abs(params.start - params.end) <= 100",
              "params": {
                "start": {
                  "match": [{
                    "title": "Elasticsearch tutorial"
                  },{
                    "title2": "Elasticsearch tutorial"
                  }]
                },
                "end": {
                  "match": {
                    "comments.text": "tutorial"
                  }
                }
              }
            }
          }
        }
      }
    }
  }
`

// match:[title comments.text]
// range:[comments.rating]

var myJson3 string = `{
    "aggs" : {
        "interactions" : {
          "adjacency_matrix" : {
            "filters" : {
              "grpA" : { "terms" : { "accounts" : ["hillary", "sidney"] }},
              "grpB" : { "terms" : { "accounts" : ["donald", "mitt"] }},
              "grpC" : { "terms" : { "accounts" : ["vladimir", "nigel"] }}
            }
          }
        }
      }    
}`

//terms : [accounts]

var myJson4 = `
{
    "aggs": {
      "sales_over_time": {
        "auto_date_histogram": {
          "field": "date",
          "buckets": 10
        }
      },
      "range": {
        "date_range": {
          "field": "date",
          "format": "MM-yyyy",
          "ranges": [
            { "to": "now-10M/M" },  
            { "from": "now-10M/M" } 
          ]
        }
      }
    }
  }`

// agg_auto_date_histogram:[date]
var myjson5 string = `
{
    "aggs": {
      "avg_price": { "avg": { "field": "price" } },
      "t_shirts": {
        "filter": { "term": { "type": "t-shirt" } },
        "aggs": {
          "avg_price": { "avg": { "field": "price" } }
        }
      }
    }
  }`

//agg_avg:[price]
//agg_term:[type]

var myJson6 string = `
{
    "size": 0,
    "aggs" : {
      "messages" : {
        "filters" : {
          "filters" : {
            "errors" :   { "match" : { "body" : "error"   }},
            "warnings" : { "match" : { "body" : "warning" }}
          }
        }
      }
    }
  }`

// agg_match:[body]

var myJson7 string = `
{
    "aggs": {
      "price_ranges": {
        "range": {
          "field": "price",
          "ranges": [
            { "to": 100.0 },
            { "from": 100.0, "to": 200.0 },
            { "from": 200.0 }
          ]
        }
      }
    }
  }`

//agg_range:[price]
var myJson8 string = `{
  "query": {
      "bool": {
          "must": [
              {
                  "match": {
                      "title": "Elasticsearch tutorial"
                  }
              },
              {
                  "nested": {
                      "path": "comments",
                      "query": {
                          "bool": {
                              "must": [
                                  {
                                      "range": {
                                          "comments.rating": {
                                              "gte": 5
                                          }
                                      },
                                      "terms" : {
                                            "title" : 45
                                      }
                                  }
                              ]
                          }
                      }
                  }
              }
          ],
          "filter": {
              "script": {
                  "script": {
                      "source": "Math.abs(params.start - params.end) <= 100",
                      "params": {
                          "start": {
                              "match": {
                                  "title": "Elasticsearch tutorial"
                              }
                          },
                          "end": {
                              "match": {
                                  "comments.text": "tutorial"
                              }
                          }
                      }
                  }
              }
          }
      }
  }
}
`

func TestReadJson(t *testing.T) {
	m := make(map[string]interface{})

	// Use json.Unmarshal() to parse the JSON string into the map
	if err := json.Unmarshal([]byte(myJson8), &m); err != nil {
		fmt.Println(err)
		return
	}
	keys := &[]string{}

	myMap := NewUsageMap()

	ReadJson(m, keys, &myMap)

	myMap.Print()

}

func TestSome(t *testing.T) {
	type args struct {
		usecases []string
		list     []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "positive",
			args: args{[]string{"a", "b", "c"}, []string{"c"}},
			want: true,
		},
		{
			name: "negative",
			args: args{[]string{"a", "b", "c"}, []string{"d"}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utils.Some(tt.args.usecases, tt.args.list); got != tt.want {
				t.Errorf("Some() = %v, want %v", got, tt.want)
			}
		})
	}
}
