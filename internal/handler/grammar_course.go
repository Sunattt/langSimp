package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"lang/pkg/models"
	"lang/pkg/utils"
	"net/http"
	"strconv"
)

func (h *Handler) createCourse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	level := vars["level_content"]

	articleID, err := strconv.Atoi(vars["article_id"])

	if err != nil {
		utils.NewResponseError(w, http.StatusInternalServerError, "Invalid article ID format", err, h.logger)
		return
	}

	levelId, err := h.service.GetLevelId(level)
	if err != nil {
		utils.NewResponseError(w, http.StatusInternalServerError, "error with level", err, h.logger)
		return
	}

	var input models.GrammarContent

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.NewResponseError(w, http.StatusBadRequest, "Invalid request body format", err, h.logger)
		return
	}

	input.LevelId = levelId
	input.ArticleId = articleID

	id, errResp := h.service.CoursePost.CreateContent(input)
	if errResp.Code != 0 {
		utils.NewResponseError(w, errResp.Code, errResp.Message, errResp.Error, h.logger)
		return
	}

	utils.ResponseServer(map[string]int{"content_id": id}, w, h.logger)
}

func (h *Handler) getCourse(w http.ResponseWriter, r *http.Request) {
	level := mux.Vars(r)["level_content"]

	levelId, err := h.service.GetLevelId(level)
	if err != nil {
		utils.NewResponseError(w, http.StatusInternalServerError, "Invalid level content ID format", err, h.logger)
		return
	}

	courseId, err := strconv.Atoi(mux.Vars(r)["article_id"])
	if err != nil {
		utils.NewResponseError(w, http.StatusBadRequest, "invalid id", err, h.logger)
		return
	}

	data, err := h.service.CoursePost.GetCourseById(courseId, levelId)
	if err != nil {
		utils.NewResponseError(w, http.StatusInternalServerError, "failed", err, h.logger)
		return
	}

	response := map[string]interface{}{
		"data": data,
	}

	utils.ResponseServer(response, w, h.logger)
}

func (h *Handler) updateCourse(w http.ResponseWriter, r *http.Request) {

	level := mux.Vars(r)["level_content"]

	levelId, err := h.service.GetLevelId(level)
	if err != nil {

		return
	}

	contentId, err := strconv.Atoi(mux.Vars(r)["content_id"])

	if err != nil {
		utils.NewResponseError(w, http.StatusBadRequest, "Invalid content ID format", err, h.logger)
		return
	}

	var input models.UpdateContentInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.NewResponseError(w, http.StatusBadRequest, "invalid dody update data", err, h.logger)
		return
	}

	err = h.service.CoursePost.UpdateContent(contentId, levelId, input)
	if err != nil {
		utils.NewResponseError(w, http.StatusInternalServerError, "failed to update article", err, h.logger)
		return
	}

	utils.ResponseServer(map[string]string{"status": "ok"}, w, h.logger)
}

func (h *Handler) deleteCourse(w http.ResponseWriter, r *http.Request) {

	contentId, err := strconv.Atoi(mux.Vars(r)["content_id"])
	if err != nil {
		utils.NewResponseError(w, http.StatusBadRequest, "Invalid content ID format", err, h.logger)
		return
	}
	level := mux.Vars(r)["level_content"]
	levelId, err := h.service.GetLevelId(level)
	if err != nil {
		utils.NewResponseError(w, http.StatusInternalServerError, "Failed to get level ID ", err, h.logger)
		return
	}

	err = h.service.CoursePost.DeleteContent(contentId, levelId)
	if err != nil {
		utils.NewResponseError(w, http.StatusInternalServerError, "failed to delete content", err, h.logger)
		return
	}

	utils.ResponseServer(map[string]string{"status": "ok"}, w, h.logger)
}

func (h *Handler) createExercise(w http.ResponseWriter, r *http.Request) {
	contentId, err := strconv.Atoi(mux.Vars(r)["content_id"])
	if err != nil {
		utils.NewResponseError(w, http.StatusBadRequest, "Invalid content ID format", err, h.logger)
		return
	}

	var input models.GrammarContentExercises
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.NewResponseError(w, http.StatusBadRequest, "invalid body required", err, h.logger)
		return
	}

	input.GrammarContentId = contentId

	id, err := h.service.CoursePost.CreateExercise(&input)
	if err != nil {
		utils.NewResponseError(w, http.StatusInternalServerError, "failed to create", err, h.logger)
		return
	}

	utils.ResponseServer(map[string]int{"exercise_id": id}, w, h.logger)
}

