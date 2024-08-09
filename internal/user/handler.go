package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"mime/multipart"
	"net/http"
	"nine-dubz/internal/helper"
	"nine-dubz/internal/response"
	"nine-dubz/internal/token"
	"nine-dubz/pkg/language"
	"nine-dubz/pkg/tokenauthorize"
	"strconv"
)

type Handler struct {
	UserUseCase    *UseCase
	TokenUseCase   *token.UseCase
	TokenAuthorize *tokenauthorize.TokenAuthorize
}

func NewHandler(uc *UseCase, tuc *token.UseCase, ta *tokenauthorize.TokenAuthorize) *Handler {
	return &Handler{
		UserUseCase:    uc,
		TokenUseCase:   tuc,
		TokenAuthorize: ta,
	}
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	loginRequest := &LoginRequest{}
	if err := json.NewDecoder(r.Body).Decode(loginRequest); err != nil {
		response.RenderError(w, r, http.StatusBadRequest, "LOGIN_INVALID_FIELDS")
		return
	}

	loginPayload := NewLoginRequest(loginRequest)
	userId := h.UserUseCase.Login(loginPayload)
	if userId > 0 {
		tokenCookie, err := h.TokenAuthorize.GetTokenCookie(loginPayload.Email)
		if err != nil {
			return
		}

		h.TokenUseCase.Add(loginPayload.ID, tokenCookie.Value)

		http.SetCookie(w, tokenCookie)

		response.RenderSuccess(w, r, http.StatusOK, "")
		return
	}

	response.RenderError(w, r, http.StatusBadRequest, "LOGIN_USER_NOT_FOUND")
}

func (h *Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	tokenCookie, err := r.Cookie("token")
	if err != nil {
		response.RenderError(w, r, http.StatusUnauthorized, "Token cookie not found")
		return
	}
	userId := r.Context().Value("userId").(uint)

	err = h.TokenUseCase.Delete(userId, tokenCookie.Value)
	if err != nil {
		response.RenderError(w, r, http.StatusBadRequest, "Can't logout")
		return
	}

	http.SetCookie(w, h.TokenAuthorize.GetEmptyTokenCookie())
	response.RenderSuccess(w, r, http.StatusOK, "")
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	registrationRequest := &RegistrationRequest{}
	if err := json.NewDecoder(r.Body).Decode(registrationRequest); err != nil {
		response.RenderError(w, r, http.StatusBadRequest, "REGISTRATION_INVALID_FIELDS")
		return
	}

	registrationPayload := NewRegistrationRequest(registrationRequest)
	userId, err := h.UserUseCase.Register(registrationPayload)
	if err != nil {
		response.RenderError(w, r, http.StatusBadRequest, err.Error())
		return
	} else if userId > 0 {
		h.SendRegistrationEmail(r, registrationPayload)

		response.RenderSuccess(w, r, http.StatusOK, "")
		return
	}

	response.RenderError(w, r, http.StatusInternalServerError, "REGISTRATION_INTERNAL_ERROR")
}

func (h *Handler) SendRegistrationEmail(r *http.Request, user *User) {
	languageCode := language.GetLanguageCode(r)

	subject, _ := language.GetMessage("EMAIL_REGISTRATION_CONFIRMATION", languageCode)
	link := fmt.Sprintf("%s/api/authorize/inner/confirm/?email=%s&hash=%s", "https://"+r.Host, user.Email, user.Hash)
	contentValues := map[string]string{"userName": user.Name, "link": link}
	content, _ := language.GetFormattedMessage("EMAIL_REGISTRATION_CONFIRMATION_CONTENT", contentValues, languageCode)

	h.UserUseCase.SendRegistrationEmail(user.Email, subject, content)
}

func (h *Handler) CheckUserWithNameExistsHandler(w http.ResponseWriter, r *http.Request) {
	userName := r.URL.Query().Get("userName")
	if userName == "" {
		return
	}

	if ok := helper.ValidateUserName(userName); !ok {
		return
	}

	isUserExists := h.UserUseCase.CheckUserWithNameExists(userName)

	render.JSON(w, r, struct {
		IsUserExists bool `json:"isUserExists"`
	}{isUserExists})
}

func (h *Handler) GetUserShortHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(uint)

	user, err := h.UserUseCase.GetById(userId)
	if err != nil {
		response.RenderError(w, r, http.StatusNotFound, "User not found")
		return
	}

	render.JSON(w, r, NewShortResponse(user))
}

func (h *Handler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.ParseUint(chi.URLParam(r, "userId"), 10, 32)
	if err != nil {
		response.RenderError(w, r, http.StatusBadRequest, "Invalid user id")
		return
	}

	user, err := h.UserUseCase.GetById(uint(userId))
	if err != nil {
		response.RenderError(w, r, http.StatusNotFound, "User not found")
		return
	}

	render.JSON(w, r, NewGetPublicResponse(user))
}

func (h *Handler) ConfirmRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	hash := r.URL.Query().Get("hash")
	if email == "" || hash == "" {
		response.RenderError(w, r, http.StatusBadRequest, "Missing email or hash")
		return
	}

	userId, ok := h.UserUseCase.ConfirmRegistration(email, hash)
	if !ok {
		response.RenderError(w, r, http.StatusBadRequest, "Can't confirm registration")
		return
	}

	tokenCookie, err := h.TokenAuthorize.GetTokenCookie(email)
	if err != nil {
		return
	}

	h.TokenUseCase.Add(userId, tokenCookie.Value)

	http.SetCookie(w, tokenCookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) UpdatePictureHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(2 << 20); err != nil {
		if errors.Is(err, multipart.ErrMessageTooLarge) {
			response.RenderError(w, r, http.StatusBadRequest, "File is too large")
			return
		}

		response.RenderError(w, r, http.StatusBadRequest, "Failed to parse form data")
		return
	}

	file, fileHeader, err := r.FormFile("picture")
	if err != nil {
		response.RenderError(w, r, http.StatusBadRequest, "No picture file")
		return
	}

	userId := r.Context().Value("userId").(uint)
	if err = h.UserUseCase.UpdatePicture(userId, file, fileHeader); err != nil {
		response.RenderError(w, r, http.StatusBadRequest, "Failed to update picture")
		return
	}
}

func (h *Handler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(uint)

	userUpdateRequest := &UpdateRequest{}
	if err := json.NewDecoder(r.Body).Decode(userUpdateRequest); err != nil {
		response.RenderError(w, r, http.StatusBadRequest, "USER_UPDATE_INVALID_FIELDS")
		return
	}
	userUpdateRequest.ID = userId

	if err := h.UserUseCase.Update(userUpdateRequest); err != nil {
		response.RenderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.RenderSuccess(w, r, http.StatusOK, "")
}
