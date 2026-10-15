package ebook

import (
	"context"

	ebookClient "github.com/web-rabis/ebook-client/client"
	ebookModel "github.com/web-rabis/ebook-client/model"
	"github.com/web-rabis/ebook-client/model/ebook"
)

// IManager — тонкая обёртка над ebook-client: reader больше не читает
// таблицы книг напрямую, а ходит за ними в сервис ebook по gRPC (как это
// уже сделано в circulation-api).
type IManager interface {
	EbookById(ctx context.Context, id int64) (*ebook.Ebook, error)
	InvById(ctx context.Context, ebookId, invId int64) (*ebook.Inv, error)
}

type Manager struct {
	cl ebookClient.EbookService
}

func NewManager(cl ebookClient.EbookService) *Manager {
	return &Manager{cl: cl}
}

func (m *Manager) EbookById(ctx context.Context, id int64) (*ebook.Ebook, error) {
	return m.cl.EbookById(ctx, id, false)
}

// InvById — у ebook-client нет метода "по id", экземпляр ищется в списке
// по EbookId (так же сделано в circulation-api/internal/domain/manager/order).
func (m *Manager) InvById(ctx context.Context, ebookId, invId int64) (*ebook.Inv, error) {
	_, list, err := m.cl.InvList(ctx, &ebookModel.InvFilters{EbookId: ebookId}, nil)
	if err != nil {
		return nil, err
	}
	for _, inv := range list {
		if inv.Id == invId {
			return inv, nil
		}
	}
	return nil, nil
}
