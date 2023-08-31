package model

import "time"

type Segment struct {
	Id        int
	Slug      string
	Percent   int
	DeletedAt time.Time
}
