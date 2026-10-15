package v1

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/web-rabis/httperrors"
)

// @Summary Dictionary Sex
// @Description Dictionary Sex
// @Accept json
// @Produce json
// @Tags Dictionary
// @Success 200 {object} reader.Dictionary
// @Failure 400 {object} httperrors.Response
// @Failure 401 {object} httperrors.Response
// @Failure 403 {object} httperrors.Response
// @Failure 404 {object} httperrors.Response
// @Failure 500 {object} httperrors.Response
// @Router /admin/api/v1/auth/sign-in [get]
func (res *Resource) renewalRegistration(w http.ResponseWriter, r *http.Request) {
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
	err = res.readerMan.RenewalRegistration(r.Context(), _id, 5, 0, 0)
	if err != nil {
		_ = render.Render(w, r, httperrors.Internal(err))
		return
	}
	render.JSON(w, r, nil)
}
