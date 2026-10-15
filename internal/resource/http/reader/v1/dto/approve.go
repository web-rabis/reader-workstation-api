package dto

import (
	"time"
)

type ApproveReaderResponse struct {
	OldTicketNumber int64     `json:"oldTicketNumber"`
	TicketNumber    int64     `json:"ticketNumber"`
	Barcode         string    `json:"barcode"`
	Status          string    `json:"status"`
	ExpirationDate  time.Time `json:"expirationDate"`
}
