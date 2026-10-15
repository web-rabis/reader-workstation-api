package control

import "time"

type Filter struct {
	VisitDateLow   *time.Time
	VisitDateHigh  *time.Time
	ExitDateLow    *time.Time
	ExitDateHigh   *time.Time
	ExitDateIsNull *bool
}

func (f *Filter) Sql() string {
	if f == nil {
		return ""
	}
	var sql = ""
	if f.VisitDateLow != nil || f.VisitDateHigh != nil {
		if sql != "" {
			sql += " and "
		}
		if f.VisitDateLow != nil && f.VisitDateHigh != nil {
			sql = sql + "vizit_date between '" + f.VisitDateLow.Format("2006-01-02 15:04:05") + "' and '" + f.VisitDateHigh.Format("2006-01-02 15:04:05") + "'"
		} else if f.VisitDateLow != nil {
			sql = sql + "vizit_date >= '" + f.VisitDateLow.Format("2006-01-02 15:04:05") + "'"
		} else if f.VisitDateHigh != nil {
			sql = sql + "vizit_date <= '" + f.VisitDateHigh.Format("2006-01-02 15:04:05") + "'"
		}
	}
	if f.ExitDateLow != nil || f.ExitDateHigh != nil {
		if sql != "" {
			sql += " and "
		}
		if f.ExitDateLow != nil && f.ExitDateHigh != nil {
			sql = sql + "vyhod between '" + f.ExitDateLow.Format("2006-01-02 15:04:05") + "' and '" + f.ExitDateHigh.Format("2006-01-02 15:04:05") + "'"
		} else if f.ExitDateLow != nil {
			sql = sql + "vyhod >= '" + f.ExitDateLow.Format("2006-01-02 15:04:05") + "'"
		} else if f.ExitDateHigh != nil {
			sql = sql + "vyhod <= '" + f.ExitDateHigh.Format("2006-01-02 15:04:05") + "'"
		}
	}
	if f.ExitDateIsNull != nil {
		if sql != "" {
			sql += " and "
		}
		if *f.ExitDateIsNull {
			sql = sql + "vyhod is null"
		} else {
			sql = sql + "vyhod is not null"
		}
	}
	return sql
}
