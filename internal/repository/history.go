package repository

import (
	"avito-segments/internal/model"
	"avito-segments/pkg/db/postgres"
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
)

type HistoryRepo struct {
	*postgres.Postgres
}

func NewHistoryRepo(pg *postgres.Postgres) *HistoryRepo {
	return &HistoryRepo{pg}
}

func (r *HistoryRepo) GetHistoryByDate(ctx context.Context, from, to int) ([]model.History, error) {
	const fn = "repository.HistoryRepo.GetHistoryByDate"

	sql, args, err := r.Builder.
		Select("user_id", "slug", "is_deleted", "executed_at").
		From("history").
		Where(sq.And{
			sq.Expr("executed_at >= to_timestamp($1)", from),
			sq.Expr("executed_at <= to_timestamp($2)", to),
		},
		).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	rows, err := r.DB.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	defer rows.Close()

	var records []model.History
	for rows.Next() {
		var record model.History
		err = rows.Scan(
			&record.UserId,
			&record.Slug,
			&record.IsDeleted,
			&record.ExecutedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}

		records = append(records, record)
	}

	return records, nil
}
