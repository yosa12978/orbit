package endpoints

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"orbit-app/internal/dto"
	"orbit-app/internal/services"
	"orbit-app/pkg"
)

func PutSnippet(snippetService services.SnippetService, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// curl
		data, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		req := dto.SnippetCreateRequest{Content: string(data)}
		resp, err := snippetService.Create(r.Context(), req)
		if err != nil {
			if errors.As(err, &pkg.ValidationError{}) {
				w.WriteHeader(http.StatusBadRequest)
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(err.Error()))
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
			logger.Error("error while creating snippet", "error", err.Error())
			return
		}
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%v", pkg.EncodeBase62(resp.ID))
	}
}
