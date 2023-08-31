package service

import (
	"avito-segments/internal/model"
	"context"
)

type SegmentRepo interface {
	CreateSegment(context.Context, *model.Segment) error
	DeleteSegmentBySlug(context.Context, string) error
}

type SegmentService struct {
	repo SegmentRepo
}

func NewSegmentService(repo SegmentRepo) *SegmentService {
	return &SegmentService{
		repo: repo,
	}
}

func (s *SegmentService) CreateSegment(ctx context.Context, segment *model.Segment) error {
	return s.repo.CreateSegment(ctx, segment)
}

func (s *SegmentService) DeleteSegmentBySlug(ctx context.Context, slug string) error {
	return s.repo.DeleteSegmentBySlug(ctx, slug)
}
