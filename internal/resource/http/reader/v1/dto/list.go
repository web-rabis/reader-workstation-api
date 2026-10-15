package dto

import reader "github.com/web-rabis/reader-client/model"

type ReaderListResponse struct {
	Result []*reader.Reader `json:"result"`
	Count  int64            `json:"count"`
}
