package model

import (
	"time"

	"github.com/web-rabis/reader-client/protobuf"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ReaderFilter struct {
	Status             string
	Query              string
	Iin                string
	ExpirationDateLow  *time.Time
	ExpirationDateHigh *time.Time
	IsEmployee         *bool
}
type ControlFilter struct {
	VisitDateLow   *time.Time
	VisitDateHigh  *time.Time
	ExitDateLow    *time.Time
	ExitDateHigh   *time.Time
	ExitDateIsNull *bool
}

func (f *ReaderFilter) ToProto() *protobuf.ReaderFilter {
	if f == nil {
		return nil
	}
	var expirationDateLow *timestamppb.Timestamp
	var expirationDateHigh *timestamppb.Timestamp
	if f.ExpirationDateLow != nil {
		expirationDateLow = timestamppb.New(*f.ExpirationDateLow)
	}
	if f.ExpirationDateHigh != nil {
		expirationDateHigh = timestamppb.New(*f.ExpirationDateHigh)
	}
	return &protobuf.ReaderFilter{
		Status:             f.Status,
		Query:              f.Query,
		Iin:                f.Iin,
		ExpirationDateLow:  expirationDateLow,
		ExpirationDateHigh: expirationDateHigh,
		IsEmployee:         f.IsEmployee,
	}
}

func (f *ControlFilter) ToProto() *protobuf.ControlFilter {
	if f == nil {
		return nil
	}
	var visitDateLow *timestamppb.Timestamp
	if f.VisitDateLow != nil {
		visitDateLow = timestamppb.New(*f.VisitDateLow)
	}
	var visitDateHigh *timestamppb.Timestamp
	if f.VisitDateHigh != nil {
		visitDateHigh = timestamppb.New(*f.VisitDateHigh)
	}
	var exitDateLow *timestamppb.Timestamp
	if f.ExitDateLow != nil {
		exitDateLow = timestamppb.New(*f.ExitDateLow)
	}
	var exitDateHigh *timestamppb.Timestamp
	if f.ExitDateHigh != nil {
		exitDateHigh = timestamppb.New(*f.ExitDateHigh)
	}
	return &protobuf.ControlFilter{
		VisitDateLow:   visitDateLow,
		VisitDateHigh:  visitDateHigh,
		ExitDateLow:    exitDateLow,
		ExitDateHigh:   exitDateHigh,
		ExitDateIsNull: f.ExitDateIsNull,
	}
}
