package v1

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/web-rabis/httperrors"
)

// @Summary Dictionary Type Card Backgrounds
// @Description Returns the background photo for the given type card id
// @Accept json
// @Produce json
// @Tags Dictionary
// @Param id path int true "Type Card ID"
// @Success 200 {string} string "binary"
// @Failure 400 {object} httperrors.Response
// @Failure 401 {object} httperrors.Response
// @Failure 403 {object} httperrors.Response
// @Failure 404 {object} httperrors.Response
// @Failure 500 {object} httperrors.Response
// @Router /api/v1/dictionary/type-card/{id}/backgrounds [get]
func (res *Resource) backgrounds(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		_ = render.Render(w, r, httperrors.BadRequest(errors.New("id is required")))
		return
	}
	_id, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		_ = render.Render(w, r, httperrors.BadRequest(err))
		return
	}
	background, err := res.dictionaryMan.TypeCardById(r.Context(), _id)
	if err != nil {
		_ = render.Render(w, r, httperrors.Internal(err))
		return
	}
	render.JSON(w, r, background.Photo)
}
