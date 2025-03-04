package endpoints

import (
	"log/slog"
	"net/http"
	"orbit-app/internal/services"
	"orbit-app/pkg"
)

func GetSnippetByID(snippetService services.SnippetService, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := pkg.DecodeBase62(r.PathValue("id"))
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte("invalid id"))
			return
		}
		resp, err := snippetService.GetByID(r.Context(), id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("not found"))
			logger.Error("error while fetching snippet by id", "error", err.Error())
			return
		}
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(resp.Content))
	}
}
