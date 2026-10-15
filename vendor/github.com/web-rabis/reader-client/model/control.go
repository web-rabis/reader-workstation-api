package model

import (
	"time"

	"github.com/web-rabis/reader-client/protobuf"
)

type Control struct {
	TicketNumber  int64      `json:"ticketNumber,omitempty"`
	VisitDate     *time.Time `json:"visitDate,omitempty"`
	ExitDate      *time.Time `json:"exitDate,omitempty"`
	UserVisit     string     `json:"userVisit,omitempty"`
	UserExit      string     `json:"userExit,omitempty"`
	DepartmentF51 string     `json:"departmentF51,omitempty"`
	Department    string     `json:"department,omitempty"`
	OrdersCount   int64      `json:"ordersCount,omitempty"`
	Count         int64      `json:"count,omitempty"`
}

func ControlFromProto(c *protobuf.Control) *Control {
	if c == nil {
		return nil
	}
	var visitDate *time.Time
	if c.VisitDate != nil {
		t := c.VisitDate.AsTime()
		visitDate = &t
	}
	var exitDate *time.Time
	if c.ExitDate != nil {
		t := c.ExitDate.AsTime()
		exitDate = &t
	}
	return &Control{
		TicketNumber:  c.TicketNumber,
		VisitDate:     visitDate,
		ExitDate:      exitDate,
		UserVisit:     c.UserVisit,
		UserExit:      c.UserExit,
		DepartmentF51: c.DepartmentF51,
		Department:    c.Department,
		OrdersCount:   c.OrdersCount,
		Count:         c.Count,
	}
}

func ControlListFromProto(l []*protobuf.Control) []*Control {
	list := make([]*Control, len(l))
	for i, d := range l {
		list[i] = ControlFromProto(d)
	}
	return list
}
