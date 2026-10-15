package dto

import (
	reader "github.com/web-rabis/reader-client/model"
)

type ByIdResponse struct {
	Result *reader.Reader `json:"result"`
}
