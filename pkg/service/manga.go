package service

import (
	"mangacentry/internal/core"
	"mangacentry/pkg/repository"
	"os"
	"path/filepath"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type MangaService struct {
	repo repository.Manga
}

func NewMangaService(repo repository.Manga) *MangaService {
	return &MangaService{repo: repo}
}

func (r *MangaService) GetById(mangaId int) (core.Manga, error) {
	return r.repo.GetById(mangaId)
}

func (r *MangaService) Create(userId int, manga core.Manga) (int, error) {
	manga.UploaderId = userId

	return r.repo.Create(userId, manga)
}

func (r *MangaService) Update(userId, mangaId int, updateInput core.UpdateMangaInput) error {
	if err := updateInput.Validate(); err != nil {
		return err
	}
	return r.repo.Update(userId, mangaId, updateInput)
}

func (r *MangaService) Delete(userId, mangaId int) error {
	err := r.repo.Delete(userId, mangaId)
	if err != nil {
		return err
	}
	mangaPath := filepath.Join(viper.GetString("uploadPath"), "manga", strconv.Itoa(mangaId))
	logrus.Debugf("deleting manga folder: %s", mangaPath)
	err = os.RemoveAll(mangaPath)
	if err != nil {
		return err
	}

	return err
}
