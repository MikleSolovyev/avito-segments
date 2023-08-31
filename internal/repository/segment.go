package repository

import (
	"avito-segments/internal/model"
	"avito-segments/pkg/db/postgres"
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

type SegmentRepo struct {
	*postgres.Postgres
}

func NewSegmentRepo(pg *postgres.Postgres) *SegmentRepo {
	return &SegmentRepo{pg}
}

func (r *SegmentRepo) CreateSegment(ctx context.Context, segment *model.Segment) error {
	const fn = "repository.SegmentRepo.CreateSegment"

	sql, args, err := r.Builder.
		Insert("segments").
		Columns("slug", "percent").
		Values(segment.Slug, segment.Percent).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	err = r.DB.QueryRow(ctx, sql, args...).Scan(&segment.Id)
	if err != nil {
		if r.IsEqualErrors(err, "23505") {
			return model.ErrSegmentAlreadyExists
		}

		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (r *SegmentRepo) DeleteSegmentBySlug(ctx context.Context, slug string) error {
	const fn = "repository.SegmentRepo.DeleteSegmentBySlug"

	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// check if segment exists
	sql, args, err := r.Builder.
		Select("1").
		From("segments").
		Where(sq.Eq{"slug": slug}).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	err = tx.QueryRow(ctx, sql, args...).Scan(nil)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.ErrSegmentDoesNotExist
		}

		return fmt.Errorf("%s: %w", fn, err)
	}

	// insert delete action in history
	sql, args, err = r.Builder.
		Insert("history (user_id, slug, is_deleted, executed_at)").
		Select(
			sq.Select("user_id",
				sq.Placeholders(1),
				"TRUE",
				"(CASE WHEN expired_at <= CURRENT_TIMESTAMP THEN expired_at ELSE CURRENT_TIMESTAMP END)").
				From("users_segments").
				Where(sq.Eq{"slug": slug})).
		ToSql()
	args = append(args, slug)

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	// delete segment
	sql, args, err = r.Builder.
		Delete("segments").
		Where(sq.Eq{"slug": slug}).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}
