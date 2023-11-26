package usecase

import "github.com/markgregr/RIP/internal/http/repository"

type UseCase struct {
	Repository *repository.Repository
}

func NewUseCase(repo *repository.Repository) *UseCase {
	return &UseCase{
		Repository: repo,
	}
}
