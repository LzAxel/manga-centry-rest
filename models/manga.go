package models

import "errors"

type MangaImage struct {
	Id      int    `json:"id"`
	Url     string `json:"url"`
	MangaId int    `json:"mangaId"`
}

type Manga struct {
	Id         int    `json:"id" db:"id"`
	Name       string `json:"name" binding:"required" db:"name"`
	Views      int    `json:"views" db:"views"`
	Preview    string `json:"preview" binding:"required" db:"preview"`
	UploaderId int    `json:"uploaderId" db:"uploader_id"`
}

type UpdateMangaInput struct {
	Name    *string `json:"name"`
	Preview *string `json:"preview"`
}

func (i UpdateMangaInput) Validate() error {
	if i.Name == nil && i.Preview == nil {
		return errors.New("update input is empty")
	}

	return nil
}
