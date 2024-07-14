package user

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"golang.org/x/net/context"
	"net/http"
	"nine-dubz/internal/helper"
	"nine-dubz/internal/token"
	"nine-dubz/pkg/language"
	"nine-dubz/pkg/tokenauthorize"
)

type Handler struct {
	UserUseCase    *UseCase
	TokenUseCase   *token.UseCase
	TokenAuthorize *tokenauthorize.TokenAuthorize
	Language       *language.Repository
}

func NewHandler(uc *UseCase, tuc *token.UseCase, ta *tokenauthorize.TokenAuthorize, lang *language.Repository) *Handler {
	return &Handler{
		UserUseCase:    uc,
		TokenUseCase:   tuc,
		TokenAuthorize: ta,
		Language:       lang,
	}
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	loginRequest := &LoginRequest{}
	if err := json.NewDecoder(r.Body).Decode(loginRequest); err != nil {
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

		render.JSON(w, r, NewLoginResponse(true))
	}
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	registrationRequest := &RegistrationRequest{}
	if err := json.NewDecoder(r.Body).Decode(registrationRequest); err != nil {
		return
	}

	registrationPayload := NewRegistrationRequest(registrationRequest)
	userId := h.UserUseCase.Register(registrationPayload)
	if userId > 0 {
		h.SendRegistrationEmail(r, registrationPayload)

		render.JSON(w, r, NewRegistrationResponse(true))
	}
}

func (h *Handler) SendRegistrationEmail(r *http.Request, user *User) {
	languageCode := h.Language.GetLanguageCode(r)

	subject, _ := h.Language.GetStringByCode(languageCode, "EMAIL_REGISTRATION_CONFIRMATION")
	link := fmt.Sprintf("%s/api/authorize/inner/confirm/?email=%s&hash=%s", "https://"+r.Host, user.Email, user.Hash)
	contentValues := map[string]string{"userName": user.Name, "link": link}
	content, _ := h.Language.GetFormattedStringByCode("EMAIL_REGISTRATION_CONFIRMATION_CONTENT", contentValues, languageCode)

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

	if ok := h.UserUseCase.ConfirmRegistration(email, hash); !ok {
		http.Error(w, "Can't confirm registration", http.StatusBadRequest)
		return
	}

	tokenCookie, err := h.TokenAuthorize.GetTokenCookie(email)
	if err != nil {
		return
	}

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

func (h *Handler) PermissionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Context().Value("token").(string)

		routePattern := chi.RouteContext(r.Context()).RoutePattern()
		method := r.Method

		userId, ok := h.UserUseCase.CheckUserPermission(tokenString, routePattern, method)
		if !ok {
			http.Error(w, "Permission denied", http.StatusForbidden)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "userId", userId)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
