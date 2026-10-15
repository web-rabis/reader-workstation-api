package model

type IFilters interface {
	Sql() string
}
