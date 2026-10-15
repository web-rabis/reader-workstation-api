package v1

import (
	"encoding/json"
	"errors"
	"net/http"

	readerModel "github.com/web-rabis/reader-client/model"
	dto2 "github.com/web-rabis/reader-workstation-api/internal/resource/http/reader/v1/dto"
	"github.com/web-rabis/reader-workstation-api/internal/utils"

	"github.com/go-chi/render"
	"github.com/web-rabis/httperrors"
)

// @Summary Create Reader User
// @Description Create Reader User
// @Accept json
// @Produce json
// @Tags ReaderUser
// @Param body body dto.CreateReaderUserRequest true "Request"
// @Success 200 {object} dto.CreateReaderUserResponse
// @Failure 400 {object} httperrors.Response
// @Failure 401 {object} httperrors.Response
// @Failure 403 {object} httperrors.Response
// @Failure 404 {object} httperrors.Response
// @Failure 500 {object} httperrors.Response
// @Router /reader-user [post]
func (res *Resource) create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var request dto2.CreateReaderRequest
	var err error

	if err = json.NewDecoder(r.Body).Decode(&request); err != nil {
		_ = render.Render(w, r, httperrors.BadRequest(err))
		return
	}

	if err := request.Validate(); err != nil {
		_ = render.Render(w, r, httperrors.BadRequest(err))
		return
	}
	req, err := request.Reader()
	if err := request.Validate(); err != nil {
		_ = render.Render(w, r, httperrors.BadRequest(err))
		return
	}

	filter := &readerModel.ReaderFilter{
		Iin: request.Iin,
	}
	_, readers, err := res.readerMan.List(ctx, filter, nil)
	if err != nil {
		_ = render.Render(w, r, httperrors.Internal(err))
		return
	}
	for _, reader_ := range readers {
		if reader_.Status != "rejected" {
			_ = render.Render(w, r, httperrors.Internal(errors.New("reader already exists")))
			return
		}
	}
	photoBytes, err := utils.ResizeImage(request.Photo, 1024, 1024, utils.ModeFit, 90)
	if err != nil {
		_ = render.Render(w, r, httperrors.Internal(err))
		return
	}

	ticketNumber, err := res.readerMan.Create(ctx, req, photoBytes)
	if err != nil {
		_ = render.Render(w, r, httperrors.Internal(err))
		return
	}
	render.JSON(w, r, ticketNumber)
}
