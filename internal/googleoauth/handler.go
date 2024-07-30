package googleoauth

import (
	"github.com/go-chi/render"
	"net/http"
	"nine-dubz/internal/response"
	"nine-dubz/internal/token"
	"nine-dubz/internal/user"
	"nine-dubz/pkg/tokenauthorize"
	"strings"
)

type Handler struct {
	GoogleOAuthUseCase *UseCase
	UserHandler        *user.Handler
	TokenUseCase       *token.UseCase
	TokenAuthorize     *tokenauthorize.TokenAuthorize
}

func NewHandler(uc *UseCase, uh *user.Handler, tuc *token.UseCase, ta *tokenauthorize.TokenAuthorize) *Handler {
	return &Handler{
		GoogleOAuthUseCase: uc,
		UserHandler:        uh,
		TokenUseCase:       tuc,
		TokenAuthorize:     ta,
	}
}

func (h *Handler) GetConsentPageUrlHandler(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, struct {
		Url string `json:"url"`
	}{h.GoogleOAuthUseCase.GetConsentPageUrl()})
}

func (h *Handler) Authorize(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	if code == "" || state == "" {
		return
	}

	googleUser, err := h.GoogleOAuthUseCase.Authorize(code, state)
	if err != nil {
		return
	}

	// Try to log in
	loginRequest := &UserLoginRequest{
		Name:  strings.Split(googleUser.Email, "@")[0],
		Email: googleUser.Email,
	}

	userId := h.GoogleOAuthUseCase.Login(loginRequest)
	if userId > 0 {
		tokenCookie, err := h.TokenAuthorize.GetTokenCookie(loginRequest.Email)
		if err != nil {
			return
		}

		h.TokenUseCase.Add(userId, tokenCookie.Value)

		http.SetCookie(w, tokenCookie)

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Try to register
	registrationRequest := &UserRegistrationRequest{
		Name:       strings.Split(googleUser.Email, "@")[0],
		Email:      googleUser.Email,
		PictureUrl: googleUser.Picture,
	}

	userId, err = h.GoogleOAuthUseCase.Register(registrationRequest)
	if err != nil {
		response.RenderError(w, r, http.StatusInternalServerError, "")
		return
	} else if userId > 0 {
		tokenCookie, err := h.TokenAuthorize.GetTokenCookie(loginRequest.Email)
		if err != nil {
			return
		}

		h.TokenUseCase.Add(userId, tokenCookie.Value)

		http.SetCookie(w, tokenCookie)

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	response.RenderError(w, r, http.StatusInternalServerError, "")
}
