package seo

type GetSeoResponse struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func NewGetSeoResponse(seo map[string]string) *GetSeoResponse {
	seoResponse := &GetSeoResponse{}

	for key, val := range seo {
		switch key {
		case "title":
			seoResponse.Title = val
			break
		case "description":
			seoResponse.Description = val
			break
		}
	}

	return seoResponse
}
