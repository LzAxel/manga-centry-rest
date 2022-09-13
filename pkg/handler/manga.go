package handler

import (
	"mangacentry/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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

	ctx.JSON(http.StatusOK, manga)
}

func (h *Handler) createManga(ctx *gin.Context) {
	logrus.Debugln("creating manga")
	var input models.Manga
	userId, err := h.getUserId(ctx)
	if err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	logrus.Debugln("binding input to model")
	if err := ctx.BindJSON(&input); err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	mangaId, err := h.services.Manga.Create(userId, input)

	if err != nil {
		NewErrorResponse(ctx, http.StatusBadGateway, err.Error())
		return
	}
	logrus.Debugln(input)
	ctx.JSON(http.StatusCreated, map[string]interface{}{
		"id": mangaId,
	})
}

type updateMangaInput struct {
	Name    *string
	Preview *string
}

func (h *Handler) updateManga(ctx *gin.Context) {
	logrus.Debugln("updating manga")
	var input models.UpdateMangaInput
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
	userId, err := h.getUserId(ctx)
	if err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

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

	userId, err := h.getUserId(ctx)
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, "can't get user id")
		return
	}

	err = h.services.Manga.Delete(userId, idInt)
	if err != nil {
		NewErrorResponse(ctx, http.StatusBadGateway, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
	})
}
