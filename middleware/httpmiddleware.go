package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/afif0808/sagara-test/contexts"
	"github.com/afif0808/sagara-test/internal/domain"
	"github.com/afif0808/sagara-test/wrapper"
)

type usecase interface {
	Authenticate(ctx context.Context, token string) (domain.User, error)
}

type HTTPMiddleware struct {
	uc usecase
}

func NewHTTPMiddleware(uc usecase) HTTPMiddleware {
	return HTTPMiddleware{
		uc: uc,
	}
}

func (mw *HTTPMiddleware) BearerAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if !strings.HasPrefix(token, "Bearer ") {
			wrapper.NewHTTPResponse(http.StatusUnauthorized, "invalid token").JSON(w)
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")
		ctx := r.Context()
		user, err := mw.uc.Authenticate(ctx, token)
		if err != nil {
			wrapper.NewHTTPResponse(http.StatusUnauthorized, "invalid token", err).JSON(w)
			return
		}
		
		ctx = context.WithValue(ctx, contexts.UserContextKey, user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
