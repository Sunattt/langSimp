package handler

import (
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"lang/internal/service"
	"net/http"
)

type Handler struct {
	service *service.Ser
	logger  *zap.Logger
}

func NewHandler(service *service.Ser, logger *zap.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) InitRoutes() *mux.Router {
	router := mux.NewRouter()

	router = router.PathPrefix("/langSimple.com").Subrouter()

	auth := router.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/sign-in", h.signIn).Methods(http.MethodPost, http.MethodOptions)
	auth.HandleFunc("/sign-up", h.signUp).Methods(http.MethodPost, http.MethodOptions)

	languageCode := router.PathPrefix("/{lang_code}").Subrouter()
	languageCode.Use(h.langMiddleware)

	languageCode.HandleFunc("/home", h.homePage).Methods(http.MethodGet, http.MethodOptions)

	grammarUser := languageCode.PathPrefix("/grammar").Subrouter()
	grammarUser.HandleFunc("/", h.getAllChapters).Methods(http.MethodGet, http.MethodOptions)
	grammarUser.HandleFunc("/{chapter_id}", h.getChapterById).Methods(http.MethodGet, http.MethodOptions)

	articleUser := grammarUser.PathPrefix("/{chapter_id}/articles").Subrouter()
	articleUser.HandleFunc("/", h.getAllArticles).Methods(http.MethodGet, http.MethodOptions)
	articleUser.HandleFunc("/{article_id}", h.getArticleById).Methods(http.MethodGet, http.MethodOptions)

	courseUser := languageCode.PathPrefix("/course/{article_id}/{level_content}").Subrouter()
	courseUser.HandleFunc("/", h.getCourse).Methods(http.MethodGet, http.MethodOptions)

	//-------------------------------------------------Client Pouters------------------------------------------------------------------------

	client := languageCode.PathPrefix("/grammar").Subrouter()
	client.Use(h.userIdentity)

	grammarClient := client.PathPrefix("/grammar").Subrouter()
	grammarClient.HandleFunc("/", h.getAllChapters).Methods(http.MethodGet, http.MethodOptions)
	grammarClient.HandleFunc("/{chapter_id}", h.getChapterById).Methods(http.MethodGet, http.MethodOptions)

	articleClient := grammarUser.PathPrefix("/{chapter_id}/articles").Subrouter()
	articleClient.HandleFunc("/", h.getAllArticles).Methods(http.MethodGet, http.MethodOptions)
	articleClient.HandleFunc("/{article_id}", h.getArticleById).Methods(http.MethodGet, http.MethodOptions)

	courseClient := client.PathPrefix("/course/{article_id}/{level_content}").Subrouter()
	courseClient.HandleFunc("/", h.getCourse).Methods(http.MethodGet, http.MethodOptions)

	//--------------------------------------------Moderator Routers-------------------------------------------------------------------------------------

	moderator := languageCode.PathPrefix("/moder").Subrouter()
	moderator.Use(h.moderIdentity)

	grammarModer := moderator.PathPrefix("/grammar").Subrouter()
	grammarModer.HandleFunc("/", h.createChapter).Methods(http.MethodPost, http.MethodOptions)
	grammarModer.HandleFunc("/", h.getAllChapters).Methods(http.MethodGet, http.MethodOptions)
	grammarModer.HandleFunc("/{chapter_id}", h.getChapterById).Methods(http.MethodGet, http.MethodOptions)
	grammarModer.HandleFunc("/{chapter_id}", h.updateChapter).Methods(http.MethodPut, http.MethodOptions)
	grammarModer.HandleFunc("/{chapter_id}", h.deleteChapter).Methods(http.MethodDelete, http.MethodOptions)

	articleModer := grammarModer.PathPrefix("/{chapter_id}/articles").Subrouter()
	articleModer.HandleFunc("/", h.createArticle).Methods(http.MethodPost, http.MethodOptions)
	articleModer.HandleFunc("/", h.getAllArticles).Methods(http.MethodGet, http.MethodOptions)
	articleModer.HandleFunc("/{article_id}", h.getArticleById).Methods(http.MethodGet, http.MethodOptions)
	articleModer.HandleFunc("/{chapter_id}", h.updateArticle).Methods(http.MethodPut, http.MethodOptions)
	articleModer.HandleFunc("/{chapter_id}", h.deleteArticle).Methods(http.MethodDelete, http.MethodOptions)

	courseModer := moderator.PathPrefix("/course/{article_id}/{level_content}").Subrouter()
	courseModer.HandleFunc("/", h.createCourse).Methods(http.MethodPost, http.MethodOptions)
	courseModer.HandleFunc("/", h.getCourse).Methods(http.MethodGet, http.MethodOptions)
	courseModer.HandleFunc("/{content_id}", h.updateCourse).Methods(http.MethodPut, http.MethodOptions)
	courseModer.HandleFunc("/{content_id}", h.deleteArticle).Methods(http.MethodDelete, http.MethodOptions)

	//----------------------------------------ADMIN ROUTERS!!!!-----------------------------------------------------------------
	admin := languageCode.PathPrefix("/admin").Subrouter()
	admin.Use(h.adminIdentity)

	grammarAdmin := admin.PathPrefix("/grammar").Subrouter()
	grammarAdmin.HandleFunc("/", h.createChapter).Methods(http.MethodPost, http.MethodOptions)
	grammarAdmin.HandleFunc("/", h.getAllChapters).Methods(http.MethodGet, http.MethodOptions)
	grammarAdmin.HandleFunc("/{chapter_id}", h.getChapterById).Methods(http.MethodGet, http.MethodOptions)
	grammarAdmin.HandleFunc("/{chapter_id}", h.updateChapter).Methods(http.MethodPut, http.MethodOptions)
	grammarAdmin.HandleFunc("/{chapter_id}", h.deleteChapter).Methods(http.MethodDelete, http.MethodOptions)

	articleAdmin := grammarAdmin.PathPrefix("/{chapter_id}/articles").Subrouter()
	articleAdmin.HandleFunc("/", h.createArticle).Methods(http.MethodPost, http.MethodOptions)
	articleAdmin.HandleFunc("/", h.getAllArticles).Methods(http.MethodGet, http.MethodOptions)
	articleAdmin.HandleFunc("/{article_id}", h.getArticleById).Methods(http.MethodGet, http.MethodOptions)
	articleAdmin.HandleFunc("/{article_id}", h.updateArticle).Methods(http.MethodPut, http.MethodOptions)
	articleAdmin.HandleFunc("/{article_id}", h.deleteArticle).Methods(http.MethodDelete, http.MethodOptions)

	course := admin.PathPrefix("/course/{article_id}/{level_content}").Subrouter()
	course.HandleFunc("/", h.createCourse).Methods(http.MethodPost, http.MethodOptions)
	course.HandleFunc("/", h.getCourse).Methods(http.MethodGet, http.MethodOptions)
	course.HandleFunc("/{content_id}", h.updateCourse).Methods(http.MethodPut, http.MethodOptions)
	course.HandleFunc("/{content_id}", h.deleteArticle).Methods(http.MethodDelete, http.MethodOptions)

	return router
}
