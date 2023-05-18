{
    "from": 0,
    "query": {
        "constant_score": {
            "filter": {
                "bool": {
                    "filter": [
                        {
                            "terms": {
                                "_name": "FILTER_BY_SHOP_ID",
                                "shop.id": [
                                    15134018
                                ]
                            }
                        },
                        {
                            "bool": {
                                "_name": "FILTER_EXC_PARENT_ID",
                                "must_not": {
                                    "exists": {
                                        "field": "parent_id"
                                    }
                                }
                            }
                        }
                    ]
                }
            }
        }
    },
    "size": 0,
    "sort": [
        {
            "id": {
                "order": "desc"
            }
        }
    ],
    "timeout": "2s",
    "track_total_hits": true
}
{
    "from": 16,
    "query": {
        "function_score": {
            "boost_mode": "replace",
            "functions": [
                {
                    "script_score": {
                        "script": {
                            "source": "doc['bestmatch_score_default'].value"
                        }
                    },
                    "weight": 1
                },
                {
                    "filter": {
                        "terms": {
                            "_name": "18",
                            "department_id": [
                                3092
                            ]
                        }
                    },
                    "script_score": {
                        "script": {
                            "source": "25"
                        }
                    }
                },
                {
                    "filter": {
                        "terms": {
                            "_name": "18",
                            "department_id": [
                                3091
                            ]
                        }
                    },
                    "script_score": {
                        "script": {
                            "source": "30"
                        }
                    }
                },
                {
                    "filter": {
                        "bool": {
                            "_name": "21",
                            "filter": {
                                "terms": {
                                    "warehouse.city_id": [
                                        252,
                                        248
                                    ]
                                }
                            }
                        }
                    },
                    "script_score": {
                        "script": {
                            "source": "1"
                        }
                    }
                },
                {
                    "filter": {
                        "term": {
                            "city_affinities": {
                                "_name": "29",
                                "value": 248
                            }
                        }
                    },
                    "script_score": {
                        "script": {
                            "source": "7"
                        }
                    }
                },
                {
                    "filter": {
                        "terms": {
                            "_name": "15",
                            "near_districts": [
                                3498
                            ]
                        }
                    },
                    "script_score": {
                        "script": {
                            "source": "4"
                        }
                    }
                },
                {
                    "filter": {
                        "terms": {
                            "_name": "15",
                            "nearer_districts": [
                                3498
                            ]
                        }
                    },
                    "script_score": {
                        "script": {
                            "source": "2"
                        }
                    }
                },
                {
                    "filter": {
                        "term": {
                            "warehouse.city_id": {
                                "_name": "5",
                                "value": 248
                            }
                        }
                    },
                    "script_score": {
                        "script": {
                            "source": "1"
                        }
                    }
                },
                {
                    "filter": {
                        "term": {
                            "free_cluster": {
                                "_name": "28",
                                "value": 206
                            }
                        }
                    },
                    "script_score": {
                        "script": {
                            "source": "4"
                        }
                    }
                },
                {
                    "filter": {
                        "term": {
                            "cheap_cluster": {
                                "_name": "28",
                                "value": 206
                            }
                        }
                    },
                    "script_score": {
                        "script": {
                            "source": "2"
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
                                    "bool": {
                                        "_name": "FILTER_NEARBY",
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
                                                    "nearby_districts": 3498
                                                }
                                            }
                                        ]
                                    }
                                },
                                {
                                    "bool": {
                                        "_name": "FILTER_NOW_GLOBAL_SRCH",
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
                                    "_name": "MUST_KEYWORD",
                                    "default_operator": "AND",
                                    "fields": [
                                        "name"
                                    ],
                                    "fuzziness": "0",
                                    "query": "galaxy ((jam tangan) OR arloji OR watch) 4 classic"
                                }
                            },
                            "must_not": [
                                {
                                    "exists": {
                                        "_name": "NOT_EXC_VAR",
                                        "field": "parent_id"
                                    }
                                },
                                {
                                    "terms": {
                                        "_name": "NOT_BAN_CAT",
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
                                        "_name": "NOT_BL_BY_DEVICE",
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
    "size": 8,
    "track_total_hits": true
}
{
    "query": {
        "constant_score": {
            "filter": {
                "bool": {
                    "filter": {
                        "exists": {
                            "field": "shop.is_official"
                        }
                    },
                    "must": {
                        "match": {
                            "extended_name": {
                                "operator": "AND",
                                "query": "silicone cover strap s22"
                            }
                        }
                    }
                }
            }
        }
    }
}
{
    "query": {
        "bool": {
            "filter": [
                {
                    "bool": {
                        "must_not": {
                            "exists": {
                                "field": "parent_id"
                            }
                        }
                    }
                },
                {
                    "query_string": {
                        "default_operator": "and",
                        "escape": false,
                        "fields": [
                            "name"
                        ],
                        "fuzziness": "0",
                        "query": "rk (cushion OR cushioned OR cushions) premium"
                    }
                }
            ]
        }
    }
}
{
    "from": 0,
    "query": {
        "constant_score": {
            "filter": {
                "bool": {
                    "must": [
                        {
                            "terms": {
                                "shop_id": [
                                    "9195008"
                                ]
                            }
                        },
                        {
                            "match": {
                                "name_ngram": {
                                    "analyzer": "standard",
                                    "operator": "AND",
                                    "query": "milk"
                                }
                            }
                        }
                    ],
                    "must_not": [
                        {
                            "exists": {
                                "field": "parent_id"
                            }
                        },
                        {
                            "exists": {
                                "field": "shop.is_tokonow"
                            }
                        }
                    ]
                }
            }
        }
    },
    "size": 3,
    "sort": [
        {
            "bestmatch_score_experiment3": {
                "order": "desc"
            }
        }
    ],
    "timeout": "300ms"
}
{
    "from": 0,
    "query": {
        "function_score": {
            "boost_mode": "replace",
            "functions": [
                {
                    "script_score": {
                        "script": {
                            "source": "doc['bestmatch_score_default'].value"
                        }
                    },
                    "weight": 1
                },
                {
                    "filter": {
                        "terms": {
                            "_name": "17",
                            "department_id": [
                                1855
                            ]
                        }
                    },
                    "script_score": {
                        "script": {
                            "source": "10"
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
                                    "bool": {
                                        "_name": "FILTER_NOW_GLOBAL_SRCH",
                                        "must_not": {
                                            "exists": {
                                                "field": "shop.is_tokonow"
                                            }
                                        }
                                    }
                                },
                                {
                                    "range": {
                                        "bestmatch_score_default": {
                                            "from": 6,
                                            "include_lower": false,
                                            "include_upper": true,
                                            "to": 110
                                        }
                                    }
                                }
                            ],
                            "must": {
                                "query_string": {
                                    "_name": "MUST_KEYWORD",
                                    "default_operator": "AND",
                                    "fields": [
                                        "name"
                                    ],
                                    "fuzziness": "0",
                                    "query": "balenciaga (bag OR tas)"
                                }
                            },
                            "must_not": [
                                {
                                    "exists": {
                                        "_name": "NOT_EXC_VAR",
                                        "field": "parent_id"
                                    }
                                },
                                {
                                    "exists": {
                                        "_name": "NOT_ADULT",
                                        "field": "is_adult"
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
    "size": 60,
    "track_total_hits": false
}
{
    "from": 60,
    "query": {
        "constant_score": {
            "filter": {
                "bool": {
                    "filter": [
                        {
                            "terms": {
                                "_name": "FILTER_BY_SHOP_ID",
                                "shop.id": [
                                    5487621
                                ]
                            }
                        },
                        {
                            "bool": {
                                "_name": "FILTER_EXC_PARENT_ID",
                                "must_not": {
                                    "exists": {
                                        "field": "parent_id"
                                    }
                                }
                            }
                        },
                        {
                            "bool": {
                                "_name": "FILTER_NEARBY_SHOPPAGE",
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
                                            "nearby_districts": 2167
                                        }
                                    }
                                ]
                            }
                        }
                    ]
                }
            }
        }
    },
    "size": 12,
    "sort": [
        {
            "id": {
                "order": "desc"
            }
        }
    ],
    "timeout": "2s",
    "track_total_hits": true
}
{
    "from": 0,
    "query": {
        "constant_score": {
            "filter": {
                "bool": {
                    "filter": {
                        "bool": {
                            "_name": "FILTER_NOW_GLOBAL_SRCH",
                            "must_not": {
                                "exists": {
                                    "field": "shop.is_tokonow"
                                }
                            }
                        }
                    },
                    "must": {
                        "query_string": {
                            "_name": "MUST_KEYWORD",
                            "default_operator": "AND",
                            "fields": [
                                "name"
                            ],
                            "fuzziness": "0",
                            "query": "(sticker OR stiker) (brand OR branding) (car OR mobil)"
                        }
                    },
                    "must_not": [
                        {
                            "exists": {
                                "_name": "NOT_EXC_VAR",
                                "field": "parent_id"
                            }
                        },
                        {
                            "exists": {
                                "_name": "NOT_ADULT",
                                "field": "is_adult"
                            }
                        }
                    ]
                }
            }
        }
    },
    "size": 60,
    "sort": [
        {
            "bestmatch_score_default": {
                "order": "desc"
            }
        }
    ],
    "track_total_hits": true
}
{
    "query": {
        "constant_score": {
            "filter": {
                "bool": {
                    "filter": [
                        {
                            "terms": {
                                "_name": "FILTER_BY_SHOP_ID",
                                "shop.id": [
                                    6873634
                                ]
                            }
                        },
                        {
                            "exists": {
                                "_name": "FILTER_ETALASE",
                                "field": "count_sold"
                            }
                        },
                        {
                            "bool": {
                                "_name": "FILTER_EXC_PARENT_ID",
                                "must_not": {
                                    "exists": {
                                        "field": "parent_id"
                                    }
                                }
                            }
                        }
                    ]
                }
            }
        }
    }
}
{
    "from": 140,
    "query": {
        "constant_score": {
            "filter": {
                "bool": {
                    "filter": [
                        {
                            "terms": {
                                "_name": "FILTER_BY_SHOP_ID",
                                "shop.id": [
                                    320124
                                ]
                            }
                        },
                        {
                            "bool": {
                                "_name": "FILTER_EXC_PARENT_ID",
                                "must_not": {
                                    "exists": {
                                        "field": "parent_id"
                                    }
                                }
                            }
                        },
                        {
                            "bool": {
                                "_name": "FILTER_NEARBY_SHOPPAGE",
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
                                            "nearby_districts": 1714
                                        }
                                    }
                                ]
                            }
                        }
                    ],
                    "must_not": [
                        {
                            "terms": {
                                "_name": "NOT_BAN_BY_CAT",
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
                                "_name": "NOT_BAN_PLAYSTORE_KEYWORD",
                                "keyword_type": [
                                    "2",
                                    "3"
                                ]
                            }
                        }
                    ]
                }
            }
        }
    },
    "size": 10,
    "sort": [
        {
            "id": {
                "order": "desc"
            }
        }
    ],
    "timeout": "2s",
    "track_total_hits": true
}
{
    "query": {
        "bool": {
            "filter": [
                {
                    "bool": {
                        "must_not": {
                            "exists": {
                                "field": "parent_id"
                            }
                        }
                    }
                },
                {
                    "query_string": {
                        "default_operator": "and",
                        "escape": false,
                        "fields": [
                            "extended_name"
                        ],
                        "fuzziness": "0",
                        "query": "pioneer avh x8750bt"
                    }
                },
                {
                    "bool": {
                        "filter": {
                            "terms": {
                                "department_id": [
                                    2866,
                                    2868,
                                    2870,
                                    565,
                                    564,
                                    2867,
                                    562,
                                    2883,
                                    3305,
                                    4375
                                ]
                            }
                        }
                    }
                }
            ]
        }
    }
}
{
    "from": 0,
    "query": {
        "function_score": {
            "boost_mode": "replace",
            "functions": [
                {
                    "script_score": {
                        "script": {
                            "source": "doc['bestmatch_score_default'].value"
                        }
                    },
                    "weight": 1
                },
                {
                    "filter": {
                        "terms": {
                            "_name": "18",
                            "department_id": [
                                2121
                            ]
                        }
                    },
                    "script_score": {
                        "script": {
                            "source": "30"
                        }
                    }
                },
                {
                    "filter": {
                        "terms": {
                            "_name": "18",
                            "department_id": [
                                3854
                            ]
                        }
                    },
                    "script_score": {
                        "script": {
                            "source": "20"
                        }
                    }
                },
                {
                    "filter": {
                        "exists": {
                            "_name": "23",
                            "field": "is_cod"
                        }
                    },
                    "script_score": {
                        "script": {
                            "source": "2"
                        }
                    }
                },
                {
                    "filter": {
                        "term": {
                            "city_affinities": {
                                "_name": "29",
                                "value": 150
                            }
                        }
                    },
                    "script_score": {
                        "script": {
                            "source": "7"
                        }
                    }
                },
                {
                    "filter": {
                        "terms": {
                            "_name": "15",
                            "near_districts": [
                                1697
                            ]
                        }
                    },
                    "script_score": {
                        "script": {
                            "source": "4"
                        }
                    }
                },
                {
                    "filter": {
                        "terms": {
                            "_name": "15",
                            "nearer_districts": [
                                1697
                            ]
                        }
                    },
                    "script_score": {
                        "script": {
                            "source": "2"
                        }
                    }
                },
                {
                    "filter": {
                        "term": {
                            "warehouse.city_id": {
                                "_name": "5",
                                "value": 150
                            }
                        }
                    },
                    "script_score": {
                        "script": {
                            "source": "1"
                        }
                    }
                },
                {
                    "filter": {
                        "term": {
                            "free_cluster": {
                                "_name": "28",
                                "value": 117
                            }
                        }
                    },
                    "script_score": {
                        "script": {
                            "source": "4"
                        }
                    }
                },
                {
                    "filter": {
                        "term": {
                            "cheap_cluster": {
                                "_name": "28",
                                "value": 117
                            }
                        }
                    },
                    "script_score": {
                        "script": {
                            "source": "2"
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
                                    "bool": {
                                        "_name": "FILTER_NEARBY",
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
                                                    "nearby_districts": 1697
                                                }
                                            }
                                        ]
                                    }
                                },
                                {
                                    "bool": {
                                        "_name": "FILTER_NOW_GLOBAL_SRCH",
                                        "minimum_should_match": "1",
                                        "should": [
                                            {
                                                "bool": {
                                                    "must": [
                                                        {
                                                            "term": {
                                                                "shop.is_tokonow": true
                                                            }
                                                        },
                                                        {
                                                            "term": {
                                                                "warehouse.warehouse_id": 12941747
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
                                        ]
                                    }
                                }
                            ],
                            "must": {
                                "query_string": {
                                    "_name": "MUST_KEYWORD",
                                    "default_operator": "AND",
                                    "fields": [
                                        "name"
                                    ],
                                    "fuzziness": "0",
                                    "query": "fifa 23 pc"
                                }
                            },
                            "must_not": [
                                {
                                    "exists": {
                                        "_name": "NOT_EXC_VAR",
                                        "field": "parent_id"
                                    }
                                },
                                {
                                    "exists": {
                                        "_name": "NOT_ADULT",
                                        "field": "is_adult"
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
    "size": 60,
    "track_total_hits": true
}
{
    "from": 40,
    "query": {
        "constant_score": {
            "filter": {
                "bool": {
                    "filter": [
                        {
                            "terms": {
                                "_name": "FILTER_BY_SHOP_ID",
                                "shop.id": [
                                    9209916
                                ]
                            }
                        },
                        {
                            "bool": {
                                "_name": "FILTER_EXC_PARENT_ID",
                                "must_not": {
                                    "exists": {
                                        "field": "parent_id"
                                    }
                                }
                            }
                        },
                        {
                            "bool": {
                                "_name": "FILTER_NEARBY_SHOPPAGE",
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
                                            "nearby_districts": 3481
                                        }
                                    }
                                ]
                            }
                        }
                    ],
                    "must_not": [
                        {
                            "terms": {
                                "_name": "NOT_BAN_BY_CAT",
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
                                "_name": "NOT_BAN_PLAYSTORE_KEYWORD",
                                "keyword_type": [
                                    "2",
                                    "3"
                                ]
                            }
                        }
                    ]
                }
            }
        }
    },
    "size": 10,
    "sort": [
        {
            "price": {
                "order": "asc"
            }
        },
        {
            "id": {
                "order": "desc"
            }
        }
    ],
    "timeout": "2s",
    "track_total_hits": true
}