package grpc

import (
	"context"

	"github.com/web-rabis/reader-client/client"
	"github.com/web-rabis/reader-client/model"
	"github.com/web-rabis/reader-client/protobuf"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ControlService struct {
	client protobuf.ControlSvcClient
}

func NewControlServiceClient(client protobuf.ControlSvcClient) client.ControlService {
	return &ControlService{
		client: client,
	}
}
func (c ControlService) ControlCheckIn(ctx context.Context, ticketNumber int64, user *model.User, department string) error {
	request := &protobuf.ControlCheckInRequest{
		TicketNumber: ticketNumber,
		User:         user.ToProto(),
		Department:   department,
	}
	_, err := c.client.ControlCheckIn(ctx, request)
	switch status.Code(err) {
	case codes.OK:
		return nil
	default:
		return err
	}
}

func (c ControlService) ControlCheckOut(ctx context.Context, ticketNumber int64, user *model.User, department string, ordersCount int64) error {
	request := &protobuf.ControlCheckOutRequest{
		TicketNumber: ticketNumber,
		User:         user.ToProto(),
		Department:   department,
		OrdersCount:  ordersCount,
	}
	_, err := c.client.ControlCheckOut(ctx, request)
	switch status.Code(err) {
	case codes.OK:
		return nil
	default:
		return err
	}
}

func (c ControlService) ControlList(ctx context.Context, filters *model.ControlFilter, paging *model.Paging) (int64, []*model.Control, error) {
	request := &protobuf.ControlListRequest{
		Filters: filters.ToProto(),
		Paging:  paging.ToProto(),
	}
	response, err := c.client.ControlList(ctx, request)
	switch status.Code(err) {
	case codes.OK:
		return response.Count, model.ControlListFromProto(response.Result), nil
	default:
		return 0, nil, err
	}
}

var _ client.ControlService = &ControlService{}
