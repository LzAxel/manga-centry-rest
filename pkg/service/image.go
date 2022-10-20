package service

import (
	"mangacentry/internal/core"
	"mangacentry/pkg/repository"
)

type ImageService struct {
	repo repository.Image
}

func NewImageService(repo repository.Image) *ImageService {
	return &ImageService{repo: repo}
}

func (r *ImageService) Create(image core.ChapterImage) (int, error) {
	return r.repo.Create(image)
}

func (r *ImageService) GetByChapterId(chapterId int) ([]core.ChapterImage, error) {
	return r.repo.GetByChapterId(chapterId)
}
