package model

import (
	"fmt"
	"strings"
)

// QueryBuilder собирает SQL-запрос в правильном порядке клаузул.
// Порядок: SELECT/FROM/JOIN → WHERE → GROUP BY → HAVING → ORDER BY → LIMIT/OFFSET
type QueryBuilder struct {
	base    string
	where   []string
	groupBy string
	having  []string
	orderBy string
	limit   *int64
	offset  *int64
}

func NewQueryBuilder(base string) *QueryBuilder {
	return &QueryBuilder{
		base:   base,
		where:  make([]string, 0),
		having: make([]string, 0),
	}
}

// Where добавляет условие WHERE. Несколько вызовов объединяются через AND.
func (q *QueryBuilder) Where(clause string) *QueryBuilder {
	if clause != "" {
		q.where = append(q.where, clause)
	}
	return q
}

// GroupBy задаёт GROUP BY. Повторный вызов перезаписывает.
func (q *QueryBuilder) GroupBy(clause string) *QueryBuilder {
	q.groupBy = clause
	return q
}

// Having добавляет условие HAVING. Несколько вызовов объединяются через AND.
func (q *QueryBuilder) Having(clause string) *QueryBuilder {
	if clause != "" {
		q.having = append(q.having, clause)
	}
	return q
}

// OrderBy задаёт ORDER BY. Повторный вызов перезаписывает.
func (q *QueryBuilder) OrderBy(clause string) *QueryBuilder {
	q.orderBy = clause
	return q
}

// Limit задаёт LIMIT.
func (q *QueryBuilder) Limit(n int64) *QueryBuilder {
	q.limit = &n
	return q
}

// Offset задаёт OFFSET.
func (q *QueryBuilder) Offset(n int64) *QueryBuilder {
	q.offset = &n
	return q
}

// Build собирает итоговый SQL в правильном порядке клаузул.
func (q *QueryBuilder) Build() string {
	var sb strings.Builder

	sb.WriteString(q.base)

	if len(q.where) > 0 {
		sb.WriteString(" WHERE ")
		sb.WriteString(strings.Join(q.where, " AND "))
	}

	if q.groupBy != "" {
		sb.WriteString(" GROUP BY ")
		sb.WriteString(q.groupBy)
	}

	if len(q.having) > 0 {
		sb.WriteString(" HAVING ")
		sb.WriteString(strings.Join(q.having, " AND "))
	}

	if q.orderBy != "" {
		sb.WriteString(" ORDER BY ")
		sb.WriteString(q.orderBy)
	}

	if q.limit != nil {
		sb.WriteString(fmt.Sprintf(" LIMIT %d", *q.limit))
	}

	if q.offset != nil {
		sb.WriteString(fmt.Sprintf(" OFFSET %d", *q.offset))
	}

	return sb.String()
}

// Clone возвращает независимую копию билдера.
// Используется чтобы не мутировать базовый билдер между вызовами.
func (q *QueryBuilder) Clone() *QueryBuilder {
	clone := &QueryBuilder{
		base:    q.base,
		where:   make([]string, len(q.where)),
		groupBy: q.groupBy,
		having:  make([]string, len(q.having)),
		orderBy: q.orderBy,
	}
	copy(clone.where, q.where)
	copy(clone.having, q.having)
	if q.limit != nil {
		v := *q.limit
		clone.limit = &v
	}
	if q.offset != nil {
		v := *q.offset
		clone.offset = &v
	}
	return clone
}

// ApplyFilters применяет model.Filters к билдеру.
// Автоматически обрезает ведущий WHERE если он есть в строке фильтра.
func (q *QueryBuilder) ApplyFilters(filters IFilters) *QueryBuilder {
	if filters == nil {
		return q
	}
	sql := strings.TrimSpace(filters.Sql())
	if sql == "" {
		return q
	}
	if strings.HasPrefix(strings.ToUpper(sql), "WHERE ") {
		sql = sql[6:]
	}
	return q.Where(sql)
}

// ApplyPaging применяет model.Paging к билдеру.
func (q *QueryBuilder) ApplyPaging(paging IPaging) *QueryBuilder {
	if paging == nil {
		return q
	}
	if order := paging.OrderBySql(); order != "" {
		q.OrderBy(order)
	}
	if limit := paging.GetLimit(); limit > 0 {
		q.Limit(limit)
	}
	if offset := paging.GetOffset(); offset > 0 {
		q.Offset(offset)
	}
	return q
}
