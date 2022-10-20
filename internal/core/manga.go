package core

import (
	"errors"
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/sirupsen/logrus"
)

type Manga struct {
	Id              int                  `json:"id" db:"id"`
	Name            string               `json:"name" binding:"required" db:"name"`
	AlternativeName string               `json:"alternativeName" db:"alternative_name"`
	Views           int                  `json:"views" db:"views"`
	Preview         string               `json:"previewUrl" binding:"required" db:"preview_url"`
	UploaderId      int                  `json:"uploaderId" db:"uploader_id"`
	Chapters        []ChapterListElement `json:"chapters"`
}

type CreateMangaInput struct {
	Preview         *multipart.FileHeader `json:"file" form:"file" binding:"required"`
	AlternativeName string                `json:"alternativeName" form:"alternativeName"`
	Name            string                `json:"name" form:"name" binding:"required"`
}

func (i CreateMangaInput) Validate() error {
	logrus.Debugln("validating create manga input")
	extansions := "png jpg webp"

	splitFileName := strings.Split(i.Preview.Filename, ".")
	fileExtension := splitFileName[len(splitFileName)-1]

	if !strings.Contains(extansions, fileExtension) {
		return fmt.Errorf("invalid file extension: %s", fileExtension)
	}

	return nil
}

type UpdateMangaInput struct {
	Name    *string `json:"name"`
	Preview *string `json:"previewUrl"`
}

func (i UpdateMangaInput) Validate() error {
	if i.Name == nil && i.Preview == nil {
		return errors.New("update input is empty")
	}

	return nil
}

type ChapterListElement struct {
	Id      int    `json:"id" db:"id"`
	Preview string `json:"previewUrl" binding:"required" db:"preview_url"`
	Number  int    `json:"number" binding:"required" db:"number"`
}

type MangaChapter struct {
	Id         int            `json:"id" db:"id"`
	Preview    string         `json:"previewUrl" binding:"required" db:"preview_url"`
	MangaId    int            `json:"mangaId" binding:"required" db:"manga_id"`
	Number     int            `json:"number" binding:"required" db:"number"`
	UploaderId int            `json:"uploaderId" db:"uploader_id"`
	Images     []ChapterImage `json:"images"`
}

type CreateChapterInput struct {
	File    *multipart.FileHeader `json:"file" form:"file" binding:"required"`
	MangaId int                   `json:"mangaId" form:"mangaId" binding:"required"`
	Number  int                   `json:"number" form:"number" binding:"required"`
}

func (i CreateChapterInput) Validate() error {
	logrus.Debugln("validating create chapter input")

	splitFileName := strings.Split(i.File.Filename, ".")
	fileExtension := splitFileName[len(splitFileName)-1]

	if fileExtension != "zip" {
		return fmt.Errorf("support only zip archives")
	}

	return nil
}

type UpdateChapterInput struct {
	Preview *string `json:"previewUrl"`
	Number  *int    `json:"number"`
}

type ChapterImage struct {
	Id        int    `json:"id" db:"id"`
	Url       string `json:"url" binding:"required" db:"url"`
	ChapterId int    `json:"-" db:"chapter_id"`
}

func (i ChapterImage) Validate() error {

	return nil
}
