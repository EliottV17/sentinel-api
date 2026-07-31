package config

import "os"

type Config struct {
	DatabaseUrl string
	Concurrency int
}

func Load() Config {
	concurrency := 10
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		url = "postgres://postgres:postgres@127.0.0.1:5432/sentinel_db"
	}
	return Config{
		DatabaseUrl: url,
		Concurrency: concurrency,
	}
}