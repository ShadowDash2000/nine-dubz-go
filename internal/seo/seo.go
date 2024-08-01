package seo

import (
	"errors"
	"golang.org/x/net/html"
	"net/http"
	"net/url"
	"nine-dubz/internal/movie"
	"nine-dubz/pkg/htmlcrawler"
	"nine-dubz/pkg/language"
	"nine-dubz/pkg/seometa"
	"regexp"
)

type UseCase struct {
	MovieUseCase *movie.UseCase
}

func New(movuc *movie.UseCase) *UseCase {
	return &UseCase{
		MovieUseCase: movuc,
	}
}

func (uc *UseCase) SetSeo(r *http.Request, document *html.Node) error {
	head := htmlcrawler.CrawlByTag("head", document)
	if head == nil {
		return errors.New("seo: head not found")
	}

	path := r.URL.Path
	seo, err := uc.GetSeo(path, r)
	if err != nil {
		return err
	}

	seometa.Set(head, seo)

	return nil
}

func (uc *UseCase) GetSeo(path string, r *http.Request) (map[string]string, error) {
	path, err := url.JoinPath(path, "/")
	if err != nil {
		return nil, err
	}

	movieDetailRegexp, err := regexp.Compile(`^/movie/(.*)/$`)
	if err != nil {
		return nil, err
	}
	matches := movieDetailRegexp.FindStringSubmatch(path)
	if matches != nil {
		movieCode := matches[1]
		return uc.MovieUseCase.GetMovieDetailSeo(movieCode, r)
	}

	languageCode := language.GetLanguageCode(r)
	siteName, err := language.GetMessage("SITE_NAME", languageCode)
	if err != nil {
		return nil, err
	}
	description, err := language.GetMessage("SEO_DEFAULT_DESCRIPTION", languageCode)
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"title":       siteName,
		"description": description,
	}, nil
}
