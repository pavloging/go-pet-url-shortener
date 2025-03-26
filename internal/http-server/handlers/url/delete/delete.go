package delete

import (
	"errors"
	"log/slog"
	"net/http"

	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/go-chi/render"
)

type Response struct {
	resp.Response
	CountDeleted int64 `json:"countDeleted"`
}

type DeleteURL interface {
	DeleteURL(alias string) (int64, error)
}

func New(log *slog.Logger, deleteURL DeleteURL) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.url.delete.New"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Error("empty alias")
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}

		countDeleted, err := deleteURL.DeleteURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", "alias", alias)
			render.JSON(w, r, resp.Error("internal error"))
			return
		}
		if err != nil {
			log.Error("failed to get url", "alias", alias)
			render.JSON(w, r, resp.Error("failed to get url"))
			return
		}

		log.Info("deleted url", slog.String("alias", alias), slog.Int64("count_deleted", countDeleted))

		responseOK(w, r, countDeleted)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, countDeleted int64) {
	render.JSON(w, r, Response{
		Response:     resp.OK(),
		CountDeleted: countDeleted,
	})
}
