package model

import "errors"

var (
	ErrUserAlreadyInSegment = errors.New("cannot add user in segment: user already in segment")
	ErrUserNotInSegment     = errors.New("cannot delete user from segment: user not in segment")
	ErrSegmentAlreadyExists = errors.New("segment already exists")
	ErrSegmentDoesNotExist  = errors.New("segment does not exist")

	ErrWrongPeriod        = errors.New("from must be before to")
	ErrUpdateIntersection = errors.New("added and deleted slugs sets must not intersect")
	ErrEmptyRequestBody   = errors.New("empty request body")
)
