package sse

import (
	"net/http"
	"sync"
)

type Client chan string

type Broker struct {
	clients map[string]map[Client]bool
	mu      sync.RWMutex
}

func NewBroker() *Broker {
	return &Broker{
		clients: make(map[string]map[Client]bool),
	}
}

func (b *Broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	matchID := r.URL.Query().Get("matchId")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	client := make(Client)

	b.mu.Lock()
	if b.clients[matchID] == nil {
		b.clients[matchID] = make(map[Client]bool)
	}
	b.clients[matchID][client] = true
	b.mu.Unlock()

	defer func() {
		b.mu.Lock()
		delete(b.clients[matchID], client)
		b.mu.Unlock()
		close(client)
	}()

	for {
		select {
		case message := <-client:
			_, err := w.Write([]byte("data: " + message + "\n\n"))
			if err != nil {
				return
			}
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}

func (b *Broker) Broadcast(matchID string, message string) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for client := range b.clients[matchID] {
		client <- message
	}
}
