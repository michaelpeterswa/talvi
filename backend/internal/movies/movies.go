package movies

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/michaelpeterswa/talvi/backend/internal/cockroach"
	"github.com/michaelpeterswa/talvi/backend/internal/dragonfly"
	"github.com/redis/go-redis/v9"
)

type MoviesClient struct {
	kv *dragonfly.DragonflyClient
	db *cockroach.CockroachClient
}

type Movie struct {
	Name string `json:"name"`
	Year int    `json:"year"`
}

func NewMoviesClient(kv *dragonfly.DragonflyClient, db *cockroach.CockroachClient) *MoviesClient {
	return &MoviesClient{
		kv: kv,
		db: db,
	}
}

func (mc *MoviesClient) GetMovies(ctx context.Context) ([]Movie, error) {
	moviesStr, err := mc.kv.Client.Get(ctx, "movies").Result()
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("failed to get movies from dragonfly: %w", err)
	} else if moviesStr != "" && err != redis.Nil {
		var movies []Movie
		err := json.Unmarshal([]byte(moviesStr), &movies)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal movies: %w", err)
		}
		return movies, nil
	}
	rows, err := mc.db.Client.Query(ctx, "SELECT * FROM talvi.movies LIMIT 10;")
	if err != nil {
		return nil, fmt.Errorf("failed to query movies: %w", err)
	}

	var movies []Movie
	for rows.Next() {
		var movie Movie
		err := rows.Scan(&movie.Name, &movie.Year)
		if err != nil {
			return nil, fmt.Errorf("failed to scan movie: %w", err)
		}
		movies = append(movies, movie)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over movies: %w", err)
	}

	moviesJSON, err := json.Marshal(movies)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal movies: %w", err)
	}

	err = mc.kv.Client.Set(ctx, "movies", moviesJSON, time.Minute).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to set movies in dragonfly: %w", err)
	}

	return movies, nil
}
