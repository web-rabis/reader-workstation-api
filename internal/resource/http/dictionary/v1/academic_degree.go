package v1

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/web-rabis/reader-workstation-api/internal/model"

	"github.com/web-rabis/httperrors"
	readerModel "github.com/web-rabis/reader-client/model"
)

// @Summary Dictionary Academic Degree
// @Description Dictionary Academic Degree
// @Accept json
// @Produce json
// @Tags Dictionary
// @Success 200 {array} readerModel.Dictionary
// @Failure 400 {object} httperrors.Response
// @Failure 401 {object} httperrors.Response
// @Failure 403 {object} httperrors.Response
// @Failure 404 {object} httperrors.Response
// @Failure 500 {object} httperrors.Response
// @Router /api/v1/dictionary/academic-degree [get]
func (res *Resource) getAcademicDegree(w http.ResponseWriter, r *http.Request) {
	paging, err := readerModel.PagingParseFromHttp(r)
	if err != nil {
		_ = render.Render(w, r, httperrors.Internal(err))
		return
	}
	count, l, err := res.dictionaryMan.AcademicDegreeList(r.Context(), paging)
	if err != nil {
		println(err.Error())
		_ = render.Render(w, r, httperrors.Internal(err))
		return
	}
	render.JSON(w, r, model.ListResponse{Count: count, Result: l})
}
