package handler

import (
	"mangacentry/internal/core"
	"mangacentry/pkg/utils"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func (h *Handler) getChapterById(ctx *gin.Context) {
	chapterId := ctx.Param("id")
	if chapterId == "" {
		NewErrorResponse(ctx, http.StatusBadRequest, "invalid chapter number")
		return
	}
	chapterIdInt, err := strconv.Atoi(chapterId)
	if err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, "invalid chapter number")
		return
	}

	chapter, err := h.services.Chapter.GetById(chapterIdInt)
	if err != nil {
		NewErrorResponse(ctx, http.StatusNoContent, "")
		return
	}
	chapter.Images, err = h.services.Image.GetByChapterId(chapter.Id)
	if err != nil {
		NewErrorResponse(ctx, http.StatusBadGateway, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, chapter)
}

func (h *Handler) createChapter(ctx *gin.Context) {
	var input core.CreateChapterInput

	userId := h.getUserId(ctx)

	err := ctx.ShouldBindWith(&input, binding.FormMultipart)
	if err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err = ctx.BindWith(&input, binding.FormMultipart)
	if err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	logrus.Debugf("creating chapter")
	chapterId, err := h.services.Chapter.Create(userId, input)
	if err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	logrus.Debugf("uploading archive")
	chapterPath := filepath.Join(viper.GetString("uploadPath"), "manga", strconv.Itoa(input.MangaId), strconv.Itoa(input.Number))
	err = os.MkdirAll(chapterPath, 0777)
	if err != nil {
		logrus.Debugln(err.Error())
		NewErrorResponse(ctx, http.StatusBadGateway, "failed to upload file")
		return
	}
	err = ctx.SaveUploadedFile(input.File, filepath.Join(chapterPath, "archive.zip"))
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, "failed to upload file")
		h.services.Chapter.Delete(userId, chapterId)
		return
	}

	chapterUrl := path.Join("static", "manga", strconv.Itoa(input.MangaId), strconv.Itoa(input.Number))
	archivePath := filepath.Join(chapterPath, "archive.zip")
	uploadedImages, err := utils.UnzipImages(archivePath, chapterPath)
	if err != nil {
		h.services.Chapter.Delete(userId, chapterId)
		err = os.RemoveAll(chapterPath)
		if err != nil {
			panic(err)
		}
		return
	}
	logrus.Debugf("uploaded images slice: %s", uploadedImages)

	for _, fileName := range uploadedImages {
		image := core.ChapterImage{
			Url:       path.Join(chapterUrl, fileName),
			ChapterId: chapterId,
		}
		h.services.Image.Create(image)
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"id": chapterId,
	})
}

func (h *Handler) updateChapter(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "s")
}

func (h *Handler) deleteChapter(ctx *gin.Context) {
	userId := h.getUserId(ctx)

	id := ctx.Param("id")
	if id == "" {
		NewErrorResponse(ctx, http.StatusBadRequest, "invalid chapter id")
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, "invalid chapter id")
		return
	}

	err = h.services.Chapter.Delete(userId, idInt)
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
	})
}
