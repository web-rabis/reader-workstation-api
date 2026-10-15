package model

import (
	"net/http"
	"strconv"

	dbmodel "github.com/web-rabis/db/model"
)

// Paging переиспользует общую реализацию пагинации из web-rabis/db —
// своей доменной логики здесь не было (offset/limit/sort общие для
// всех сервисов), поэтому дублировать код больше не нужно.
type Paging = dbmodel.Paging

// PagingParseFromHttp — перенос одноимённой функции из старой версии
// github.com/web-rabis/db/model (она ушла вместе с публичными полями
// Paging, когда либа перешла на приватные поля + сеттеры).
func PagingParseFromHttp(r *http.Request) (*Paging, error) {
	paging := dbmodel.NewPaging().SetSortKey("name").SetSortVal(1)

	if limitParam := r.URL.Query().Get("limit"); limitParam != "" {
		limit, err := strconv.Atoi(limitParam)
		if err != nil {
			return nil, err
		}
		paging.SetLimit(int64(limit))
	}

	if pageParam := r.URL.Query().Get("skip"); pageParam != "" {
		page, err := strconv.Atoi(pageParam)
		if err != nil {
			return nil, err
		}
		skip := page - 1
		if paging.GetLimit() != 0 {
			skip = int(paging.GetLimit()) * skip
		} else {
			skip = skip * 10
		}
		paging.SetOffset(int64(skip))
	}

	if orderParam := r.URL.Query().Get("order"); orderParam != "" {
		sortVal, err := strconv.Atoi(orderParam)
		if err != nil {
			return nil, err
		}

		if sortVal == 0 {
			sortVal = -1
		}

		//nolint:gosec
		paging.SetSortVal(int32(sortVal))
	}

	if sortKeyFromQuery := r.URL.Query().Get("orderBy"); sortKeyFromQuery != "" {
		paging.SetSortKey(sortKeyFromQuery)
	}

	return paging, nil
}
