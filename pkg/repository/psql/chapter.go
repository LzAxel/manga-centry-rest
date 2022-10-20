package psql

import (
	"fmt"
	"mangacentry/internal/core"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type ChapterPostgres struct {
	db *sqlx.DB
}

func NewChapterPostgres(db *sqlx.DB) *ChapterPostgres {
	return &ChapterPostgres{db: db}
}

func (r *ChapterPostgres) GetById(chapterId int) (core.MangaChapter, error) {
	var chapter core.MangaChapter
	logrus.Debugf("getting chapter %s in postgres db", chapterId)
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1",
		chapterTable)
	err := r.db.Get(&chapter, query, chapterId)
	if err != nil {
		return core.MangaChapter{}, err
	}

	return chapter, nil
}

func (r *ChapterPostgres) GetByMangaId(mangaId int) ([]core.ChapterListElement, error) {
	chapters := make([]core.ChapterListElement, 2)
	logrus.Debugf("getting chapters from manga %s in postgres db", mangaId)
	query := fmt.Sprintf("SELECT id, preview_url, number FROM %s WHERE manga_id = $1 ORDER BY -number",
		chapterTable)
	err := r.db.Select(&chapters, query, mangaId)
	if err != nil {
		return []core.ChapterListElement{}, err
	}
	return chapters, nil
}

func (r *ChapterPostgres) GetByNumber(mangaId, chapterNumber int) (core.MangaChapter, error) {
	var chapter core.MangaChapter
	logrus.Debugf("getting chapter %s from manga %s in postgres db", chapterNumber, mangaId)
	query := fmt.Sprintf("SELECT * FROM %s WHERE manga_id = $1 AND number = $2",
		chapterTable)
	err := r.db.Get(&chapter, query, mangaId, chapterNumber)
	if err != nil {
		return core.MangaChapter{}, err
	}

	return chapter, nil
}

func (r *ChapterPostgres) Create(userId int, chapter core.MangaChapter) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (preview_url, uploader_id, manga_id, number) values ($1, $2, $3, $4) RETURNING id",
		chapterTable)
	row := r.db.QueryRow(query, chapter.Preview, chapter.UploaderId, chapter.MangaId, chapter.Number)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *ChapterPostgres) Update(userId, chapterId int, updateInput core.UpdateChapterInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if updateInput.Preview != nil {
		setValues = append(setValues, fmt.Sprintf("preview_url=$%d", argId))
		args = append(args, *updateInput.Preview)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d AND uploader_id = $%d",
		chapterTable, setQuery, argId, argId+1)
	args = append(args, chapterId, userId)
	logrus.Debugf("args: %s", args)
	logrus.Debugf("updateQuery: %s", query)

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *ChapterPostgres) Delete(userId, chapterId int) error {
	var id int
	logrus.Debugln("deleting chapter from postgres db")
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1 AND uploader_id=$2 RETURNING id",
		chapterTable)
	row := r.db.QueryRow(query, chapterId, userId)
	if err := row.Scan(&id); err != nil {
		return err
	}

	return nil
}
