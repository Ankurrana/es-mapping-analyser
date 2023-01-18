#!/bin/bash

# Script to accept parameters

# Assign default values to variables
port=8080
es_port=9200
sample_size=0.1
output_file="output.txt"
index_filter="*"
only_agent=false
output_url=""
duration="1h"

# Loop through the input arguments
while [ $# -gt 0 ]; do
  case "$1" in
    --port)
      port="$2"
      shift 2
      ;;
    --es-port)
      es_port="$2"
      shift 2
      ;;
    --sample-size)
      sample_size="$2"
      shift 2
      ;;
    --output-file)
      output_file="$2"
      shift 2
      ;;
    --index-filter)
      index_filter="$2"
      shift 2
      ;;
    --only-agent)
      only_agent=true
      shift 1
      ;;
    --output-url)
      output_url="$2"
      shift 2
      ;;
    --duration)
      duration="$2"
      shift 2
      ;;
    *)
      echo "Unknown option: $1"
      exit 1
      ;;
  esac
done

# Print the final values of the variables
echo "Port: $port"
echo "Elasticsearch Port: $es_port"
echo "Sample Size: $sample_size"
echo "Output File: $output_file"
echo "Index Filter: $index_filter"
echo "Only Agent: $only_agent"
echo "Output URL: $output_url"
echo "Duration: $duration"



# START GO REPLAY LISTENING ON es_port 9200 AND   