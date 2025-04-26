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

	courseUser := grammarUser.PathPrefix("/course/{article_id}/{level_content}").Subrouter()
	courseUser.HandleFunc("/", h.getCourse).Methods(http.MethodGet, http.MethodOptions)

	//-------------------------------------------------Client Pouters------------------------------------------------------------------------

	client := languageCode.PathPrefix("/grammar").Subrouter()
	client.Use(h.userIdentity)

	grammarClient := client.PathPrefix("/grammar").Subrouter()
	grammarClient.HandleFunc("/", h.getAllChapters).Methods(http.MethodGet, http.MethodOptions)
	grammarClient.HandleFunc("/{chapter_id}", h.getChapterById).Methods(http.MethodGet, http.MethodOptions)

	articleClient := grammarClient.PathPrefix("/{chapter_id}/articles").Subrouter()
	articleClient.HandleFunc("/", h.getAllArticles).Methods(http.MethodGet, http.MethodOptions)
	articleClient.HandleFunc("/{article_id}", h.getArticleById).Methods(http.MethodGet, http.MethodOptions)

	courseClient := grammarClient.PathPrefix("/course/{article_id}/{level_content}").Subrouter()
	courseClient.HandleFunc("/", h.getCourse).Methods(http.MethodGet, http.MethodOptions)

	quizRouter := grammarClient.PathPrefix("/{content_id}/quizzes").Subrouter()
	quizRouter.HandleFunc("/", h.getExercises).Methods(http.MethodGet)
	quizRouter.HandleFunc("/check", h.checkAnswers).Methods(http.MethodPost)

	profileCl := languageCode.PathPrefix("/profile").Subrouter()

	savedCl := profileCl.PathPrefix("/saved").Subrouter()

	savedCl.HandleFunc("/chapters", h.getSavedChapter).Methods(http.MethodGet, http.MethodOptions)
	savedCl.HandleFunc("/articles", h.getSavedArticle).Methods(http.MethodGet, http.MethodOptions)
	savedCl.HandleFunc("/words", h.getSavedWord).Methods(http.MethodGet, http.MethodOptions)
	savedCl.HandleFunc("/{}", h.removeSavedArticle).Methods(http.MethodDelete, http.MethodOptions)
	savedCl.HandleFunc("/{}", h.removeSavedChapter).Methods(http.MethodDelete, http.MethodOptions)
	savedCl.HandleFunc("/{}", h.removeSavedWord).Methods(http.MethodDelete, http.MethodOptions)
	savedCl.HandleFunc("/{}", h.saveChapter).Methods(http.MethodPost, http.MethodOptions)
	savedCl.HandleFunc("/{}", h.saveArticle).Methods(http.MethodPost, http.MethodOptions)

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

	courseModer := grammarModer.PathPrefix("/course/{article_id}/{level_content}").Subrouter()
	courseModer.HandleFunc("/", h.createCourse).Methods(http.MethodPost, http.MethodOptions)
	courseModer.HandleFunc("/", h.getCourse).Methods(http.MethodGet, http.MethodOptions)
	courseModer.HandleFunc("/{content_id}", h.updateCourse).Methods(http.MethodPut, http.MethodOptions)
	courseModer.HandleFunc("/{content_id}", h.deleteCourse).Methods(http.MethodDelete, http.MethodOptions)

	quizRouterM := grammarModer.PathPrefix("/{content_id}/quizzes").Subrouter()
	quizRouterM.HandleFunc("/", h.createExercise).Methods(http.MethodPost)
	quizRouterM.HandleFunc("/", h.getExercises).Methods(http.MethodGet)
	quizRouterM.HandleFunc("/check", h.checkAnswers).Methods(http.MethodPost)
	quizRouterM.HandleFunc("/{exercise_id}", h.updateExercise).Methods(http.MethodPut)
	quizRouterM.HandleFunc("/{exercise_id}", h.deleteExercise).Methods(http.MethodDelete)

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

	course := grammarAdmin.PathPrefix("/course/{article_id}/{level_content}").Subrouter()
	course.HandleFunc("/", h.createCourse).Methods(http.MethodPost, http.MethodOptions)
	course.HandleFunc("/", h.getCourse).Methods(http.MethodGet, http.MethodOptions)
	course.HandleFunc("/{content_id}", h.updateCourse).Methods(http.MethodPut, http.MethodOptions)
	course.HandleFunc("/{content_id}", h.deleteCourse).Methods(http.MethodDelete, http.MethodOptions)

	quizRouterA := grammarAdmin.PathPrefix("/{content_id}/quizzes").Subrouter()
	quizRouterA.HandleFunc("/", h.createExercise).Methods(http.MethodPost)
	quizRouterA.HandleFunc("/", h.getExercises).Methods(http.MethodGet)
	quizRouterA.HandleFunc("/check", h.checkAnswers).Methods(http.MethodPost)
	quizRouterA.HandleFunc("/{exercise_id}", h.updateExercise).Methods(http.MethodPut)
	quizRouterA.HandleFunc("/{exercise_id}", h.deleteExercise).Methods(http.MethodDelete)

	return router
}
