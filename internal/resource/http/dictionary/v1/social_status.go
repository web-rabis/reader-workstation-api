package v1

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/web-rabis/httperrors"

	readerModel "github.com/web-rabis/reader-client/model"
)

// @Summary Dictionary Social Status
// @Description Dictionary Social Status
// @Accept json
// @Produce json
// @Tags Dictionary
// @Success 200 {array} readerModel.SocialStatus
// @Failure 400 {object} httperrors.Response
// @Failure 401 {object} httperrors.Response
// @Failure 403 {object} httperrors.Response
// @Failure 404 {object} httperrors.Response
// @Failure 500 {object} httperrors.Response
// @Router /api/v1/dictionary/social-status [get]
func (res *Resource) getSocialStatus(w http.ResponseWriter, r *http.Request) {
	paging, err := readerModel.PagingParseFromHttp(r)
	if err != nil {
		_ = render.Render(w, r, httperrors.Internal(err))
		return
	}
	_, l, err := res.dictionaryMan.SocialStatusList(r.Context(), paging)
	if err != nil {
		_ = render.Render(w, r, httperrors.Internal(err))
		return
	}
	render.JSON(w, r, l)
}
