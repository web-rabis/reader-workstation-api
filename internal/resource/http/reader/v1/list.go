package v1

import (
	"net/http"
	"strconv"
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
func (res *Resource) list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	paging, err := PagingParse(r)
	if err != nil {
		_ = render.Render(w, r, httperrors.BadRequest(err))
		return
	}
	filter, err := parseReaderFilter(r)
	if err != nil {
		_ = render.Render(w, r, httperrors.BadRequest(err))
		return
	}
	count, list, err := res.readerMan.List(ctx, filter, paging)
	if err != nil {
		_ = render.Render(w, r, httperrors.Internal(err))
		return
	}
	render.JSON(w, r, dto.ReaderListResponse{
		Result: list,
		Count:  count,
	})
}

func parseReaderFilter(r *http.Request) (*readerModel.ReaderFilter, error) {
	filter := readerModel.ReaderFilter{}

	if status := r.URL.Query().Get("status"); status != "" && status != "*" {
		filter.Status = status
	}
	if iin := r.URL.Query().Get("iin"); iin != "" {
		filter.Iin = iin
	}
	if expirationDate := r.URL.Query().Get("expirationDateLow"); expirationDate != "" {
		expDate, err := time.Parse("2006-01-02", expirationDate)
		if err == nil {
			filter.ExpirationDateLow = &expDate
		}
	}
	if expirationDate := r.URL.Query().Get("expirationDateHigh"); expirationDate != "" {
		expDate, err := time.Parse("2006-01-02", expirationDate)
		if err == nil {
			filter.ExpirationDateHigh = &expDate
		}
	}
	if query := r.URL.Query().Get("query"); query != "" {
		filter.Query = query
	}
	if isEmployeeStr := r.URL.Query().Get("isEmployee"); isEmployeeStr != "" {
		isEmployee, err := strconv.ParseBool(isEmployeeStr)
		if err == nil {
			filter.IsEmployee = &isEmployee
		}
	}
	return &filter, nil
}
