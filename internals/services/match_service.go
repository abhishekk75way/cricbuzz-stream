package services

import (
	"encoding/json"
	"time"

	"github.com/abhishekk75way/cricbuzz-stream/internals/models"
	"github.com/abhishekk75way/cricbuzz-stream/internals/sse"
)

type MatchService struct {
	Broker *sse.Broker
}

func (s *MatchService) UpdateScore(MatchID, Score, Overs string) {
	match := models.Match{
		ID:       MatchID,
		Score:    Score,
		Overs:    Overs,
		Status:   "LIVE",
		UpdateAt: time.Now(),
	}

	data, _ := json.Marshal(match)

	s.Broker.Broadcast(MatchID, string(data))
}
