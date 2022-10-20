package handler

import (
	"mangacentry/pkg/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services service.Service) *Handler {
	return &Handler{services: &services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	api := router.Group("/api", h.userIdentity)
	{
		manga := api.Group("/manga")
		{
			manga.GET("/:id", h.getMangaById)
			manga.POST("/", h.createManga)
			manga.PATCH("/:id", h.updateManga)
			manga.DELETE("/:id", h.deleteManga)

		}
		chapter := api.Group("/chapter")
		{
			chapter.GET("/:id", h.getChapterById)
			chapter.POST("/", h.createChapter)
			chapter.PATCH("/:id", h.updateChapter)
			chapter.DELETE("/:id", h.deleteChapter)
		}
	}

	return router
}
