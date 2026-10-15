package reader

import (
	"context"

	readerClient "github.com/web-rabis/reader-client/client"
	readerModel "github.com/web-rabis/reader-client/model"
)

// IManager — тонкая обёртка над reader-client: reader-workstation-api больше
// не читает таблицы читателей напрямую, а ходит за ними в сервис reader по
// gRPC (как это уже сделано для ebook/order/control/dictionary).
type IManager interface {
	List(ctx context.Context, filter *readerModel.ReaderFilter, paging *readerModel.Paging) (int64, []*readerModel.Reader, error)
	ByTicketNumber(ctx context.Context, ticketNumber int64) (*readerModel.Reader, error)
	ByBarcode(ctx context.Context, barcode string) (*readerModel.Reader, error)
	Create(ctx context.Context, reader *readerModel.Reader, photo []byte) (int64, error)
	Update(ctx context.Context, reader *readerModel.Reader, photo []byte) error
	Delete(ctx context.Context, ticketNumber int64) error
	PhotoById(ctx context.Context, ticketNumber int64) ([]byte, error)
	UploadPhoto(ctx context.Context, ticketNumber int64, photo []byte, user *readerModel.User) error
	Approve(ctx context.Context, ticketNumber int64) (*readerModel.ReaderApproveResult, error)
	Reject(ctx context.Context, ticketNumber int64) error
	RenewalRegistration(ctx context.Context, ticketNumber int64, years, months, days int) error
}

type Manager struct {
	cl readerClient.ReaderService
}

func NewManager(cl readerClient.ReaderService) *Manager {
	return &Manager{cl: cl}
}

func (m *Manager) List(ctx context.Context, filter *readerModel.ReaderFilter, paging *readerModel.Paging) (int64, []*readerModel.Reader, error) {
	return m.cl.ReaderList(ctx, filter, paging)
}
func (m *Manager) ByTicketNumber(ctx context.Context, ticketNumber int64) (*readerModel.Reader, error) {
	return m.cl.ReaderByTicketNumber(ctx, ticketNumber)
}
func (m *Manager) ByBarcode(ctx context.Context, barcode string) (*readerModel.Reader, error) {
	return m.cl.ReaderByBarcode(ctx, barcode)
}
func (m *Manager) Create(ctx context.Context, reader *readerModel.Reader, photo []byte) (int64, error) {
	return m.cl.ReaderCreate(ctx, reader, photo)
}
func (m *Manager) Delete(ctx context.Context, ticketNumber int64) error {
	return m.cl.ReaderDelete(ctx, ticketNumber)
}
func (m *Manager) Update(ctx context.Context, reader *readerModel.Reader, photo []byte) error {
	return m.cl.ReaderUpdate(ctx, reader, photo)
}
func (m *Manager) PhotoById(ctx context.Context, ticketNumber int64) ([]byte, error) {
	return m.cl.ReaderPhotoByTicketNumber(ctx, ticketNumber)
}
func (m *Manager) UploadPhoto(ctx context.Context, ticketNumber int64, photo []byte, user *readerModel.User) error {
	return m.cl.ReaderUploadPhoto(ctx, ticketNumber, photo, user)
}
func (m *Manager) Approve(ctx context.Context, ticketNumber int64) (*readerModel.ReaderApproveResult, error) {
	return m.cl.ReaderApprove(ctx, ticketNumber)
}
func (m *Manager) Reject(ctx context.Context, ticketNumber int64) error {
	return m.cl.ReaderReject(ctx, ticketNumber)
}
func (m *Manager) RenewalRegistration(ctx context.Context, ticketNumber int64, years, months, days int) error {
	return m.cl.ReaderRenewalRegistration(ctx, ticketNumber, years, months, days)
}
