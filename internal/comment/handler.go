package comment

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
	"nine-dubz/internal/pagination"
	"nine-dubz/internal/response"
	"nine-dubz/internal/sorting"
	"nine-dubz/internal/user"
	"strconv"
)

type Handler struct {
	CommentUseCase *UseCase
	UserHandler    *user.Handler
}

func NewHandler(cuc *UseCase, uh *user.Handler) *Handler {
	return &Handler{
		CommentUseCase: cuc,
		UserHandler:    uh,
	}
}

func (h *Handler) AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	movieCode := chi.URLParam(r, "movieCode")

	var commentId uint64
	var err error
	commentId = 0
	commentIdParam := chi.URLParam(r, "commentId")
	if commentIdParam != "" {
		commentId, err = strconv.ParseUint(commentIdParam, 10, 32)
		if err != nil {
			response.RenderError(w, r, http.StatusBadRequest, "Failed to parse comment id")
			return
		}
	}

	userId := r.Context().Value("userId").(uint)

	commentAddRequest := &AddRequest{}
	if err := json.NewDecoder(r.Body).Decode(commentAddRequest); err != nil {
		response.RenderError(w, r, http.StatusBadRequest, "")
		return
	}

	err = h.CommentUseCase.Add(userId, movieCode, commentAddRequest.Text, uint(commentId))
	if err != nil {
		response.RenderError(w, r, http.StatusBadRequest, "Can't add comment")
		return
	}

	response.RenderSuccess(w, r, http.StatusOK, "")
}

func (h *Handler) GetMultipleSubCommentsHandler(w http.ResponseWriter, r *http.Request) {
	movieCode := chi.URLParam(r, "movieCode")
	subPagination := r.Context().Value("pagination").(*pagination.Pagination)
	userId := r.Context().Value("userId")
	if userId == nil {
		userId = uint(0)
	}

	var commentId uint64
	var err error
	commentId = 0
	commentIdParam := chi.URLParam(r, "commentId")
	if commentIdParam != "" {
		commentId, err = strconv.ParseUint(commentIdParam, 10, 32)
		if err != nil {
			response.RenderError(w, r, http.StatusBadRequest, "Failed to parse comment id")
			return
		}
	}

	comments, err := h.CommentUseCase.GetMultipleSubComments(userId.(uint), movieCode, uint(commentId), subPagination)
	if err != nil {
		response.RenderError(w, r, http.StatusBadRequest, "Can't get comments")
		return
	}

	render.JSON(w, r, comments)
}

func (h *Handler) GetMultipleHandler(w http.ResponseWriter, r *http.Request) {
	movieCode := chi.URLParam(r, "movieCode")
	userId := r.Context().Value("userId").(*uint)
	pagination := r.Context().Value("pagination").(*pagination.Pagination)
	sorting := r.Context().Value("sorting").(*sorting.Sort)

	comments, err := h.CommentUseCase.GetMultiple(userId, movieCode, pagination, sorting)
	if err != nil {
		response.RenderError(w, r, http.StatusBadRequest, "Can't get comments")
		return
	}

	if comments != nil {
		render.JSON(w, r, comments)
	} else {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, make([]struct{}, 0))
	}
}

func (h *Handler) DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	commentId, err := strconv.ParseUint(chi.URLParam(r, "commentId"), 10, 32)
	if err != nil {
		response.RenderError(w, r, http.StatusBadRequest, "Invalid comment id")
		return
	}
	userId := r.Context().Value("userId").(uint)

	rowsAffected, err := h.CommentUseCase.Delete(uint(commentId), userId)
	if err != nil || rowsAffected == 0 {
		response.RenderError(w, r, http.StatusBadRequest, "Can't delete comment")
		return
	}

	response.RenderSuccess(w, r, http.StatusOK, "")
}
