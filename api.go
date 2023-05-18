package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ankur-toko/es-mapping-analyser/reports"
	"github.com/ankur-toko/es-mapping-analyser/webui"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/report", GetOptimizationReport).Methods("GET")
	r.HandleFunc("/report-ui", GetOptimizationReportUI).Methods("GET")
	r.HandleFunc("/run", RunAnalysis).Methods("GET")
	return r
}

func GetOptimizationReportUI(rw http.ResponseWriter, r *http.Request) {
	log.Print("Getting ES Analysis Report")
	index_regex := r.URL.Query().Get("index_regex")
	update_mapping := r.URL.Query().Get("update_mapping")
	var err error
	updateMapping := false
	if update_mapping != "" {
		updateMapping, err = strconv.ParseBool(update_mapping)
		if err != nil {
			fmt.Fprint(rw, "update_mapping is a boolean variable, try true,false")
		}
	}

	report := ReadReportUI(index_regex, updateMapping)
	fmt.Fprint(rw, report)
}

func GetOptimizationReport(rw http.ResponseWriter, r *http.Request) {
	log.Print("Getting ES Analysis Report")
	index_regex := r.URL.Query().Get("index_regex")
	update_mapping := r.URL.Query().Get("update_mapping")
	var err error
	updateMapping := false
	if update_mapping != "" {
		updateMapping, err = strconv.ParseBool(update_mapping)
		if err != nil {
			fmt.Fprint(rw, "update_mapping is a boolean variable, try true,false")
		}
	}

	report := ReadReport(index_regex, updateMapping)
	fmt.Fprint(rw, report)
}

func ReadReport(index_regex string, update_mapping bool) string {
	return reports.GetReportFor(index_regex, update_mapping)
}

func ReadReportUI(index_regex string, update_mapping bool) string {
	reportMap := reports.GetQUMMapReportFor(index_regex, update_mapping)
	return webui.QURMapToGraph(reportMap)
}

func RunAnalysis(rw http.ResponseWriter, r *http.Request) {

	agentPort := r.URL.Query().Get("agent-port")
	esURL := r.URL.Query().Get("es-url")

	if agentPort == "" {
		agentPort = "9500"
	}

	if esURL == "" {
		esURL = "localhost:9200"
	}
	port, _ := strconv.ParseInt(agentPort, 10, 32)
	Run(esURL, int(port))
	rw.Write([]byte(fmt.Sprintf("Analysing ES: %v\nAccepting input on port %v", esURL, port)))
}

func Run(ESURL string, port int) {
	go reports.RunAnalysis(ESURL, port)
	time.Sleep(time.Second)
	log.Printf("Analysing ES: %v\nAccepting es requests on port %v\n", ESURL, port)
}
