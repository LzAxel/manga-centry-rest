package service

import (
	"mangacentry/internal/core"
	"mangacentry/pkg/repository"
)

type Authorization interface {
	CreateUser(user core.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type Manga interface {
	Create(userId int, manga core.Manga) (int, error)
	GetById(mangaId int) (core.Manga, error)
	Update(userId, mangaId int, updateInput core.UpdateMangaInput) error
	Delete(userId, mangaId int) error
}

type Chapter interface {
	Create(userId int, input core.CreateChapterInput) (int, error)
	GetById(chapterNumber int) (core.MangaChapter, error)
	GetByMangaId(mangaId int) ([]core.ChapterListElement, error)
	GetByNumber(mangaId, chapterNumber int) (core.MangaChapter, error)
	Update(userId, chapterId int, updateInput core.UpdateChapterInput) error
	Delete(userId, chapterId int) error
}

type Image interface {
	Create(image core.ChapterImage) (int, error)
	GetByChapterId(chapterId int) ([]core.ChapterImage, error)
}

type Service struct {
	Authorization
	Manga
	Chapter
	Image
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo.Authorization),
		Manga:         NewMangaService(repo.Manga),
		Chapter:       NewChapterService(repo.Chapter),
		Image:         NewImageService(repo.Image),
	}
}
