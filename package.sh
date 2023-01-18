GOOS=linux GOARCH=amd64 go build -o ./es-mapping-analyser/bin/es_mapping_analyser
tar czf es_mapping_analyser.tar.gz es-mapping-analyser
