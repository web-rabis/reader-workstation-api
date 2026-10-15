package control

import (
	"context"

	readerClient "github.com/web-rabis/reader-client/client"
	readerModel "github.com/web-rabis/reader-client/model"
)

type IManager interface {
	CheckIn(ctx context.Context, ticketNumber int64, department string, user *readerModel.User) error
	CheckOut(ctx context.Context, ticketNumber int64, department string, user *readerModel.User, ordersCount int64) error
	List(ctx context.Context, filters *readerModel.ControlFilter, paging *readerModel.Paging) (int64, []*readerModel.Control, error)
}

type Manager struct {
	controlCl readerClient.ControlService
}

func NewManager(controlCl readerClient.ControlService) *Manager {
	return &Manager{
		controlCl: controlCl,
	}
}

func (m *Manager) CheckIn(ctx context.Context, ticketNumber int64, department string, user *readerModel.User) error {
	return m.controlCl.ControlCheckIn(ctx, ticketNumber, user, department)
}
func (m *Manager) CheckOut(ctx context.Context, ticketNumber int64, department string, user *readerModel.User, ordersCount int64) error {
	return m.controlCl.ControlCheckOut(ctx, ticketNumber, user, department, ordersCount)
}
func (m *Manager) List(ctx context.Context, filters *readerModel.ControlFilter, paging *readerModel.Paging) (int64, []*readerModel.Control, error) {
	return m.controlCl.ControlList(ctx, filters, paging)
}
