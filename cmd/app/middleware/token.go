package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"
)

//ErrNoAuthentication for
var ErrNoAuthentication = errors.New("no authentication") 
var authenticationContextKey = &contextKey{"authentication context"}

type contextKey struct {
	name string
}

func (c * contextKey) String() string {
	return c.name
}

//IDFunc for
type IDFunc func(ctx context.Context, token string) (int64, error)

//Authenticate for
func Authenticate(IdFunc IDFunc) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			token := request.Header.Get("Authorization")
			log.Print(token)
			id, err := IdFunc(request.Context(), token)
			if err != nil {
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			ctx := context.WithValue(request.Context(), authenticationContextKey, id)
			request = request.WithContext(ctx)

			handler.ServeHTTP(writer, request)
		}) 
	}
}

//Authentication for
func Authentication(ctx context.Context) (int64, error) {
	if value, ok := ctx.Value(authenticationContextKey).(int64); ok {
		return value, nil
	}
	return 0, ErrNoAuthentication
}