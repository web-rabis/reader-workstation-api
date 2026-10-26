package v1

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/web-rabis/httperrors"
	ordermodel "github.com/web-rabis/order-client/model"
	readerModel "github.com/web-rabis/reader-client/model"
	"github.com/web-rabis/reader-workstation-api/internal/resource/http/reader/v1/dto"
)

// @Summary Reader Visit Status
// @Description Returns a reader's current visit status by ticket number or barcode, including any active orders
// @Accept json
// @Produce json
// @Tags Reader
// @Param id path string true "Reader Ticket Number or Barcode"
// @Success 200 {object} dto.VisitsStatusResponse
// @Failure 400 {object} httperrors.Response
// @Failure 401 {object} httperrors.Response
// @Failure 403 {object} httperrors.Response
// @Failure 404 {object} httperrors.Response
// @Failure 500 {object} httperrors.Response
// @Router /api/v1/reader/{id}/visits/status [get]
func (res *Resource) visitsStatus(w http.ResponseWriter, r *http.Request) {
	reader_, status, err := res.getReaderByBarcodeOrTicketNumber(r)
	if err != nil {
		_ = render.Render(w, r, httperrors.BadRequest(err))
		return
	}
	if status != "in" && status != "out" {
		var lastVisitDate *time.Time
		if reader_ != nil && reader_.LastVisitDate != nil {
			lastVisitDate = reader_.LastVisitDate
		}
		render.JSON(w, r, dto.VisitsStatusResponse{
			Reader: reader_,
			Status: status,
			Since:  lastVisitDate,
			Orders: make([]dto.Order, 0),
		})
		return
	}
	orders, err := res.getOrders(r.Context(), reader_)
	if err != nil {
		_ = render.Render(w, r, httperrors.Internal(err))
		return
	}
	if len(orders) > 0 {
		status = "activeOrders"
	}
	response := dto.VisitsStatusResponse{
		Status: status,
		Since:  reader_.LastVisitDate,
		Orders: orders,
		Reader: reader_,
	}
	render.JSON(w, r, response)
}
func (res *Resource) getReaderByBarcodeOrTicketNumber(r *http.Request) (*readerModel.Reader, string, error) {
	id := chi.URLParam(r, "id")
	if id == "" {
		return nil, "ticketNumberOrBarcodeRequired", nil
	}
	var reader_ *readerModel.Reader
	var err error
	reader_, err = res.readerMan.ByBarcode(r.Context(), id)
	if err != nil {
		var ticketNumber int64
		ticketNumber, err = strconv.ParseInt(id, 10, 64)
		if err != nil {
			return nil, "ticketNumberOrBarcodeInvalid", nil
		}
		reader_, err = res.readerMan.ByTicketNumber(r.Context(), ticketNumber)
		if err != nil {
			return nil, "readerNotFound", nil
		}
	}
	if reader_.ExpirationDate != nil && reader_.ExpirationDate.Before(time.Now()) {
		return reader_, "expired", nil
	}
	filter := &readerModel.ReaderFilter{
		Iin: reader_.Iin,
	}
	_, readers, err := res.readerMan.List(r.Context(), filter, nil)
	if err != nil {
		return nil, "", errors.New("reader by iin not found")
	}
	var iinCnt = 0
	for _, rd := range readers {
		if rd.Status == "" && rd.PriznakUdal != "udalen" {
			iinCnt++
		}
	}
	if iinCnt > 1 {
		return reader_, "checkIin", nil
	}
	if reader_.Department == "*********" {
		return reader_, "out", nil
	}

	return reader_, "in", nil
}
func (res *Resource) getOrders(ctx context.Context, reader_ *readerModel.Reader) ([]dto.Order, error) {
	orderFilter := ordermodel.OrderFilters{
		TicketNumber: reader_.TicketNumber,
		States:       []string{ordermodel.OrderStateInHands},
	}
	paging := &ordermodel.Paging{SortKey: "ordered_date", SortVal: -1}
	orders, err := res.orderMan.List(ctx, &orderFilter, paging)
	if err != nil {
		return nil, errors.New("orders by ticket number not found")
	}
	dtoOrders := make([]dto.Order, len(orders), len(orders))
	for i, order := range orders {
		o := dto.Order{
			Id:         order.Id,
			OrderedAt:  order.CreatedAt,
			Department: order.Department,
			Status:     order.State,
		}
		if order.EbookId != 0 {
			o.Type = "ebook"
			b, err := res.ebookMan.EbookById(ctx, order.EbookId)
			if err != nil {
				return nil, err
			}
			o.Ebook = &dto.Ebook{
				Id:     b.Id,
				Author: b.Author,
				Title:  b.Title,
			}
			if order.EbookInvId != 0 {
				ebookInv, err := res.ebookMan.InvById(ctx, order.EbookId, order.EbookInvId)
				if err != nil {
					return nil, err
				}
				o.Ebook.Inv = &dto.EbookInv{
					Id:        order.EbookInvId,
					InvNumber: order.InvNumber,
					Barcode:   order.Barcode,
				}
				if ebookInv != nil {
					o.Ebook.Inv.InvNumber = ebookInv.InvNumber
					o.Ebook.Inv.Barcode = ebookInv.Barcode
				}
			} else {
				o.Ebook.Inv = &dto.EbookInv{
					InvNumber: order.InvNumber,
					Barcode:   order.Barcode,
				}
			}
		} else {
			o.Type = "magazine"
			o.Magazine = &dto.Magazine{
				Inv: &dto.MagazineInv{
					Id:        0,
					InvNumber: order.InvNumber,
					Barcode:   order.Barcode,
				},
			}
			if order.Periodical != nil {
				o.Magazine.Id = order.Periodical.Id
				o.Magazine.Nkr = order.Periodical.Nkr
				o.Magazine.Title = order.Periodical.Title
				o.Magazine.Number = order.Periodical.Number
				o.Magazine.Year = order.Periodical.YearEdition
			}
		}
		dtoOrders[i] = o
	}
	return dtoOrders, nil
}
