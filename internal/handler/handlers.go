package handler

import (
	"lang/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h Handler) InitRoute() *gin.Engine {
	route := gin.New()

	auth := route.Group("/auth")
	{
		auth.POST("/sign-up", h.sighUp)
		auth.POST("/sign-in", h.sighIn)
	}

	api := route.Group("api", h.userIdentity)
	{
		en := api.Group("/en")

		homePage := en.Group("")
		homePage.GET("/#")

		grammar := en.Group("/grammar")
		{
			grammar.POST("/", h.createChapter)
			grammar.GET("/", h.getAllChapters)
			grammar.GET("/:id", h.getChapterById)
			grammar.PUT("/:id", h.updateChapter)
			grammar.DELETE("/:id", h.deleteChapter)
		}
		// 			chapter := grammer.Group("/:id/:topic")
		// 			{
		// 				chapter.POST("/")
		// 				chapter.GET("/")
		// 				chapter.GET("/:id")
		// 				chapter.PUT("/:id")
		// 				chapter.DELETE("/:id")
		// 			}

		// 			article := grammer.Group("/course/:id/:name_topic")
		// 			{
		// 				article.POST("/")
		// 				article.GET("/")
		// 				article.GET("/:id")
		// 				article.PUT("/:id")
		// 				article.DELETE("/:id")

		// 			}
		// 		}

		// 		vocab := en.Group("/vocab")
		// 		{
		// 			vocab.POST("/")
		// 			vocab.GET("/")
		// 			vocab.GET("/:id")
		// 			vocab.PUT("/:id")
		// 			vocab.DELETE("/:id")
		// 		}

		// 		panel := en.Group("/dassboard")
		// 		{
		// 			profile := panel.Group("/profile")
		// 			{
		// 				profile.POST("/")
		// 				profile.GET("/")
		// 				profile.GET("/:id")
		// 				profile.PUT("/:id")
		// 				profile.DELETE("/:id")
		// 			}

		// 			marked := panel.Group("/marked")
		// 			{
		// 				marked.POST("/")
		// 				marked.GET("/")
		// 				marked.GET("/:id")
		// 				marked.PUT("/:id")
		// 				marked.DELETE("/:id")
		// 			}
		// 		}
	}
	return route
}
