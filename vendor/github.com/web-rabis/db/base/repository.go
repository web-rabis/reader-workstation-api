package base

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/web-rabis/db/model"
)

var (
	ErrRecordNotExist = errors.New("record not exists")
	// ErrConcurrentModification - запись существует, но не соответствует
	// ожидаемой версии (edit_date/аналог), переданной вызывающей стороной:
	// кто-то успел изменить её между чтением и записью. Общий sentinel для
	// optimistic-конкурентности в любом сервисе, использующем этот пакет -
	// конкретика (какая колонка играет роль версии, как отличить от
	// ErrRecordNotExist) остаётся на стороне репозитория.
	ErrConcurrentModification = errors.New("record was concurrently modified")
)

// Queryer - минимальный набор методов, которым пользуются хелперы этого
// файла. Ему одинаково удовлетворяют *pgxpool.Pool и pgx.Tx, поэтому любой
// вызов ниже можно выполнить как напрямую через пул, так и внутри
// транзакции (tx, err := pool.Begin(ctx); ...; dbbase.Exec(ctx, tx, ...)),
// без каких-либо изменений на стороне самих хелперов.
type Queryer interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

type Mapper[T any] func(values []any) *T

// IdMapper - Mapper для insert ... returning id, где id - целочисленный
// (integer/serial, в pgx приезжает как int32). Общий для всех сервисов,
// которым для Create нужен только новый id одной колонкой - чтобы не
// заводить одну и ту же замыкание-обёртку в каждом репозитории.
func IdMapper(v []any) *int64 {
	id := int64(v[0].(int32))
	return &id
}

// QueryListWithMapper - как QueryList, но T собирается вручную из
// []any (позиционно) через переданный Mapper[T], а не по имени колонки.
// Нужен для составных структур, которые сами вызывают мапперы вложенных
// типов на срезе строки (см., например, ebook.NewEbookFromResult в
// сервисе ebook) - такую вложенность по имени колонки не выразить.
func QueryListWithMapper[T any](
	ctx context.Context,
	q Queryer,
	qb *model.QueryBuilder,
	filters model.IFilters,
	paging model.IPaging,
	mapper Mapper[T],
	args ...any,
) ([]*T, error) {
	sql := qb.Clone().
		ApplyFilters(filters).
		ApplyPaging(paging).
		Build()

	log.Printf("[DEBUG] query sql: %s, args:%s", sql, args)
	rows, err := q.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*T, 0)

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, err
		}

		item := mapper(values)
		result = append(result, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// QueryOneWithMapper - как QueryOne, но T собирается вручную из []any
// (позиционно) через переданный Mapper[T], а не по имени колонки. См.
// комментарий к QueryListWithMapper.
func QueryOneWithMapper[T any](
	ctx context.Context,
	q Queryer,
	sql string,
	mapper Mapper[T],
	args ...any,
) (*T, error) {
	log.Printf("[DEBUG] query sql: %s, args:%s", sql, args)
	rows, err := q.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, ErrRecordNotExist
	}

	values, err := rows.Values()
	if err != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return mapper(values), nil
}

// QueryList сопоставляет колонки с полями T по имени (регистронезависимо,
// либо по тегу `db:"..."`) через встроенный в pgx
// pgx.RowToAddrOfStructByName - без ручного Mapper[T]. Подходит для
// плоских структур; для составных (вложенная композиция через срез
// значений строки) используйте QueryListWithMapper.
func QueryList[T any](
	ctx context.Context,
	q Queryer,
	qb *model.QueryBuilder,
	filters model.IFilters,
	paging model.IPaging,
	args ...any,
) ([]*T, error) {
	sql := qb.Clone().
		ApplyFilters(filters).
		ApplyPaging(paging).
		Build()

	log.Printf("[DEBUG] query sql: %s, args:%s", sql, args)
	rows, err := q.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows[*T](rows, pgx.RowToAddrOfStructByName[T])
}

// QueryOne - как QueryList, но для одной строки. Если строк нет -
// ErrRecordNotExist, как и у QueryOneWithMapper.
func QueryOne[T any](
	ctx context.Context,
	q Queryer,
	sql string,
	args ...any,
) (*T, error) {
	log.Printf("[DEBUG] query sql: %s, args:%s", sql, args)
	rows, err := q.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	v, err := pgx.CollectOneRow[*T](rows, pgx.RowToAddrOfStructByName[T])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrRecordNotExist
		}
		return nil, err
	}

	return v, nil
}

func QueryCount(
	ctx context.Context,
	q Queryer,
	qb *model.QueryBuilder,
	filters model.IFilters,
	args ...any,
) (int64, error) {
	sql := qb.Clone().
		ApplyFilters(filters).
		Build()
	log.Printf("[DEBUG] query sql: %s, args:%s", sql, args)
	cnt, err := QueryOneWithMapper(
		ctx,
		q,
		sql,
		func(values []any) *int64 {
			cnt := values[0].(int64)
			return &cnt
		},
		args...,
	)
	if err != nil {
		return 0, err
	}
	return *cnt, nil
}

// Exec универсален для insert/update/delete без returning.
func Exec(
	ctx context.Context,
	q Queryer,
	sql string,
	args ...any,
) (pgconn.CommandTag, error) {
	log.Printf("[DEBUG] exec sql: %s, args:%s", sql, args)
	return q.Exec(ctx, sql, args...)
}

// Create универсален для insert ... returning ...
func Create[T any](
	ctx context.Context,
	q Queryer,
	sql string,
	mapper Mapper[T],
	args ...any,
) (*T, error) {
	return QueryOneWithMapper(ctx, q, sql, mapper, args...)
}

// Update выполняет update и возвращает количество затронутых строк.
func Update(
	ctx context.Context,
	q Queryer,
	sql string,
	args ...any,
) (int64, error) {
	tag, err := Exec(ctx, q, sql, args...)
	if err != nil {
		return 0, err
	}

	return tag.RowsAffected(), nil
}

// Delete выполняет delete и возвращает количество затронутых строк.
func Delete(
	ctx context.Context,
	q Queryer,
	sql string,
	args ...any,
) (int64, error) {
	tag, err := Exec(ctx, q, sql, args...)
	if err != nil {
		return 0, err
	}

	if tag.RowsAffected() == 0 {
		return 0, ErrRecordNotExist
	}

	return tag.RowsAffected(), nil
}
