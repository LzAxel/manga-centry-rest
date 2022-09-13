package service

import (
	"mangacentry/models"
	"mangacentry/pkg/repository"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type Manga interface {
	Create(userId int, manga models.Manga) (int, error)
	GetById(mangaId int) (models.Manga, error)
	Update(userId, mangaId int, updateInput models.UpdateMangaInput) error
	Delete(userId, mangaId int) error
}

type Service struct {
	Authorization
	Manga
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo.Authorization),
		Manga:         NewMangaService(repo.Manga),
	}
}
