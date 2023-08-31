package service

import (
	"avito-segments/internal/model"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
)

type HistoryRepo interface {
	GetHistoryByDate(context.Context, int, int) ([]model.History, error)
}

type HistoryService struct {
	repo HistoryRepo
}

func NewHistoryService(repo HistoryRepo) *HistoryService {
	return &HistoryService{
		repo: repo,
	}
}

func (s *HistoryService) GetHistoryByDate(ctx context.Context, from, to int, dest io.Writer) error {
	if from >= to {
		return model.ErrWrongPeriod
	}

	records, err := s.repo.GetHistoryByDate(ctx, from, to)
	if err != nil {
		return err
	}

	err = s.writeToCSV(records, dest)
	if err != nil {
		return err
	}

	return nil
}

func (s *HistoryService) writeToCSV(records []model.History, dest io.Writer) error {
	const fn = "service.HistoryService.writeToCSV"

	w := csv.NewWriter(dest)
	defer w.Flush()

	for _, record := range records {
		var stringRecord []string
		action := "add"
		if record.IsDeleted {
			action = "del"
		}
		stringRecord = append(stringRecord, strconv.Itoa(record.UserId), record.Slug, action, record.ExecutedAt.String())

		err := w.Write(stringRecord)
		if err != nil {
			return fmt.Errorf("%s: %w", fn, err)
		}
	}

	return nil
}
