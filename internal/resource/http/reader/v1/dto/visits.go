package dto

import (
	"time"

	ordermodel "github.com/web-rabis/order-client/model"
	readerModel "github.com/web-rabis/reader-client/model"
)

type VisitsStatusResponse struct {
	Reader *readerModel.Reader `json:"reader"`
	Status string              `json:"status"`
	Since  *time.Time          `json:"since"`
	Orders []Order             `json:"orders"`
}
type VisitsTodayResponse struct {
	Result []*readerModel.Control `json:"result"`
	Count  int64                  `json:"count"`
}
type VisitsStatisticsResponse struct {
	All    int64 `json:"all"`
	InHall int64 `json:"inHall"`
}

type Order struct {
	Id         int64                  `json:"id,omitempty"`
	Type       string                 `json:"type"`
	OrderedAt  time.Time              `json:"orderedAt,omitempty"`
	Status     *ordermodel.State      `json:"status,omitempty"`
	Department *ordermodel.Department `json:"department,omitempty"`
	Ebook      *Ebook                 `json:"ebook"`
	Magazine   *Magazine              `json:"magazine"`
}
type Ebook struct {
	Id     int64     `json:"id"`
	Author string    `json:"author"`
	Title  string    `json:"title"`
	Inv    *EbookInv `json:"inv"`
}
type EbookInv struct {
	Id        int64  `json:"id"`
	InvNumber string `json:"invNumber"`
	Barcode   string `json:"barcode"`
}
type Magazine struct {
	Id     int64        `json:"id"`
	Nkr    int64        `json:"nkr"`
	Title  string       `json:"title"`
	Number string       `json:"number"`
	Year   string       `json:"year"`
	Inv    *MagazineInv `json:"inv"`
}
type MagazineInv struct {
	Id        int64  `json:"id"`
	InvNumber string `json:"invNumber"`
	Barcode   string `json:"barcode"`
}
