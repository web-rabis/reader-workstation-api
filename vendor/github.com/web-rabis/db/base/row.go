package base

import "time"

// Str, StrPtr, Int64, Int64Ptr, Time, TimePtr, BoolFromInt - типизированные
// хелперы для составных Mapper[T] func(v []any) *T (см. QueryOneWithMapper/
// QueryListWithMapper), где T собирается вручную из позиционного среза
// значений строки - обычно потому что T вкладывает другой тип через
// подсрез v[i:j] (составная композиция), что не выразить через
// QueryOne/QueryList (сопоставление по имени колонки). Убирают повтор
// вида `if v[i] != nil { f = v[i].(string) }` на каждое поле.

// Str возвращает значение колонки i как string. NULL -> "".
func Str(v []any, i int) string {
	if v[i] == nil {
		return ""
	}
	return v[i].(string)
}

// StrPtr - как Str, но NULL -> nil вместо "".
func StrPtr(v []any, i int) *string {
	if v[i] == nil {
		return nil
	}
	s := v[i].(string)
	return &s
}

// Int64 возвращает значение колонки i как int64, независимо от того,
// какой ширины integer вернул pgx (smallint/integer/bigint -> int16/
// int32/int64). NULL -> 0.
func Int64(v []any, i int) int64 {
	switch x := v[i].(type) {
	case int16:
		return int64(x)
	case int32:
		return int64(x)
	case int64:
		return x
	}
	return 0
}

// Int64Ptr - как Int64, но NULL -> nil вместо 0.
func Int64Ptr(v []any, i int) *int64 {
	if v[i] == nil {
		return nil
	}
	n := Int64(v, i)
	return &n
}

// Time возвращает значение колонки i как time.Time. NULL -> нулевое
// значение time.Time{}.
func Time(v []any, i int) time.Time {
	if v[i] == nil {
		return time.Time{}
	}
	return v[i].(time.Time)
}

// TimePtr - как Time, но NULL -> nil.
func TimePtr(v []any, i int) *time.Time {
	if v[i] == nil {
		return nil
	}
	t := v[i].(time.Time)
	return &t
}

// BoolFromInt - для колонок вида smallint 0/1 (конвенция этой БД вместо
// boolean, см. krv/is_repeat/required). NULL или любое значение кроме 1 -
// false.
func BoolFromInt(v []any, i int) bool {
	return Int64(v, i) == 1
}
