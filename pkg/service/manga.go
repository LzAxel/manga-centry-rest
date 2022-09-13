package service

import (
	"mangacentry/models"
	"mangacentry/pkg/repository"
)

type MangaService struct {
	repo repository.Manga
}

func NewMangaService(repo repository.Manga) *MangaService {
	return &MangaService{repo: repo}
}

func (r *MangaService) GetById(mangaId int) (models.Manga, error) {
	return r.repo.GetById(mangaId)
}

func (r *MangaService) Create(userId int, manga models.Manga) (int, error) {
	manga.UploaderId = userId

	return r.repo.Create(userId, manga)
}

func (r *MangaService) Update(userId, mangaId int, updateInput models.UpdateMangaInput) error {
	if err := updateInput.Validate(); err != nil {
		return err
	}
	return r.repo.Update(userId, mangaId, updateInput)
}

func (r *MangaService) Delete(userId, mangaId int) error {
	return r.repo.Delete(userId, mangaId)
}
