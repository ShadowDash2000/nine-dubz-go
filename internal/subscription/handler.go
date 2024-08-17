package subscription

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
	"nine-dubz/internal/pagination"
	"nine-dubz/internal/response"
	"nine-dubz/internal/user"
	"strconv"
)

type Handler struct {
	SubscriptionUseCase *UseCase
	UserHandler         *user.Handler
}

func NewHandler(uc *UseCase, uh *user.Handler) *Handler {
	return &Handler{
		SubscriptionUseCase: uc,
		UserHandler:         uh,
	}
}

func (h *Handler) SubscribeHandler(w http.ResponseWriter, r *http.Request) {
	channelId, err := strconv.ParseUint(chi.URLParam(r, "channelId"), 10, 32)
	if err != nil {
		response.RenderError(w, r, http.StatusBadRequest, "Invalid channel id")
		return
	}
	userId := r.Context().Value("userId").(uint)

	err = h.SubscriptionUseCase.Subscribe(userId, uint(channelId))
	if err != nil {
		response.RenderError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	response.RenderSuccess(w, r, http.StatusOK, "")
}

func (h *Handler) UnsubscribeHandler(w http.ResponseWriter, r *http.Request) {
	channelId, err := strconv.ParseUint(chi.URLParam(r, "channelId"), 10, 32)
	if err != nil {
		response.RenderError(w, r, http.StatusBadRequest, "Invalid channel id")
		return
	}
	userId := r.Context().Value("userId").(uint)

	err = h.SubscriptionUseCase.Unsubscribe(userId, uint(channelId))
	if err != nil {
		response.RenderError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	response.RenderSuccess(w, r, http.StatusOK, "")
}

func (h *Handler) GetMultipleHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(uint)
	pagination := r.Context().Value("pagination").(*pagination.Pagination)

	subscriptions, err := h.SubscriptionUseCase.GetMultiple(userId, pagination)
	if err != nil {
		response.RenderError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	if len(subscriptions) > 0 {
		render.JSON(w, r, subscriptions)
	} else {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, make([]struct{}, 0))
	}
}
