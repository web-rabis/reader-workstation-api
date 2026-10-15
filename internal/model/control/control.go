package control

import "time"

type Control struct {
	TicketNumber  int64      `json:"ticketNumber" db:"ticket_num"`
	VisitDate     *time.Time `json:"visitDate" db:"vizit_date"`
	ExitDate      *time.Time `json:"exitDate" db:"vyhod"`
	UserVisit     string     `json:"userVisit" db:"sotrudnik_vizit"`
	UserExit      string     `json:"userExit" db:"sotrudnik_vyhod"`
	DepartmentF51 string     `json:"departmentF51" db:"f51"`
	Department    string     `json:"department" db:"otdel"`
	OrdersCount   int64      `json:"ordersCount" db:"zakazy"`
	Count         int64      `json:"count" db:"count"`
}
