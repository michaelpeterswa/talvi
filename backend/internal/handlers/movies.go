package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/michaelpeterswa/talvi/backend/internal/movies"
)

type MovieHandler struct {
	MoviesClient *movies.MoviesClient
}

func NewMovieHandler(moviesClient *movies.MoviesClient) *MovieHandler {
	return &MovieHandler{
		MoviesClient: moviesClient,
	}
}

func (mh *MovieHandler) GetMovies(w http.ResponseWriter, r *http.Request) {
	moviesCtx, cancelFunc := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancelFunc()

	traceCtx, span := tracer.Start(moviesCtx, "GetMovies")
	defer span.End()

	movies, err := mh.MoviesClient.GetMovies(traceCtx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	moviesJSON, err := json.Marshal(movies)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(moviesJSON)
}
