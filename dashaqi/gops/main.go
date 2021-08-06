package main

import (
	"github.com/google/gops/agent"
	"log"
	"net/http"
)

func main() {
	if err := agent.Listen(agent.Options{}); err != nil {
		log.Fatalf("agent.listen err: %v", err)
	}

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`hello world `))
	})
	_ = http.ListenAndServe(":6060", http.DefaultServeMux)
}
