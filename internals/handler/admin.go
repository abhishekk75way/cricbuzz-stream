package handler

import (
	"encoding/json"
	"net/http"

	"github.com/abhishekk75way/cricbuzz-stream/internals/services"
)

type UpdateRequest struct {
	MatchID string `json:"matchId"`
	Score   string `json:"score"`
	Overs   string `json:"overs"`
}

func UpdateHandler(svc *services.MatchService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		svc.UpdateScore(req.MatchID, req.Score, req.Overs)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"score updated successfully"}`))
	}
}
