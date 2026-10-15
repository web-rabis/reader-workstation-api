package model

import (
	"net/http/httptest"
	"testing"
)

func parse(t *testing.T, query string) *Paging {
	t.Helper()
	p, err := PagingParseFromHttp(httptest.NewRequest("GET", "/?"+query, nil))
	if err != nil {
		t.Fatalf("%s: %v", query, err)
	}
	return p
}

func TestPagingParseFromHttp_Defaults(t *testing.T) {
	p := parse(t, "")
	if p.OrderBySql() != "name asc" || p.GetLimit() != 0 || p.GetOffset() != 0 {
		t.Errorf("defaults: %q limit=%d offset=%d", p.OrderBySql(), p.GetLimit(), p.GetOffset())
	}
}

func TestPagingParseFromHttp_PageToOffset(t *testing.T) {
	if p := parse(t, "limit=20&skip=3"); p.GetLimit() != 20 || p.GetOffset() != 40 {
		t.Errorf("limit=%d offset=%d, want 20/40", p.GetLimit(), p.GetOffset())
	}
	// без limit страница считается по 10 записей
	if p := parse(t, "skip=3"); p.GetOffset() != 20 {
		t.Errorf("offset=%d, want 20", p.GetOffset())
	}
}

func TestPagingParseFromHttp_Order(t *testing.T) {
	if got := parse(t, "orderBy=id&order=0").OrderBySql(); got != "id desc" {
		t.Errorf("order=0 -> %q, want 'id desc'", got)
	}
	if got := parse(t, "orderBy=id&order=1").OrderBySql(); got != "id asc" {
		t.Errorf("order=1 -> %q, want 'id asc'", got)
	}
}

func TestPagingParseFromHttp_BadNumbers(t *testing.T) {
	for _, q := range []string{"limit=x", "skip=x", "order=x"} {
		if _, err := PagingParseFromHttp(httptest.NewRequest("GET", "/?"+q, nil)); err == nil {
			t.Errorf("%s: ожидалась ошибка", q)
		}
	}
}
