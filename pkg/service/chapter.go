package service

import (
	"mangacentry/internal/core"
	"mangacentry/pkg/repository"
	"os"
	"strconv"

	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ChapterService struct {
	repo repository.Chapter
}

func NewChapterService(repo repository.Chapter) *ChapterService {
	return &ChapterService{repo: repo}
}

func (r *ChapterService) GetById(chapterId int) (core.MangaChapter, error) {
	manga, err := r.repo.GetById(chapterId)
	if err != nil {
		return manga, err
	}

	return manga, nil
}

func (r *ChapterService) GetByMangaId(mangaId int) ([]core.ChapterListElement, error) {
	return r.repo.GetByMangaId(mangaId)
}

func (r *ChapterService) GetByNumber(mangaId, chapterNumber int) (core.MangaChapter, error) {
	return r.repo.GetByNumber(mangaId, chapterNumber)
}

func (r *ChapterService) Create(userId int, input core.CreateChapterInput) (int, error) {
	var chapter core.MangaChapter

	chapter.UploaderId = userId
	chapter.MangaId = input.MangaId
	chapter.Number = input.Number

	chapterId, err := r.repo.Create(userId, chapter)
	if err != nil {
		return 0, nil
	}

	return chapterId, nil
}

func (r *ChapterService) Update(userId, chapterId int, updateInput core.UpdateChapterInput) error {
	return r.repo.Update(userId, chapterId, updateInput)
}

func (r *ChapterService) Delete(userId, chapterId int) error {
	chapter, err := r.repo.GetById(chapterId)
	if err != nil {
		return err
	}
	chapterPath := filepath.Join(viper.GetString("uploadPath"), "manga", strconv.Itoa(chapter.MangaId), strconv.Itoa(chapter.Number))
	logrus.Debugf("deleting chapter folder: %s", chapterPath)
	err = os.RemoveAll(chapterPath)
	if err != nil {
		return err
	}

	return r.repo.Delete(userId, chapterId)
}
