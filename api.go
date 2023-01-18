package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ankur-toko/es-mapping-analyser/reports"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/report", GetOptimizationReport).Methods("GET")
	r.HandleFunc("/run", RunAnalysis).Methods("GET")
	return r
}

func GetOptimizationReport(rw http.ResponseWriter, r *http.Request) {
	log.Print("Getting ES Analysis Report")
	identifier := r.URL.Query().Get("id")
	report := ReadReport(identifier)
	fmt.Fprint(rw, report)
}

func ReadReport(id string) string {
	return reports.GetReportFor(id)
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
