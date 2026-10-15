package v1

import (
	"net/http"

	"github.com/go-chi/render"

	"github.com/web-rabis/httperrors"
	ordermodel "github.com/web-rabis/order-client/model"
)

// @Summary Reader Check-Out
// @Description Registers a reader's check-out visit from the reading hall
// @Accept json
// @Produce json
// @Tags Reader
// @Param id path int true "Reader Ticket Number"
// @Success 200 {string} string "OK"
// @Failure 400 {object} httperrors.Response
// @Failure 401 {object} httperrors.Response
// @Failure 403 {object} httperrors.Response
// @Failure 404 {object} httperrors.Response
// @Failure 500 {object} httperrors.Response
// @Router /api/v1/reader/{id}/visits/check-out [post]
func (res *Resource) visitsCheckOut(w http.ResponseWriter, r *http.Request) {
	reader, err := res.readerCheckInOut(r, "*********")
	if err != nil {
		_ = render.Render(w, r, httperrors.BadRequest(err))
		return
	}
	user, err := res.getUser(r)
	if err != nil {
		_ = render.Render(w, r, httperrors.BadRequest(err))
		return
	}
	orderFilter := ordermodel.OrderFilters{
		TicketNumber: reader.TicketNumber,
		States:       []string{ordermodel.OrderStateInHands},
	}
	ordersCount, _ := res.orderMan.Count(r.Context(), &orderFilter)
	err = res.controlMan.CheckOut(r.Context(), reader.TicketNumber, "*********", user, ordersCount)
	if err != nil {
		_ = render.Render(w, r, httperrors.Internal(err))
		return
	}
}
