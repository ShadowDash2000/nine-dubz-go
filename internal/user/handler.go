package user

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/render"
	"net/http"
	"nine-dubz/internal/helper"
	"nine-dubz/internal/response"
	"nine-dubz/internal/token"
	"nine-dubz/pkg/language"
	"nine-dubz/pkg/tokenauthorize"
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
	} else {
		response.RenderError(w, r, http.StatusBadRequest, "LOGIN_USER_NOT_FOUND")
	}
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
		return
	}

	registrationPayload := NewRegistrationRequest(registrationRequest)
	userId, err := h.UserUseCase.Register(registrationPayload)
	if err != nil {
		response.RenderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if userId > 0 {
		h.SendRegistrationEmail(r, registrationPayload)

		response.RenderSuccess(w, r, http.StatusOK, "")
	}
}

func (h *Handler) SendRegistrationEmail(r *http.Request, user *User) {
	languageCode := language.GetLanguageCode(r)

	subject, _ := language.GetMessage(languageCode, "EMAIL_REGISTRATION_CONFIRMATION")
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
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	render.JSON(w, r, NewShortResponse(user))
}

func (h *Handler) ConfirmRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	hash := r.URL.Query().Get("hash")
	if email == "" || hash == "" {
		http.Error(w, "Missing email or hash", http.StatusBadRequest)
		return
	}

	userId, ok := h.UserUseCase.ConfirmRegistration(email, hash)
	if !ok {
		http.Error(w, "Can't confirm registration", http.StatusBadRequest)
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
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("picture")
	if err != nil {
		http.Error(w, "No picture file", http.StatusBadRequest)
		return
	}

	userId := r.Context().Value("userId").(uint)
	if err = h.UserUseCase.UpdatePicture(userId, file, fileHeader); err != nil {
		http.Error(w, "Failed to update picture", http.StatusBadRequest)
		return
	}
}