func (h *Handler) getExercises(w http.ResponseWriter, r *http.Request) {
	contentId, err := strconv.Atoi(mux.Vars(r)["content_id"])
	if err != nil {
		utils.NewResponseError(w, http.StatusBadRequest, "Invalid content ID format", err, h.logger)
		return
	}

	exercise, err := h.service.GetExerciseById(contentId)
	if err != nil {
		utils.NewResponseError(w, http.StatusInternalServerError, "failed with get exercise", err, h.logger)
		return
	}

	utils.ResponseServer(map[string]interface{}{
		"quizzes": exercise,
		"count":   len(exercise),
	}, w, h.logger)

}

func (h *Handler) checkAnswers(w http.ResponseWriter, r *http.Request) {
	var answers []models.UserAnswer
	if err := json.NewDecoder(r.Body).Decode(&answers); err != nil {
		utils.NewResponseError(w, http.StatusBadRequest, "Invalid request body", err, h.logger)
		return
	}

	responses, err := h.service.CheckAnswers(answers)
	if err != nil {
		utils.NewResponseError(w, http.StatusInternalServerError, "Failed to check answers", err, h.logger)
		return
	}

	utils.ResponseServer(map[string]interface{}{
		"results": responses,
	}, w, h.logger)
}

func (h *Handler) updateExercise(w http.ResponseWriter, r *http.Request) {

	contentId, err := strconv.Atoi(mux.Vars(r)["content_id"])
	if err != nil {
		utils.NewResponseError(w, http.StatusBadRequest, "Invalid content ID format", err, h.logger)
		return
	}

	var input models.UpdateGrammarExercise
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.NewResponseError(w, http.StatusBadRequest, "invalid body required", err, h.logger)
		return
	}

	err = h.service.UpdateExercise(contentId, input)
	if err != nil {
		utils.NewResponseError(w, http.StatusInternalServerError, "failed to update exercise ", err, h.logger)
		return
	}

	utils.ResponseServer(map[string]string{"status": "ok"}, w, h.logger)
}

func (h *Handler) deleteExercise(w http.ResponseWriter, r *http.Request) {

	contentId, err := strconv.Atoi(mux.Vars(r)["content_id"])
	if err != nil {
		utils.NewResponseError(w, http.StatusBadRequest, "Invalid content ID format", err, h.logger)
		return
	}

	level := mux.Vars(r)["level_content"]
	levelId, err := h.service.GetLevelId(level)
	if err != nil {
		utils.NewResponseError(w, http.StatusInternalServerError, "invalid level id ", err, h.logger)
		return
	}

	err = h.service.CoursePost.DeleteExercise(contentId, levelId)
	if err != nil {
		utils.NewResponseError(w, http.StatusInternalServerError, "failed to delete exercise", err, h.logger)
		return
	}

	utils.ResponseServer(map[string]string{"status": "ok"}, w, h.logger)
}

func (h *Handler) createComment(w http.ResponseWriter, r *http.Request) {
	contentId, err := strconv.Atoi(mux.Vars(r)["content_id"])
	if err != nil {
		utils.NewResponseError(w, http.StatusBadRequest, "Invalid content ID format", err, h.logger)
		return
	}

	userId := h.getUserId(w, r)
	if userId == 0 {
		return
	}

	var comment models.GrammarComment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		utils.NewResponseError(w, http.StatusBadRequest, "Invalid request body format", err, h.logger)
		return
	}

	commentId, errResp := h.service.CreateComment(contentId, userId, comment)
	if errResp.Code == 0 {
		utils.NewResponseError(w, errResp.Code, errResp.Message, errResp.Error, h.logger)
		return
	}

	utils.ResponseServer(map[string]int{"commentId": commentId}, w, h.logger)
}

func (h *Handler) getAllComment(w http.ResponseWriter, r *http.Request) {
	contentId, err := strconv.Atoi(mux.Vars(r)["content_id"])
	if err != nil {
		utils.NewResponseError(w, http.StatusBadRequest, "Failed to content ID", err, h.logger)
		return
	}

	comments, err := h.service.GetAllComments(contentId)
	if err != nil {
		utils.NewResponseError(w, http.StatusInternalServerError, "failed with get comment", err, h.logger)
		return
	}

	response := map[string]interface{}{
		"data": comments,
	}

	utils.ResponseServer(response, w, h.logger)
}

