package v1

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/web-rabis/httperrors"
	ordermodel "github.com/web-rabis/order-client/model"
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
