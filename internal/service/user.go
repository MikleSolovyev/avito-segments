package service

import (
	"avito-segments/internal/model"
	"context"
	"fmt"
)

type UserRepo interface {
	GetUserSegments(context.Context, *model.User) ([]string, error)
	UpdateUserSegments(context.Context, *model.User, map[string]string, map[string]struct{}) error
}

type UserService struct {
	repo UserRepo
}

func NewUserService(repo UserRepo) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) GetUserSegments(ctx context.Context, user *model.User) ([]string, error) {
	return s.repo.GetUserSegments(ctx, user)
}

func (s *UserService) UpdateUserSegments(ctx context.Context, user *model.User, addSlugs map[string]int, delSlugs []string) error {
	// check if both empty
	if len(addSlugs) == 0 && len(delSlugs) == 0 {
		return model.ErrEmptyRequestBody
	}

	addSlugsRepo := make(map[string]string)
	for slug, ttl := range addSlugs {
		if ttl == 0 {
			addSlugsRepo[slug] = "NULL"
		} else {
			addSlugsRepo[slug] = fmt.Sprintf("CURRENT_TIMESTAMP + INTERVAL '%d seconds'", ttl)
		}
	}

	// check intersection of added and deleted slugs
	delSlugsSet := make(map[string]struct{})
	for _, slug := range delSlugs {
		if _, ok := addSlugsRepo[slug]; ok {
			return model.ErrUpdateIntersection
		}

		delSlugsSet[slug] = struct{}{}
	}

	return s.repo.UpdateUserSegments(ctx, user, addSlugsRepo, delSlugsSet)
}
