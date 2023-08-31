package repository

import "avito-segments/pkg/db/postgres"

type Repository struct {
	*UserRepo
	*SegmentRepo
	*HistoryRepo
}

func New(pg *postgres.Postgres) *Repository {
	return &Repository{
		UserRepo:    NewUserRepo(pg),
		SegmentRepo: NewSegmentRepo(pg),
		HistoryRepo: NewHistoryRepo(pg),
	}
}
