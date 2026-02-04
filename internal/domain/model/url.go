package model

import "time"

type URL struct {
	ID        int64
	LongURL   string
	Alias     string
	CreatedAt time.Time
}