func (h *Handler) updateComment(w http.ResponseWriter, r *http.Request) {
	commentId, err := strconv.Atoi(mux.Vars(r)["comment_id"])
	if err != nil {
		utils.NewResponseError(w, http.StatusBadRequest, "Failed to comment ID", err, h.logger)
		return
	}

	userId := h.getUserId(w, r)
	if userId == 0 {
		return
	}

	commentExists, err := h.service.GetComments(commentId, userId)
	if err != nil {
		utils.NewResponseError(w, http.StatusNotFound, "failed with comment", err, h.logger)
		return
	}

	if !commentExists {
		utils.NewResponseError(w, http.StatusNotFound, "comment does not exist", err, h.logger)
		return
	}

	var input models.UpdateGrammarComment
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.NewResponseError(w, http.StatusBadRequest, "Invalid request body format", err, h.logger)
		return
	}

	err = h.service.UpdateComment(commentId, input)
	if err != nil {
		utils.NewResponseError(w, http.StatusInternalServerError, "Failed to update comment", err, h.logger)
		return
	}

	utils.ResponseServer(map[string]string{"status": "ok"}, w, h.logger)
}

func (h *Handler) deleteComment(w http.ResponseWriter, r *http.Request) {
	commentId, err := strconv.Atoi(mux.Vars(r)["comment_id"])
	if err != nil {
		utils.NewResponseError(w, http.StatusBadRequest, "Failed to get comment ID", err, h.logger)
		return
	}

	userId := h.getUserId(w, r)
	if userId == 0 {
		return
	}

	commentExists, err := h.service.GetComments(commentId, userId)
	if err != nil {
		utils.NewResponseError(w, http.StatusNotFound, "failed with comment", err, h.logger)
		return
	}

	if !commentExists {
		utils.NewResponseError(w, http.StatusNotFound, "comment does not exist", err, h.logger)
		return
	}
	admin, err := h.service.IsAdmin(userId)

	if err != nil {
		utils.NewResponseError(w, http.StatusInternalServerError, "", err, h.logger)
		return
	}

	moder, err := h.service.IsModerator(userId)
	if err != nil {
		utils.NewResponseError(w, http.StatusInternalServerError, "", err, h.logger)
		return
	}

	if admin == true || moder == true {
		err := h.service.DeleteComment(commentId)
		if err != nil {
			utils.NewResponseError(w, http.StatusInternalServerError, "failed with delete comment", err, h.logger)
			return
		}
	}

	err = h.service.DeleteComment(commentId)
	if err != nil {
		utils.NewResponseError(w, http.StatusInternalServerError, "failed with delete comment", err, h.logger)
		return
	}

	utils.ResponseServer(map[string]string{"status": "ok"}, w, h.logger)

}

func (h *Handler) likeComment(w http.ResponseWriter, r *http.Request) {
	commentId, err := strconv.Atoi(mux.Vars(r)["comment_id"])
	if err != nil {
		utils.NewResponseError(w, http.StatusBadRequest, "Invalid comment id: %s ", err, h.logger)
		return
	}

	userId := h.getUserId(w, r)
	if userId == 0 {
		return
	}

	err = h.service.LikeComment(userId, commentId)
	if err != nil {
		utils.NewResponseError(w, http.StatusInternalServerError, "", err, h.logger)
		return
	}

	utils.ResponseServer(map[string]string{"status": "ok"}, w, h.logger)
}

func (h *Handler) removeLikeComment(w http.ResponseWriter, r *http.Request) {
	commentId, err := strconv.Atoi(mux.Vars(r)["comment_id"])
	if err != nil {
		utils.NewResponseError(w, http.StatusBadRequest, "Invalid comment id", err, h.logger)
		return
	}

	userId := h.getUserId(w, r)
	if userId == 0 {
		return
	}

	err = h.service.RemoveLike(userId, commentId)
	if err != nil {
		utils.NewResponseError(w, http.StatusInternalServerError, "", err, h.logger)
		return
	}

	utils.ResponseServer(map[string]string{"status": "ok"}, w, h.logger)
}
