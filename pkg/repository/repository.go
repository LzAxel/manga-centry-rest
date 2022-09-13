package repository

import (
	"mangacentry/models"

	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GetUser(username, password string) (models.User, error)
}

type Manga interface {
	GetById(mangaId int) (models.Manga, error)
	Create(userId int, manga models.Manga) (int, error)
	Update(userId, mangaId int, updateInput models.UpdateMangaInput) error
	Delete(userId, mangaId int) error
}

type Repository struct {
	Authorization
	Manga
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Manga:         NewMangaPostgres(db),
	}
}
