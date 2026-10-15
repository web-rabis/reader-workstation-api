package v1

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/web-rabis/httperrors"
)

func (res *Resource) getTypeCard(w http.ResponseWriter, r *http.Request) {
	_, list, err := res.dictionaryMan.TypeCardList(r.Context(), nil)
	if err != nil {
		_ = render.Render(w, r, httperrors.Internal(err))
		return
	}
	render.JSON(w, r, list)
}
