package v1

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/web-rabis/httperrors"
	readerModel "github.com/web-rabis/reader-client/model"
	"github.com/web-rabis/reader-workstation-api/internal/manager/auth"
	"github.com/web-rabis/reader-workstation-api/internal/resource/http/reader/v1/dto"
	"github.com/web-rabis/reader-workstation-api/internal/utils"
)

// @Summary Dictionary Sex
// @Description Dictionary Sex
// @Accept json
// @Produce json
// @Tags Dictionary
// @Success 200 {object} reader.Dictionary
// @Failure 400 {object} httperrors.Response
// @Failure 401 {object} httperrors.Response
// @Failure 403 {object} httperrors.Response
// @Failure 404 {object} httperrors.Response
// @Failure 500 {object} httperrors.Response
// @Router /admin/api/v1/auth/sign-in [get]
func (res *Resource) update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		_ = render.Render(w, r, httperrors.BadRequest(errors.New("id is required")))
		return
	}
	_id, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		_ = render.Render(w, r, httperrors.BadRequest(err))
		return
	}

	var request dto.UpdateReaderRequest

	if err = json.NewDecoder(r.Body).Decode(&request); err != nil {
		_ = render.Render(w, r, httperrors.BadRequest(err))
		return
	}

	reader_, err := res.readerMan.ByTicketNumber(r.Context(), _id)
	if err != nil {
		_ = render.Render(w, r, httperrors.BadRequest(err))
		return
	}

	req, err := request.Reader(reader_)

	if err != nil {
		_ = render.Render(w, r, httperrors.BadRequest(err))
		return
	}
	if req.Barcode == "" && req.Status == "" {
		req.Barcode, err = generateBarcode(_id)
		if err != nil {
			_ = render.Render(w, r, httperrors.Internal(err))
			return
		}
	}

	err = res.readerMan.Update(r.Context(), req, nil)
	if err != nil {
		_ = render.Render(w, r, httperrors.BadRequest(err))
		return
	}
	if request.Photo != "" {
		photoBytes, err := utils.ResizeImage(request.Photo, 1024, 1024, utils.ModeFit, 90)
		if err != nil {
			log.Printf("[ERROR] decode base64 photo error: %s", err)
		} else {
			var username = "unknown"
			if r.Context().Value(auth.UserIDKey) != nil {
				userId := int(r.Context().Value(auth.UserIDKey).(float64))
				user, err := res.userMan.ById(r.Context(), userId)
				if err != nil {
					log.Printf("[ERROR] get user error: %s", err)
				} else {
					username = user.Username
				}
			} else {
				log.Printf("[ERROR] user is required")
			}
			err = res.readerMan.UploadPhoto(r.Context(), _id, photoBytes, &readerModel.User{Username: username})
			if err != nil {
				log.Printf("[ERROR] upload photo error: %s", err)
			}
		}
	}
	render.JSON(w, r, nil)
}
