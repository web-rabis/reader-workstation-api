package eorder

import (
	"context"

	orderClient "github.com/web-rabis/order-client/client"
	ordermodel "github.com/web-rabis/order-client/model"
)

// IManager — тонкая обёртка над order-client: reader больше не читает
// таблицы заказов напрямую, а ходит за ними в сервис orders по gRPC (как
// это уже сделано в circulation-api).
type IManager interface {
	List(ctx context.Context, filter *ordermodel.OrderFilters, paging *ordermodel.Paging) ([]*ordermodel.Order, error)
	Count(ctx context.Context, filter *ordermodel.OrderFilters) (int64, error)
}

type Manager struct {
	cl orderClient.OrderService
}

func NewManager(cl orderClient.OrderService) *Manager {
	return &Manager{cl: cl}
}

func (m *Manager) List(ctx context.Context, filter *ordermodel.OrderFilters, paging *ordermodel.Paging) ([]*ordermodel.Order, error) {
	_, orders, err := m.cl.List(ctx, paging, filter)
	return orders, err
}
func (m *Manager) Count(ctx context.Context, filter *ordermodel.OrderFilters) (int64, error) {
	count, _, err := m.cl.List(ctx, &ordermodel.Paging{Limit: 1}, filter)
	return count, err
}
