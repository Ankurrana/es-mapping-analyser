# Elasticsearch Mapping Analyzer
## Introduction
The Elasticsearch Mapping Analyzer (EMA) is a tool that automates the process of identifying potential optimizations for any Elasticsearch index, based on set of queries that are run on that index.

EMA creates a usage map for each field in the index, based on the input queries. By analyzing this usage map, the tool can identify fields that can be converted to a more optimized type or configuration, which can save disk space, improve indexing speed, and potentially enhance search performance.

## Installation
EMA is installed on a machine that handles requests (usually master or coordinator nodes) in a Elasticsearch cluster.

## Usage
EMA comes with a profiler that intercepts a percentage of the queries sent to the Elasticsearch server and forwards them to an internal analyzer. This analyzer then creates a usage map for each field for every index of the cluster.

The analyzer is equipped with the knowledge needed to identify fields that can be converted to a more optimized type without impacting any existing queries. These recommendations are made available through an API and can be accessed in real-time. For most production systems, if the analyzer is run for a sufficient amount of time, it should cover most use cases that occur on the cluster.

By default, the analyzer captures only 10% of requests sent to the Elasticsearch server. You can configure this setting to intercept higher %age of requests for clusters with low RPS.

## Installation
Download the latest version of EMA from the releases. The package includes a binary file and an installation script.

1. EMA uses GoReplay as a sniffing tool on the index, which is installed as a separate service called sp002-profiler. Another service, es-mapping-analyser, is installed to analyze the intercepted requests.
2. To install EMA, run the installation script with the following command:    
    `./install.sh -es-port [es-port] -agent-port [agent-port] -api-port [api-port] -traffic [traffic]`

3. The installation command takes four optional parameters as input:
    * es-port: The Elasticsearch port that needs to be sniffed.
    * agent-port: An internal port used by the profiler and analyzer. Default is 9500
    * api-port: The port on which the EMA report will be accessible via the URL http://localhost:[api-port]/report.
    * traffic: The percentage of traffic to be sniffed.

4. After the installation is complete, the sp002-profiler and es-mapping-analyzer services will be started automatically.


## How to read results
EMA exposes an API `/report` to check the current status of the analysis and recommendations

Here is the sample report
```
Assuming you have a Index named "products" with mapping like
{
	"name" : {
		"type: "text"
	},
	"cost" : {
		"type" : long
	},
	"shop" : {
		"properties" : {
			"id" : {
				"type" : "long"
			}
			"name" : {
				"type" : "keyword"
			}
		}
	},
	"product_description" : {
		"type" : "text"
	}
}

and your usecases are something like

1. Text matches on name
2. Sort queries on cost
3. Term matches and aggregations on shop.id 


Assuming EMA was able to capture all the usecases, the analysis report would look something like this: 

{
	"products": {
		"name": "products",
		"queries_analysed": 100,
		"optimizations": {
			"+OBJECT": [
				"shop.name",
				"product_description"
			],
			"-DOCVALUES": [
				"shop.name"
			],
			"-INDEX": [
				"cost"
			],
            "+KEYWORD": [
                "shop.id"
            ]
		},
		"usage_frequency_map": {
			"name": {
				"match": 70
			},
			"cost": {
				"sort": 70
			},
			"shop.id" : {
				"terms" : 10,
				"terms_aggs" : 20 
			}
		},
		"hash_code": "3203411xxx-1-1-2"
	}
}
```

For each index for which queries were recieved, it returns an object with 3 critical parts, queries_analysed, optimizations and hash_code. 

* queries_analysed is the number of queries analysed to generate this report. 
* usage_frequency_map lists how each of the field were used in the sampled queries. Fields with no use were ignored in this map
* optimizations: As of how, EMA gives out 4 possible optimizations "+OBJECT" , "-DOCVALUES" , "-INDEX", "+KEYWORD"
    - "+OBJECT" : these fields are recommended to be converted to type:object and enabled:false as no use was identified for these fields.
    - "-DOCVALUES" : fields for which we can disable doc_values are listed in this section `doc_values:false`
    - "-INDEX" : fields that doesn't require creating inverted index. `index:false`
    - "+KEYWORD" : this is useful for those numeric fields which are not utilzed for range/sort queries. such fields can be coverted to `type:keyword`
* Hashcode is a simple hash of the optimization map, so as to quickly identify if there are any changes in the results from the previous tests instead of manually checking for any new changes in the result.



# Options 
1. Report Request Param: index_regex. Accepts regex to get report for indices matching regex
`http://localhost:8123/report?index_regex=product_retail_v5`

2. Update Mapping Param: If you have updated the mapping of the index and need to get the report with updated mapping use update_mapping=true param.
`http://localhost:8123/report?update_mapping=true`

3. WebUI for usage frequency map can be seen using the below url
	`http://localhost:8123/report-ui`
	![Example](/example_graph.png "Example Usage Frequency Map")



# Architecture
![EMA Architecture Diagram](/ema.svg "EMA Architecture Diagram")


# Compatibility
Currently supports ES6 and ES7



# Build Command
```
	./package.sh
	
```

## Contact
Contact ankur.rana@tokopedia.com for any concerns or improvements