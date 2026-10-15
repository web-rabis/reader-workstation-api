package v1

import (
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/web-rabis/httperrors"
	readerModel "github.com/web-rabis/reader-client/model"
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
func (res *Resource) visitsCheckIn(w http.ResponseWriter, r *http.Request) {
	reader, err := res.readerCheckInOut(r, "КР")
	if err != nil {
		_ = render.Render(w, r, httperrors.BadRequest(err))
		return
	}
	user, err := res.getUser(r)
	if err != nil {
		_ = render.Render(w, r, httperrors.BadRequest(err))
		return
	}
	err = res.controlMan.CheckIn(r.Context(), reader.TicketNumber, "КР", user)
	if err != nil {
		_ = render.Render(w, r, httperrors.Internal(err))
		return
	}
}
func (res *Resource) readerCheckInOut(r *http.Request, department string) (*readerModel.Reader, error) {
	reader, err := res.getReader(r)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	reader.Department = department
	reader.LastVisitDate = &now
	err = res.readerMan.Update(r.Context(), reader, nil)
	if err != nil {
		return nil, err
	}
	return reader, nil
}
