package middleware

import (
	"log/slog"
	"net/http"
	"regexp"
	"slices"
)

var supportedTypes = []string{
	"text/html",
	"text/plain",
}

func ContentType(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if slices.Contains(supportedTypes, r.Header.Get("Content-Type")) {
				next.ServeHTTP(w, r)
				return
			}
			r.Header.Set("Content-Type", "text/html")
			re, err := regexp.Compile(`(?i)^(curl|wget|httpie)`)
			if err != nil {
				logger.Error("regexp compilation error", "error", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Internal Server Error"))
				return
			}
			if re.MatchString(r.UserAgent()) {
				r.Header.Set("Content-Type", "text/plain")
			}
			next.ServeHTTP(w, r)
		})
	}
}
