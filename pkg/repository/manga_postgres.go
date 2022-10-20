package repository

import (
	"fmt"
	"mangacentry/internal/core"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type MangaPostgres struct {
	db *sqlx.DB
}

func NewMangaPostgres(db *sqlx.DB) *MangaPostgres {
	return &MangaPostgres{db: db}
}

func (r *MangaPostgres) GetById(mangaId int) (core.Manga, error) {
	var manga core.Manga
	logrus.Debugln("getting manga by id in postgres db")
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1",
		mangaTable)
	err := r.db.Get(&manga, query, mangaId)
	if err != nil {
		return core.Manga{}, err
	}

	return manga, nil
}

func (r *MangaPostgres) Create(userId int, manga core.Manga) (int, error) {
	var id int
	logrus.Debugln("inserting manga into postgres db")
	query := fmt.Sprintf("INSERT INTO %s (name, preview_url, uploader_id, alternative_name) values ($1, $2, $3, $4) RETURNING id",
		mangaTable)
	row := r.db.QueryRow(query, manga.Name, manga.Preview, manga.UploaderId, manga.AlternativeName)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *MangaPostgres) Update(userId, mangaId int, updateInput core.UpdateMangaInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if updateInput.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argId))
		args = append(args, *updateInput.Name)
		argId++
	}

	if updateInput.Preview != nil {
		setValues = append(setValues, fmt.Sprintf("preview_url=$%d", argId))
		args = append(args, *updateInput.Preview)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d AND uploader_id = $%d",
		mangaTable, setQuery, argId, argId+1)
	args = append(args, mangaId, userId)
	logrus.Debugf("args: %s", args)
	logrus.Debugf("updateQuery: %s", query)

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *MangaPostgres) Delete(userId, mangaId int) error {
	var id int
	logrus.Debugln("deleting manga from postgres db")
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1 and uploader_id=$2 RETURNING id",
		mangaTable)
	row := r.db.QueryRow(query, mangaId, userId)
	if err := row.Scan(&id); err != nil {
		return err
	}

	return nil
}
