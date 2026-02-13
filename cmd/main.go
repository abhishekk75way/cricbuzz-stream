package main

import (
	"log"
	"net/http"

	"github.com/abhishekk75way/cricbuzz-stream/internals/handler"
	"github.com/abhishekk75way/cricbuzz-stream/internals/services"
	"github.com/abhishekk75way/cricbuzz-stream/internals/sse"
)

func main() {
	broker := sse.NewBroker()
	matchService := &services.MatchService{Broker: broker}

	http.Handle("/events", broker)
	http.HandleFunc("/admin/update-score", handler.UpdateHandler(matchService))

	log.Println("SSE server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
