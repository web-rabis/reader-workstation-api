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

// @Summary Approve Reader
// @Description Approves a reader registration request and issues a ticket number and barcode
// @Accept json
// @Produce json
// @Tags Reader
// @Param id path int true "Reader Ticket Number"
// @Success 200 {object} dto.ApproveReaderResponse
// @Failure 400 {object} httperrors.Response
// @Failure 401 {object} httperrors.Response
// @Failure 403 {object} httperrors.Response
// @Failure 404 {object} httperrors.Response
// @Failure 500 {object} httperrors.Response
// @Router /api/v1/reader/{id}/approve [put]
func (res *Resource) approve(w http.ResponseWriter, r *http.Request) {
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
	result, err := res.readerMan.Approve(r.Context(), _id)
	if err != nil {
		_ = render.Render(w, r, httperrors.Internal(err))
		return
	}
	response := dto.ApproveReaderResponse{
		OldTicketNumber: result.OldTicketNumber,
		Barcode:         result.Barcode,
		TicketNumber:    result.TicketNumber,
		Status:          result.Status,
	}
	if result.ExpirationDate != nil {
		response.ExpirationDate = *result.ExpirationDate
	}
	render.JSON(w, r, response)
}
