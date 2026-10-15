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

// @Summary List Readers
// @Description Returns a paginated and filterable list of readers
// @Accept json
// @Produce json
// @Tags Reader
// @Param status query string false "Reader status"
// @Param iin query string false "Reader IIN"
// @Param query query string false "Search query"
// @Param isEmployee query bool false "Filter by employee flag"
// @Param expirationDateLow query string false "Expiration date from (YYYY-MM-DD)"
// @Param expirationDateHigh query string false "Expiration date to (YYYY-MM-DD)"
// @Param limit query int false "Page size"
// @Param skip query int false "Page number"
// @Param order query int false "Sort order: 1 or -1"
// @Param orderBy query string false "Sort key"
// @Success 200 {object} dto.ReaderListResponse
// @Failure 400 {object} httperrors.Response
// @Failure 401 {object} httperrors.Response
// @Failure 403 {object} httperrors.Response
// @Failure 404 {object} httperrors.Response
// @Failure 500 {object} httperrors.Response
// @Router /api/v1/reader [get]
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
