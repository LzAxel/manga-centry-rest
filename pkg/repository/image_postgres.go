package repository

import (
	"fmt"
	"mangacentry/internal/core"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type ImagePostgres struct {
	db *sqlx.DB
}

func NewImagePostgres(db *sqlx.DB) *ImagePostgres {
	return &ImagePostgres{db: db}
}

func (r *ImagePostgres) Create(image core.ChapterImage) (int, error) {
	var id int

	logrus.Debugln("inserting image into postgres db")
	query := fmt.Sprintf("INSERT INTO %s (chapter_id, url) values ($1, $2) RETURNING id",
		chapterImageTable)
	err := r.db.Get(&id, query, image.ChapterId, image.Url)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *ImagePostgres) GetByChapterId(chapterId int) ([]core.ChapterImage, error) {
	var images []core.ChapterImage
	logrus.Debugf("getting images from chapter %s in postgres db", chapterId)
	query := fmt.Sprintf("SELECT * FROM %s WHERE chapter_id = $1",
		chapterImageTable)
	err := r.db.Select(&images, query, chapterId)
	if err != nil {
		return images, err
	}

	return images, nil
}
