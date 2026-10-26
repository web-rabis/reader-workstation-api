package v1

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/web-rabis/httperrors"
	"github.com/web-rabis/reader-workstation-api/internal/model"

	readerModel "github.com/web-rabis/reader-client/model"
)

// @Summary Dictionary Education
// @Description Dictionary Education
// @Accept json
// @Produce json
// @Tags Dictionary
// @Success 200 {array} readerModel.Dictionary
// @Failure 400 {object} httperrors.Response
// @Failure 401 {object} httperrors.Response
// @Failure 403 {object} httperrors.Response
// @Failure 404 {object} httperrors.Response
// @Failure 500 {object} httperrors.Response
// @Router /api/v1/dictionary/education [get]
func (res *Resource) getEducation(w http.ResponseWriter, r *http.Request) {
	paging, err := readerModel.PagingParseFromHttp(r)
	if err != nil {
		_ = render.Render(w, r, httperrors.Internal(err))
		return
	}
	count, l, err := res.dictionaryMan.EducationList(r.Context(), paging)
	if err != nil {
		_ = render.Render(w, r, httperrors.Internal(err))
		return
	}
	render.JSON(w, r, model.ListResponse{Count: count, Result: l})
}
