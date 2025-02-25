package models

type LogHistory struct {
	ID      int64
	UserId  int64
	Command string
	Count   int
}
