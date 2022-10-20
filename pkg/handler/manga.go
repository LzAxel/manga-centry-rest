package handler

import (
	"mangacentry/internal/core"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func (h *Handler) getMangaById(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		NewErrorResponse(ctx, http.StatusBadRequest, "invalid manga id")
		return
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, "invalid manga id")
		return
	}

	manga, err := h.services.Manga.GetById(idInt)
	if err != nil {
		NewErrorResponse(ctx, http.StatusBadGateway, err.Error())
		return
	}
	manga.Chapters, err = h.services.GetByMangaId(manga.Id)
	if err != nil {
		NewErrorResponse(ctx, http.StatusBadGateway, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, manga)
}

func (h *Handler) createManga(ctx *gin.Context) {
	logrus.Debugln("creating manga")

	var formInput core.CreateMangaInput
	var input core.Manga

	userId := h.getUserId(ctx)

	logrus.Debugf("userId: %d", userId)
	logrus.Debugln("checking input valid")
	if err := ctx.ShouldBindWith(&formInput, binding.FormMultipart); err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	logrus.Debugln("validating input")
	err := formInput.Validate()
	if err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	logrus.Debugln("binding input to model")
	if err := ctx.BindWith(&formInput, binding.FormMultipart); err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	fileName := formInput.Preview.Filename
	fileNameSplit := strings.Split(fileName, ".")
	fileFormattedName := "preview." + fileNameSplit[len(fileNameSplit)-1]

	input.Name = formInput.Name
	if formInput.AlternativeName != "" {
		input.AlternativeName = formInput.AlternativeName
	}
	input.UploaderId = userId

	mangaId, err := h.services.Manga.Create(userId, input)
	if err != nil {
		NewErrorResponse(ctx, http.StatusBadGateway, err.Error())
		return
	}
	mangaFolderPath := filepath.Join(viper.GetString("uploadPath"), "manga", strconv.Itoa(mangaId))
	logrus.Debugf("uplaoding file %s", filepath.Join(mangaFolderPath, fileFormattedName))

	err = os.MkdirAll(mangaFolderPath, 0777)
	if err != nil {
		logrus.Debugln(err.Error())
		NewErrorResponse(ctx, http.StatusBadGateway, "failed uploding file")
		return
	}

	err = ctx.SaveUploadedFile(formInput.Preview, filepath.Join(mangaFolderPath, fileFormattedName))
	if err != nil {
		logrus.Debugln(err.Error())
		h.services.Manga.Delete(userId, mangaId)
		NewErrorResponse(ctx, http.StatusBadRequest, "invalid file")
		return
	}

	filePath := filepath.Join("static", "manga", strconv.Itoa(mangaId), fileFormattedName)
	err = h.services.Manga.Update(userId, mangaId, core.UpdateMangaInput{
		Preview: &filePath,
	})
	if err != nil {
		logrus.Debugln(err.Error())
		h.services.Manga.Delete(userId, mangaId)
		NewErrorResponse(ctx, http.StatusBadRequest, "invalid file")
		return
	}

	logrus.Debugln(input)
	ctx.JSON(http.StatusCreated, map[string]interface{}{
		"id": mangaId,
	})
}

func (h *Handler) updateManga(ctx *gin.Context) {
	logrus.Debugln("updating manga")
	var input core.UpdateMangaInput
	id := ctx.Param("id")
	if id == "" {
		NewErrorResponse(ctx, http.StatusBadRequest, "invalid manga id")
		return
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, "invalid manga id")
		return
	}
	err = ctx.BindJSON(&input)
	if err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, "invalid update body")
		return
	}
	userId := h.getUserId(ctx)

	err = h.services.Manga.Update(userId, idInt, input)
	if err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
	})
	logrus.Debug("manga updated!")
}

func (h *Handler) deleteManga(ctx *gin.Context) {
	logrus.Debugln("deleting manga")
	id := ctx.Params.ByName("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, "invalid manga id")
		return
	}

	userId := h.getUserId(ctx)

	err = h.services.Manga.Delete(userId, idInt)
	if err != nil {
		NewErrorResponse(ctx, http.StatusBadGateway, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
	})
}
