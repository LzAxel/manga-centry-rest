package repository

import (
	"mangacentry/internal/core"
	"mangacentry/pkg/repository/psql"

	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user core.User) (int, error)
	GetUser(username, password string) (core.User, error)
}

type Manga interface {
	GetById(mangaId int) (core.Manga, error)
	Create(userId int, manga core.Manga) (int, error)
	Update(userId, mangaId int, updateInput core.UpdateMangaInput) error
	Delete(userId, mangaId int) error
}

type Chapter interface {
	Create(userId int, chapter core.MangaChapter) (int, error)
	GetById(chapterId int) (core.MangaChapter, error)
	GetByMangaId(mangaId int) ([]core.ChapterListElement, error)
	GetByNumber(mangaId, chapterNumber int) (core.MangaChapter, error)
	Update(userId, chapterId int, updateInput core.UpdateChapterInput) error
	Delete(userId, chapterId int) error
}

type Image interface {
	Create(image core.ChapterImage) (int, error)
	GetByChapterId(chapterId int) ([]core.ChapterImage, error)
}

type Repository struct {
	Authorization
	Manga
	Chapter
	Image
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: psql.NewAuthPostgres(db),
		Manga:         psql.NewMangaPostgres(db),
		Chapter:       psql.NewChapterPostgres(db),
		Image:         psql.NewImagePostgres(db),
	}
}
