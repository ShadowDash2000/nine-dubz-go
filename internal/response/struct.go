package response

import (
	"github.com/go-chi/render"
	"log"
	"net/http"
	"nine-dubz/pkg/language"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

func RenderError(w http.ResponseWriter, r *http.Request, status int, message string) {
	messageLang, err := language.GetMessage(message, language.GetLanguageCode(r))
	if err != nil {
		messageLang = message
		log.Println(err)
	}

	render.Status(r, status)
	render.JSON(w, r, &Response{
		Status:  "error",
		Message: messageLang,
	})
}

func RenderSuccess(w http.ResponseWriter, r *http.Request, status int, message string) {
	messageLang, err := language.GetMessage(message, language.GetLanguageCode(r))
	if err != nil {
		messageLang = message
		log.Println(err)
	}

	render.Status(r, status)
	render.JSON(w, r, &Response{
		Status:  "success",
		Message: messageLang,
	})
}
