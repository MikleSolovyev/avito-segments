package model

import "time"

type History struct {
	Id         int
	UserId     int
	Slug       string
	IsDeleted  bool
	ExecutedAt time.Time
}
