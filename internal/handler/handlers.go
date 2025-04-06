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

	router = router.PathPrefix("/langSimple.com").Subrouter()

	auth := router.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/sign-in", h.signIn).Methods(http.MethodPost, http.MethodOptions)
	auth.HandleFunc("/sign-up", h.signUp).Methods(http.MethodPost, http.MethodOptions)

	languageCode := router.PathPrefix("/{lang_code}").Subrouter()

	grammarUser := languageCode.PathPrefix("/grammar").Subrouter()
	grammarUser.HandleFunc("/", h.getAllChapters).Methods(http.MethodGet, http.MethodOptions)
	grammarUser.HandleFunc("/{chapter_id}", h.getChapterById).Methods(http.MethodGet, http.MethodOptions)

	articleUser := grammarUser.PathPrefix("/{chapter_id}/{title_chapter}").Subrouter()
	articleUser.HandleFunc("/", h.getAllChapters).Methods(http.MethodGet, http.MethodOptions)
	articleUser.HandleFunc("/{article_id}", h.getArticleById).Methods(http.MethodGet, http.MethodOptions)

	_ = languageCode.PathPrefix("/{course_id}").Subrouter()

	//-------------------------------------------------Client Pouters------------------------------------------------------------------------

	client := languageCode.PathPrefix("/grammar").Subrouter()
	client.Use(h.userIdentity)

	grammarClient := client.PathPrefix("/grammar").Subrouter()
	grammarClient.HandleFunc("/", h.getAllChapters).Methods(http.MethodGet, http.MethodOptions)
	grammarClient.HandleFunc("/{chapter_id}", h.getChapterById).Methods(http.MethodGet, http.MethodOptions)

	articleClient := grammarUser.PathPrefix("/{chapter_id}/{title_chapter}").Subrouter()
	articleClient.HandleFunc("/", h.getAllChapters).Methods(http.MethodGet, http.MethodOptions)
	articleClient.HandleFunc("/{article_id}", h.getArticleById).Methods(http.MethodGet, http.MethodOptions)

	_ = languageCode.PathPrefix("/{course_id}").Subrouter()

	//--------------------------------------------Moderator Routers-------------------------------------------------------------------------------------

	moderator := languageCode.PathPrefix("/admin").Subrouter()
	moderator.Use(h.adminIdentity)

	grammarModer := moderator.PathPrefix("/grammar").Subrouter()
	grammarModer.HandleFunc("/", h.createChapter).Methods(http.MethodPost, http.MethodOptions)
	grammarModer.HandleFunc("/", h.getAllChapters).Methods(http.MethodGet, http.MethodOptions)
	grammarModer.HandleFunc("/{chapter_id}", h.getChapterById).Methods(http.MethodGet, http.MethodOptions)
	grammarModer.HandleFunc("/{chapter_id}", h.updateChapter).Methods(http.MethodPut, http.MethodOptions)
	grammarModer.HandleFunc("/{chapter_id}", h.deleteChapter).Methods(http.MethodDelete, http.MethodOptions)

	articleModer := grammarModer.PathPrefix("/{chapter_id}/{title_chapter}").Subrouter()
	articleModer.HandleFunc("/", h.createArticle).Methods(http.MethodPost, http.MethodOptions)
	articleModer.HandleFunc("/", h.getAllChapters).Methods(http.MethodGet, http.MethodOptions)
	articleModer.HandleFunc("/{article_id}", h.getArticleById).Methods(http.MethodGet, http.MethodOptions)
	articleModer.HandleFunc("/{chapter_id}", h.updateChapter).Methods(http.MethodPut, http.MethodOptions)
	articleModer.HandleFunc("/{chapter_id}", h.deleteArticle).Methods(http.MethodDelete, http.MethodOptions)

	_ = moderator.PathPrefix("/{course_id}").Subrouter()

	//----------------------------------------ADMIN ROUTERS!!!!-----------------------------------------------------------------
	admin := languageCode.PathPrefix("/admin").Subrouter()
	admin.Use(h.adminIdentity)

	grammarAdmin := admin.PathPrefix("/grammar").Subrouter()
	grammarAdmin.HandleFunc("/", h.createChapter).Methods(http.MethodPost, http.MethodOptions)
	grammarAdmin.HandleFunc("/", h.getAllChapters).Methods(http.MethodGet, http.MethodOptions)
	grammarAdmin.HandleFunc("/{chapter_id}", h.getChapterById).Methods(http.MethodGet, http.MethodOptions)
	grammarAdmin.HandleFunc("/{chapter_id}", h.updateChapter).Methods(http.MethodPut, http.MethodOptions)
	grammarAdmin.HandleFunc("/{chapter_id}", h.deleteChapter).Methods(http.MethodDelete, http.MethodOptions)

	articleAdmin := grammarAdmin.PathPrefix("/{chapter_id}/{title_chapter}").Subrouter()
	articleAdmin.HandleFunc("/", h.createArticle).Methods(http.MethodPost, http.MethodOptions)
	articleAdmin.HandleFunc("/", h.getAllChapters).Methods(http.MethodGet, http.MethodOptions)
	articleAdmin.HandleFunc("/{article_id}", h.getArticleById).Methods(http.MethodGet, http.MethodOptions)
	articleAdmin.HandleFunc("/{chapter_id}", h.updateChapter).Methods(http.MethodPut, http.MethodOptions)
	articleAdmin.HandleFunc("/{chapter_id}", h.deleteArticle).Methods(http.MethodDelete, http.MethodOptions)

	_ = admin.PathPrefix("/{course_id}").Subrouter()

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
