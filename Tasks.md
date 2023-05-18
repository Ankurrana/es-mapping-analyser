**ES Mapping Profiler**

Inputs the queries used on an index, it finds appropriate mapping for an ES index which can be used as an optimizer for the mapping definition



Functional Requirements
1. Read queries used on an index [Done]
2. Determines an appropriate mapping for the index [Done]
3. Store result in a file 
4. Read Mapping to analyse
5. Add Close method in the Fetcher Library [DONE]
6. How about stored_fields? [DONE]
   // Stored Fields are generally discouraged, instead the recommendation is to use _source or doc_values 
7. How about fields which are not searched but only sorted on? Should we do index:false
    // We can set index false only if the field is only being used in sort, or source. In every other case we need index:true

8. What is the difference between keyword with (index:false,doc_values:false) and object with enabled:False? 
9. If a numeric field is being sorted on and no range queries, then we cannot convert it to keyword as it will change its sorting order [DONE]

9. What is object and object with enabled:false? [OBJECT TYPE ACTUALLY PARSES JSON INTO FIELDS BUT DOES'T INDEX][DONE]
10. What is text with index:false?
11. How could we identify geo_point? 
12. How to identify range of numeric fields?

13. Frequency of each field? 
14. Size of each field?
15. Add a check while analysing queries if all queries for a particular field are always boolean

16. Check if keywords are actually cheaper than long? [YES, PREFER KEYWORD FOR EXACT MATCHES][DONE] 
18. Do we need doc_values for aggregation queries? [YES] [DONE]
19. Study the frequency of usage for every fields, we may be able to decide which usecases are rarely used.
20. Difference between removing mapping entirely or keep it in source? 
21. Can we use sort with keyword for numeric type fields [NO, WE CANNOT] [DONE]
22. A system which can quickly test any 2 queries bruteforcelly
23. Send some data to newrelic and create charts for anybody
24. Develop an aggregator of the analysis and push the report in a file 
25. Invert Optimizations like I want to know all fields which are not used at all, and fields for which we can remove doc_values! [Done]

Improvements:
1. Can we get the frequency of most used item? [DONE]
2. Validate multi fields support [DONE]
3. Doesn't support all aggregations like multi_terms // (todo)
5. Instead of this service reading  queries from DB, we can create a server which accepts queries and generate usage map in realtime? [Done]
6. Need to evaluate all the alises of the index under analysis [DONE]
9. Fields which have docvalues false are also recommended to remove docvalues, which is unnecessary [DONE]
10. Fields which are only sorted cannot be converted to object, rather they need to be some type with index:false. [DONE]
11. Can we sort on fields where index:false is false? [YES It works] [DONE]
12. How does sorting for numeric fields work if type changed to keyword! [It sorts in alphanumeric, so we cannot convert numeric ] [DONE]



14. How about exists query? There are fields which only need run exists query! [We need to index such fields but no need to create doc_values] [DONE]
15. Do fields using range queries only need doc_values? [Yes, range queries may use doc_values sometimes on the basis of data and internal scoring system of ES] [DONE]
16. Can we also exclude fields from the source in shop index? [DONE] [YES, IT IS POSSIBLE IF WE DO NOT NEED SOME OF THE FIELDS]

17. query_string parameter actually searches in all fields mentioned in the fields param. [IMP][DONE]

18. Dont wait in the tail file indefinitely if the number of rows are less than batch size
19. Check why the count is incorrect in the report [Its correct][DONE] 
20. Include all terms level queries https://www.elastic.co/guide/en/elasticsearch/reference/master/term-level-queries.html




21. Currently it does't support multiple indices by default... for handling product_inactive and other's I need to run the function separately... I don't want that anymore..fix that. [Done]
22. I should not be required to hard code anything.System should handle that itself.[Done]
23. Only accept connection with path _search, _count and _msearch 

24. Generate HashValue and Count of fields for every type of field type change, this can be used as a quick check if there was any change in the results of the optimization [Done]


25. Fix the order of keys in hashcode of the query report [DONE]
26. Add a debug flag in the tool


27. Remove msearch queries from getting logged [DONE]
28. Do not print no incoming messages everytime, once per 5 seconds (maybe once every minute if no incoming messages) [Done]
29. Add a packaging system with eS mapper [Done]
30. Add this to github [Done]


31. What happens when ES is not reachable? or doesn't respond to mappings request? 
32. Logs first 10 queries as a debugging mechanism to be used when installing the agent [Done]

33. Support ES6 as well  [Done]
34. Add support for scroll queries and search_after [Done]







35. Recommend the best index sort strategy for each index. 
36. For fields with many fields, recommend if the entire object can be converted to type:object 
37. Find patterns of queries from the data. 
38. Updated mappings do not reflect correctly, better to reload them in some time [Done]
39. Validate if dynamic is enabled or not, recommend to disable it. [Done]
41. sp002 profiler sometimes stops, need to find the reason. 
42. Add a new boolean request_param update_mapping which when true will refetch updated mapping from the ES and return updated mapping [Done]

43. Add a new Static recommendation engine which performs static analysis on the basis of mapping definition. for ex. dyanmic should be strict or false and other static settings like default index sort is also not available in the mapping object.   


44. Add Benchmarking results in the Readme.  
46. If usage is match, then also no need to index, set doc_values:false [Done]
47. Can we recommend mapping as well instead of giving out result? 
48. Add a --debug clause to sp002 profiler as well. 


49. 


Analyzer


1. Use gor to capture realtime traffic and push it out to a file/tcp port/ http port.
2. A different binary which reads data from the output stream in step 1, analyzes it. 
3. It exposes report in json format on a different port which user can access directly. 
4. Or At the end of the process, it can output a json report to a file


./es_query_analyzer
    --port 8080    // gor to do post to this port itself
    --es-port 9200  // the port on which eqa sniffs
    --sample-size 0.1 // a number between 0 and 1 to capture the percent of queries on your cluster
    --output-file <CLUSTER_ID>
    --index-filter * // considers all indexes by default 
    --only-agent  // use eqa to only sniff data
    --output-url 156:8050 // Only if you are enabling --only-agent
    --duration 1h  // default duration for which you want to run eqa .. no need to run it forever
    
    


Check current report at: 
-- localhost:8080/report
:D Thinking on how you can improve your ES... 


What would user do
1. I'll give them a to download the latest version of eqa, which will be a gzip file with a few files and a binary
2. Once you have a binary, you can run it on your machine recieving direct traffic on ES
3. 


.sh file
1. downloads and installs gor
2. Go application spawns a server on output port and starts recieving traffic
3. Start gor to capture /search,/count and _msearch traffic on es-port
4. Current report will be exposed to --output-port /report
5. If user does Ctrl+C, a json_report will be exported to a file 



 


# Custom Query Analyser

1. Identifies the conditions which are most common in the set of all queries that we run. Can we identify the part of the search space which aren most commonly accessed?
2. If we can identify those then 





