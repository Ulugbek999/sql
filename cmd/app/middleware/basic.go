package middleware

import (
	"log"
	"net/http"


)


//Basic for
func Basic(s *security.Service) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

			user, pass, _ := request.BasicAuth()
            log.Print("hhhhhhhhhhh", user, pass)
            // value := s.Auth(user, pass)
			// if value != true {
			// 	writer.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			// 	http.Error(writer, "Unauthorized.", http.StatusUnauthorized)
			// 	return
			// }
			handler.ServeHTTP(writer, request)
		})
	}
}
