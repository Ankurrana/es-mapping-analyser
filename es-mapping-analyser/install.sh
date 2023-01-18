#!/bin/bash

echo "Usage: $0 -es-port [es-port] -agent-port [agent-port] -api-port [api-port] -traffic [traffic]"

es_port=9200
agent_port=9500
api_port=8123
traffic_percent=10

while [[ $# -gt 0 ]]
do
key="$1"

case $key in
    -es-port)
    es_port="$2"
    shift
    ;;
    -agent-port)
    agent_port="$2"
    shift
    ;;
    -api-port)
    api_port="$2"
    shift
    ;;
    -traffic)
    traffic_percent="$2"
    shift
    ;;
    *)
    # unknown option
    ;;
esac
shift
done


echo "es-port: $es_port"
echo "agent-port: $agent_port"
echo "api-port: $api_port"
echo "traffic-percent: $traffic_percent"

# gor installation!
# bash <(curl -s 'http://10.41.0.156:8002/install/install.sh')
wget -nv 'https://github.com/buger/goreplay/releases/download/1.3.3/gor_1.3.3_x64.tar.gz' && tar -xzf gor_1.3.3_x64.tar.gz && mv gor /usr/local/bin/go_replay
chmod 755 /usr/local/bin/go_replay

# moving 
cp bin/es_mapping_analyser /usr/local/bin/es_mapping_analyser

#if service --status-all | grep -Fq 'sp002-profiler'; then    
service sp002-profiler stop    
#fi

#if service --status-all | grep -Fq 'es-mapping-analyzer'; then    
service es-mapping-analyzer stop    
#fi



# install systemd
cat >/etc/systemd/system/sp002-profiler.service <<EOL
[Unit]
Description=sp002 Profiling Agent
After=network.target
StartLimitBurst=12
StartLimitIntervalSec=60


[Service]
Restart=always
RestartSec=5
ExecStart=/usr/local/bin/go_replay -http-original-host \
    -http-allow-method GET \
    -http-allow-method POST \
    -http-allow-url _search \
    -http-allow-url _msearch \
    -http-allow-url _count \
    -http-disallow-url ^/\\. \
    -prettify-http \
    -stats \
    -output-http-stats \
    -input-raw :${es_port} \
    -input-raw-buffer-size 10485760 \
    -input-raw-realip-header X-Real-IP \
    -input-raw-engine raw_socket \
    -verbose 1 \
    -output-http 'http://localhost:${agent_port}|${traffic_percent}%' 

[Install]
WantedBy=multi-user.target

EOL


# install systemd
cat >/etc/systemd/system/es-mapping-analyzer.service <<EOL
[Unit]
Description=Elasticsearch Mapping Analyser
After=network.target
StartLimitBurst=12
StartLimitIntervalSec=60

[Service]
Restart=always
RestartSec=5
ExecStart=/usr/local/bin/es_mapping_analyser \
    -port ${api_port} \
    -agent-port ${agent_port} \
    --es-url 127.0.0.1:${es_port}

[Install]
WantedBy=multi-user.target

EOL

systemctl daemon-reload
systemctl enable es-mapping-analyzer
systemctl restart es-mapping-analyzer



systemctl enable sp002-profiler
systemctl restart sp002-profiler



echo "to stop analyser: service es-mapping-analyzer stop"
echo "to stop profiler: service sp002-profiler stop"
echo "to check report: curl localhost:${api_port}/report"
echo "EMA installation complete thanks!"



