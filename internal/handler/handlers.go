package handler

import (
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"lang/internal/service"
	"net/http"
)

type Handler struct {
	service *service.Service
	logger  *zap.Logger
}

func NewHandler(service *service.Service, logger *zap.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) InitRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/home", h.homePage)
	auth := router.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/sign-in", h.signIn).Methods(http.MethodPost, http.MethodOptions)
	auth.HandleFunc("/sign-up", h.signUp).Methods(http.MethodPost, http.MethodOptions)

	admin := router.PathPrefix("/admin").Subrouter()
	admin.Use(h.userIdentity)
	chap := admin.PathPrefix("/grammar").Subrouter()
	chap.HandleFunc("/", h.createChapter).Methods(http.MethodPost, http.MethodOptions)
	chap.HandleFunc("/{chapter_id}", h.updateChapter).Methods(http.MethodPut, http.MethodOptions)
	chap.HandleFunc("/{chapter_id}", h.deleteChapter).Methods(http.MethodDelete, http.MethodOptions)

	article := chap.PathPrefix("/article").Subrouter()
	article.HandleFunc("/", h.createArticle).Methods(http.MethodPost, http.MethodOptions)

	user := router.PathPrefix("/").Subrouter()
	user.HandleFunc("/", h.getAllChapters).Methods(http.MethodGet, http.MethodOptions)
	user.HandleFunc("/{chapter_id}", h.getAllChapters).Methods(http.MethodGet, http.MethodOptions)

	return router
}

/*
	api := router.Group("/api")

	//api.GET("/home", h.homePage)

	{
		grammar := api.Group("/grammar")
		{
			grammar.GET("/", h.getAllChapters)
			grammar.GET("/:chapter_id", h.getChapterById)

			admin := grammar.Group("/admin", h.adminIdentity)
			{
				admin.POST("/", h.createChapter)
				admin.PUT("/:chapter_id", h.updateChapter)
				admin.DELETE("/:chapter_id", h.deleteChapter)
			}

			//topics := grammar.Group("/topics")
			//{
			//	topics.GET("/", h.getAll)
			//	topics.GET("/:topic_id", h.getTopicById)
			//}

		}
		//client := api.Group("/user")
		//{
		//
		//}

	}
	//admin's handlers

	return route
}
*/
