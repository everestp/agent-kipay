package middleware

import (
	"context"
	"net/http"
	"strings"

	apikeyservice "github.com/everest/bheri/modules/api_key/service"
)



const (
	APIKeyIDContextKey contextKey = "api_key_id"
	UserIDContextKey   contextKey = "user_id"
)

func APIKey(
	apiKeyService apikeyservice.APIKeyService,
) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {

			// =====================================================
			// READ API KEY
			// =====================================================

			apiKey := strings.TrimSpace(
				r.Header.Get("X-API-Key"),
			)

			if apiKey == "" {
				http.Error(
					w,
					"missing API key",
					http.StatusUnauthorized,
				)
				return
			}

			// =====================================================
			// VALIDATE API KEY
			// =====================================================

			key, err := apiKeyService.Validate(
				r.Context(),
				apiKey,
			)

			if err != nil {
				http.Error(
					w,
					"invalid API key",
					http.StatusUnauthorized,
				)
				return
			}

			if key == nil {
				http.Error(
					w,
					"invalid API key",
					http.StatusUnauthorized,
				)
				return
			}

			// =====================================================
			// CHECK STATUS
			// =====================================================

			if key.Status != "active" {
				http.Error(
					w,
					"API key is not active",
					http.StatusUnauthorized,
				)
				return
			}

			// =====================================================
			// CONTEXT
			// =====================================================

			ctx := context.WithValue(
				r.Context(),
				APIKeyIDContextKey,
				key.ID,
			)

			ctx = context.WithValue(
				ctx,
				UserIDContextKey,
				key.UserID,
			)

			// =====================================================
			// NEXT
			// =====================================================

			next.ServeHTTP(
				w,
				r.WithContext(ctx),
			)
		})
	}
}
