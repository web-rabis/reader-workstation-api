package model

type IPaging interface {
	OrderBySql() string
	GetLimit() int64
	GetOffset() int64
}

// Paging — готовая реализация IPaging без каких-либо доменных полей.
// Сервисы, которым не нужна собственная логика пагинации/сортировки,
// могут использовать её напрямую вместо того, чтобы копировать этот
// код у себя (как раньше было в internal/domain/model.Paging в ebook).
type Paging struct {
	offset  int64
	limit   int64
	sortKey string
	sortVal int32
}

func NewPaging() *Paging {
	return &Paging{}
}

func (p *Paging) SetLimit(limit int64) *Paging {
	p.limit = limit
	return p
}

func (p *Paging) SetOffset(offset int64) *Paging {
	p.offset = offset
	return p
}

func (p *Paging) SetSortKey(sortKey string) *Paging {
	p.sortKey = sortKey
	return p
}

func (p *Paging) SetSortVal(sortVal int32) *Paging {
	p.sortVal = sortVal
	return p
}

func (p *Paging) OrderBySql() string {
	if p == nil || p.sortKey == "" {
		return ""
	}

	if p.sortVal == -1 {
		return p.sortKey + " desc"
	}
	return p.sortKey + " asc"
}

func (p *Paging) GetLimit() int64 {
	if p == nil {
		return 0
	}
	return p.limit
}

func (p *Paging) GetOffset() int64 {
	if p == nil {
		return 0
	}
	return p.offset
}
