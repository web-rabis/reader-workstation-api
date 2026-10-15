package v1

import (
	"net/http"

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
// @Router /sex [get]
func (res *Resource) getSex(w http.ResponseWriter, r *http.Request) {
	_, sexList, err := res.dictionaryMan.SexList(r.Context(), nil)
	if err != nil {
		_ = render.Render(w, r, httperrors.Internal(err))
		return
	}
	render.JSON(w, r, sexList)
}
