package middleware

import (
	"net/http"

	"github.com/ezcnrmn/quiz_service/internal/app/utils/consts"
)

func AddCorsHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Access-Control-Allow-Origin", consts.FRONTEND_ADDRESS)
		writer.Header().Add("Access-Control-Allow-Headers", "Content-Type")

		next.ServeHTTP(writer, request)
	})
}
