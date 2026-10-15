package v1

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	
	"github.com/web-rabis/httperrors"
	"github.com/web-rabis/reader-workstation-api/internal/resource/http/reader/v1/dto"
)

// @Summary Reader By Ticket Number
// @Description Returns a reader by their ticket number
// @Accept json
// @Produce json
// @Tags Reader
// @Param id path int true "Reader Ticket Number"
// @Success 200 {object} dto.ByIdResponse
// @Failure 400 {object} httperrors.Response
// @Failure 401 {object} httperrors.Response
// @Failure 403 {object} httperrors.Response
// @Failure 404 {object} httperrors.Response
// @Failure 500 {object} httperrors.Response
// @Router /api/v1/reader/{id} [get]
func (res *Resource) byId(w http.ResponseWriter, r *http.Request) {
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
	reader, err := res.readerMan.ByTicketNumber(r.Context(), _id)
	if err != nil {
		_ = render.Render(w, r, httperrors.BadRequest(err))
		return
	}
	render.JSON(w, r, dto.ByIdResponse{
		Result: reader,
	})
}
