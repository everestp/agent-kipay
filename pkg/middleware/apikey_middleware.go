// pkg/middleware/api_key_middleware.go

package middleware

import (
    "context"
    "net/http"
    "strings"

    apikeyservice "github.com/everest/bheri/modules/api_key/service"
    "github.com/everest/bheri/pkg/utils"
)

type contextKey string

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
                utils.Error(
                    w,
                    http.StatusUnauthorized,
                    "missing API key",
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

            if err != nil || key == nil {
                utils.Error(
                    w,
                    http.StatusUnauthorized,
                    "invalid API key",
                )
                return
            }

            // =====================================================
            // CHECK STATUS
            // =====================================================

            if key.Status != "active" {
                utils.Error(
                    w,
                    http.StatusUnauthorized,
                    "API key is not active",
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
