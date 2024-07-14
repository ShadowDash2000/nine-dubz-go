package middleware

import (
	"context"
	"net/http"
	"nine-dubz/model"
	"strconv"
)

func PaginationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil || limit <= 0 {
			limit = -1
		}

		offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil || offset < 0 || limit <= 0 {
			offset = -1
		}

		pagination := &model.Pagination{
			Limit:  limit,
			Offset: offset,
		}

		ctx := context.WithValue(r.Context(), "pagination", pagination)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
