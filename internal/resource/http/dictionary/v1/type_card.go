package v1

import (
	"net/http"

	"github.com/go-chi/render"

	"github.com/web-rabis/httperrors"
	readerModel "github.com/web-rabis/reader-client/model"
)

// @Summary Dictionary Type Card
// @Description Dictionary Type Card
// @Accept json
// @Produce json
// @Tags Dictionary
// @Success 200 {array} readerModel.TypeCard
// @Failure 400 {object} httperrors.Response
// @Failure 401 {object} httperrors.Response
// @Failure 403 {object} httperrors.Response
// @Failure 404 {object} httperrors.Response
// @Failure 500 {object} httperrors.Response
// @Router /api/v1/dictionary/type-card [get]
func (res *Resource) getTypeCard(w http.ResponseWriter, r *http.Request) {
	var list []*readerModel.TypeCard
	_, list, err := res.dictionaryMan.TypeCardList(r.Context(), nil)
	if err != nil {
		_ = render.Render(w, r, httperrors.Internal(err))
		return
	}
	render.JSON(w, r, list)
}
