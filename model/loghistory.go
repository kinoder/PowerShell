package models

import "time"

type LogHistory struct {
	ID        int64
	UserId    int64
	Command   string
	Count     int
	CreatedAt time.Time
}
