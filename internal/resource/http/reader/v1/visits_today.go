package v1

import (
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/web-rabis/httperrors"
	readerModel "github.com/web-rabis/reader-client/model"
	"github.com/web-rabis/reader-workstation-api/internal/resource/http/reader/v1/dto"
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
func (res *Resource) visitsToday(w http.ResponseWriter, r *http.Request) {
	exitDateIsNull := true
	start, end := DayRangeHalfOpen(time.Now())
	//start = start.AddDate(0, 0, -30)
	filter := readerModel.ControlFilter{
		VisitDateLow:   &start,
		VisitDateHigh:  &end,
		ExitDateIsNull: &exitDateIsNull,
	}
	paging, err := readerModel.PagingParseFromHttp(r)
	if err != nil {
		_ = render.Render(w, r, httperrors.BadRequest(err))
		return
	}
	count, list, err := res.controlMan.List(r.Context(), &filter, paging)
	if err != nil {
		_ = render.Render(w, r, httperrors.BadRequest(err))
		return
	}
	response := dto.VisitsTodayResponse{
		Result: list,
		Count:  count,
	}
	render.JSON(w, r, response)
}
func StartOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

func EndOfDayInclusive(t time.Time) time.Time {
	return StartOfDay(t).Add(24*time.Hour - time.Second)
}

func DayRangeHalfOpen(t time.Time) (start, end time.Time) {
	start = StartOfDay(t)
	end = start.Add(24 * time.Hour) // [start, end)
	return
}
