package repository

import (
	"avito-segments/internal/model"
	"avito-segments/pkg/db/postgres"
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

type UserRepo struct {
	*postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

func (r *UserRepo) createUserInTx(ctx context.Context, user *model.User, tx pgx.Tx) error {
	const fn = "repository.UserRepo.createUserInTx"

	sql, args, err := r.Builder.
		Insert("users").
		Columns("id").
		Values(user.Id).
		Suffix("ON CONFLICT DO NOTHING").
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (r *UserRepo) GetUserSegments(ctx context.Context, user *model.User) ([]string, error) {
	const fn = "repository.UserRepo.GetUserSegments"

	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// create user, if he exists - nothing will change
	err = r.createUserInTx(ctx, user, tx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	// get user segments
	sql, args, err := r.Builder.
		Select("slug").
		From("users_segments").
		Where(sq.And{
			sq.Eq{"user_id": user.Id},
			sq.Or{
				sq.Eq{"expired_at": nil},
				sq.Expr("expired_at >= CURRENT_TIMESTAMP"),
			},
		}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	rows, err := tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	defer rows.Close()

	var slugs []string
	for rows.Next() {
		var slug string
		err = rows.Scan(&slug)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}

		slugs = append(slugs, slug)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return slugs, nil
}

func (r *UserRepo) UpdateUserSegments(ctx context.Context, user *model.User, addSlugs map[string]string, delSlugs map[string]struct{}) error {
	const fn = "repository.UserRepo.UpdateUserSegments"

	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// create user, if he exists - nothing will change
	err = r.createUserInTx(ctx, user, tx)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	// get current user segments
	sql, args, err := r.Builder.
		Select("slug").
		From("users_segments").
		Where(sq.And{
			sq.Eq{"user_id": user.Id},
			sq.Or{
				sq.Eq{"expired_at": nil},
				sq.Expr("expired_at >= CURRENT_TIMESTAMP"),
			},
		}).ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	rows, err := tx.Query(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}
	defer rows.Close()

	curSlugsSet := make(map[string]struct{})
	for rows.Next() {
		var slug string
		err = rows.Scan(&slug)
		if err != nil {
			return fmt.Errorf("%s: %w", fn, err)
		}

		curSlugsSet[slug] = struct{}{}
	}

	// check for ErrUserAlreadyInSegment and add segments
	for slug, expiredAt := range addSlugs {
		if _, ok := curSlugsSet[slug]; ok {
			return fmt.Errorf("%w %s", model.ErrUserAlreadyInSegment, slug)
		}

		sql, args, err = r.Builder.
			Insert("users_segments (user_id, slug, expired_at)").
			Values(sq.Expr(fmt.Sprintf("%s, %s", sq.Placeholders(2), expiredAt), user.Id, slug)).
			ToSql()
		if err != nil {
			return fmt.Errorf("%s: %w", fn, err)
		}

		_, err = tx.Exec(ctx, sql, args...)
		if err != nil {
			return fmt.Errorf("%s: %w", fn, err)
		}
	}

	// check for ErrUserNotInSegment and delete segments
	for slug := range delSlugs {
		if _, ok := curSlugsSet[slug]; !ok {
			return fmt.Errorf("%w %s", model.ErrUserNotInSegment, slug)
		}

		sql, args, err = r.Builder.
			Delete("users_segments").
			Where(sq.Eq{"user_id": user.Id, "slug": slug}).
			ToSql()
		if err != nil {
			return fmt.Errorf("%s: %w", fn, err)
		}

		_, err = tx.Exec(ctx, sql, args...)
		if err != nil {
			return fmt.Errorf("%s: %w", fn, err)
		}
	}

	// save additions in history
	for slug := range addSlugs {
		sql, args, err = r.Builder.
			Insert("history (user_id, slug, is_deleted, executed_at)").
			Values(sq.Expr(fmt.Sprintf("%s, FALSE, CURRENT_TIMESTAMP", sq.Placeholders(2)), user.Id, slug)).
			ToSql()

		_, err = tx.Exec(ctx, sql, args...)
		if err != nil {
			return fmt.Errorf("%s: %w", fn, err)
		}
	}

	// save deletions in history
	for slug := range delSlugs {
		sql, args, err = r.Builder.
			Insert("history (user_id, slug, is_deleted, executed_at)").
			Values(sq.Expr(fmt.Sprintf("%s, TRUE, CURRENT_TIMESTAMP", sq.Placeholders(2)), user.Id, slug)).
			ToSql()

		_, err = tx.Exec(ctx, sql, args...)
		if err != nil {
			return fmt.Errorf("%s: %w", fn, err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}
