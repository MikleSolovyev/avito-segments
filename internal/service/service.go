package service

type Repository interface {
	UserRepo
	SegmentRepo
	HistoryRepo
}

type Service struct {
	*UserService
	*SegmentService
	*HistoryService
}

func New(repo Repository) *Service {
	return &Service{
		UserService:    NewUserService(repo),
		SegmentService: NewSegmentService(repo),
		HistoryService: NewHistoryService(repo),
	}
}
