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
func (res *Resource) visitsStatistics(w http.ResponseWriter, r *http.Request) {
	exitDateIsNull := true
	start, end := DayRangeHalfOpen(time.Now())
	//start = start.AddDate(0, 0, -30)
	filterAll := &readerModel.ControlFilter{
		VisitDateLow:  &start,
		VisitDateHigh: &end,
	}
	filterInHall := &readerModel.ControlFilter{
		VisitDateLow:   &start,
		VisitDateHigh:  &end,
		ExitDateIsNull: &exitDateIsNull,
	}
	paging := &readerModel.Paging{
		Offset: 0,
		Limit:  1,
	}
	countAll, _, err := res.controlMan.List(r.Context(), filterAll, paging)
	if err != nil {
		_ = render.Render(w, r, httperrors.BadRequest(err))
		return
	}
	countInHall, _, err := res.controlMan.List(r.Context(), filterInHall, paging)
	if err != nil {
		_ = render.Render(w, r, httperrors.BadRequest(err))
		return
	}
	response := dto.VisitsStatisticsResponse{
		All:    countAll,
		InHall: countInHall,
	}
	render.JSON(w, r, response)
}
