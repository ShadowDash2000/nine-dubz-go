package sort

import (
	"context"
	"net/http"
	"slices"
)

func SetSortContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sortBy := r.URL.Query().Get("sort-by")
		sortVal := r.URL.Query().Get("sort")

		if !slices.Contains([]string{"asc", "desc"}, sortBy) {
			sortBy = "desc"
		}

		sort := &Sort{
			SortBy:  sortBy,
			SortVal: sortVal,
		}

		ctx := context.WithValue(r.Context(), "sort", sort)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
