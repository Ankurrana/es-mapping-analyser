package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

var ApiPort, AgentPort *int
var ESURL *string

func main() {
	ApiPort = flag.Int("port", 80, "Application port")
	AgentPort = flag.Int("agent-port", 9500, "agent listener")
	ESURL = flag.String("es-url", "localhost:9200", "Port to listen on")
	flag.Parse()
	go StartWebServer(*ApiPort)
	time.Sleep(2 * time.Second)
	Run(*ESURL, *AgentPort)

	// Blocking go code forever
	ch := make(chan bool)
	<-ch
}

func StartWebServer(port int) {
	log.Printf("Initializing Webapp on port %v", port)
	r := Router()
	err := http.ListenAndServe(fmt.Sprintf(":%v", port), r)
	if err != nil {
		log.Fatalf("WebApp failed with error %v", err.Error())
	}
}
