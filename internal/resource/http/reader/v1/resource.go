package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"

	readerModel "github.com/web-rabis/reader-client/model"
	"github.com/web-rabis/reader-workstation-api/internal/manager/auth"
	"github.com/web-rabis/reader-workstation-api/internal/manager/control"
	"github.com/web-rabis/reader-workstation-api/internal/manager/ebook"
	"github.com/web-rabis/reader-workstation-api/internal/manager/eorder"
	"github.com/web-rabis/reader-workstation-api/internal/manager/reader"
	"github.com/web-rabis/reader-workstation-api/internal/manager/user"
	ssoModel "github.com/web-rabis/sso-client/model"
)

type Resource struct {
	path       string
	authMan    *auth.Manager
	readerMan  reader.IManager
	userMan    *user.Manager
	orderMan   eorder.IManager
	ebookMan   ebook.IManager
	controlMan control.IManager
}

func NewAdminReaderResource(basePath string,
	authMan *auth.Manager,
	readerMan reader.IManager,
	userMan *user.Manager,
	orderMan eorder.IManager,
	ebookMan ebook.IManager,
	controlMan control.IManager) *Resource {
	return &Resource{
		path:       basePath,
		authMan:    authMan,
		readerMan:  readerMan,
		userMan:    userMan,
		orderMan:   orderMan,
		ebookMan:   ebookMan,
		controlMan: controlMan,
	}
}

func (res *Resource) Path() string {
	return res.path
}

func (res *Resource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(res.authMan.JWTAuth()))
		r.Use(auth.NewUserAccessCtx(res.authMan.JWTKey()).ChiMiddleware)
		r.Get("/", res.list)
		r.Get("/{id}", res.byId)
		r.Get("/{id}/photo", res.photo)
		r.Put("/{id}/photo/rotate90", res.photoRotate90)
		r.Put("/{id}/photo/rotate270", res.photoRotate270)
		r.Put("/{id}", res.update)
		r.Put("/{id}/approve", res.approve)
		r.Put("/{id}/reject", res.reject)
		r.Put("/{id}/renewal-registration", res.renewalRegistration)
		r.Post("/", res.create)
		r.Delete("/{id}", res.delete)
		r.Get("/visits/today", res.visitsToday)
		r.Get("/visits/statistics", res.visitsStatistics)
		r.Get("/{id}/visits/status", res.visitsStatus)
		r.Post("/{id}/visits/check-in", res.visitsCheckIn)
		r.Post("/{id}/visits/check-out", res.visitsCheckOut)
	})

	return r
}

func PagingParse(r *http.Request) (*readerModel.Paging, error) {
	paging := &readerModel.Paging{SortKey: "ticket_num", SortVal: -1}

	if limitParam := r.URL.Query().Get("limit"); limitParam != "" {
		limit, err := strconv.Atoi(limitParam)
		if err != nil {
			return nil, err
		}
		paging.Limit = int64(limit)
	}

	if pageParam := r.URL.Query().Get("skip"); pageParam != "" {
		page, err := strconv.Atoi(pageParam)
		if err != nil {
			return nil, err
		}
		skip := page - 1
		if paging.Limit != 0 {
			skip = int(paging.Limit) * skip
		} else {
			skip = skip * 10
		}
		paging.Offset = int64(skip)
	}

	if orderParam := r.URL.Query().Get("order"); orderParam != "" {
		sortVal, err := strconv.Atoi(orderParam)
		if err != nil {
			return nil, err
		}

		if sortVal == 0 {
			sortVal = -1
		}

		//nolint:gosec
		paging.SortVal = int32(sortVal)
	}

	if sortKeyFromQuery := r.URL.Query().Get("orderBy"); sortKeyFromQuery != "" {
		if sortKeyFromQuery == "ticketNumber" {
			sortKeyFromQuery = "ticket_num"
		}
		paging.SortKey = sortKeyFromQuery
	}

	return paging, nil
}
func generateBarcode(ticketNumber int64) (string, error) {
	// формируем 12-значную основу: "22001" + ticketNumber с ведущими нулями до 7 символов
	base := fmt.Sprintf("22001%07d", ticketNumber)

	if len(base) != 12 {
		return "", fmt.Errorf("invalid base length: %s", base)
	}

	var sumOdd, sumEven int
	for i, r := range base {
		digit, err := strconv.Atoi(string(r))
		if err != nil {
			return "", err
		}
		if i%2 == 0 {
			sumOdd += digit // чётные индексы
		} else {
			sumEven += digit // нечётные индексы
		}
	}

	checkDigit := 10 - (sumEven*3+sumOdd)%10
	if checkDigit == 10 {
		checkDigit = 0
	}

	return base + strconv.Itoa(checkDigit), nil
}

func (res *Resource) getReader(r *http.Request) (*readerModel.Reader, error) {
	id := chi.URLParam(r, "id")
	if id == "" {
		return nil, errors.New("id is required")
	}
	_id, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, errors.New("invalid id")
	}
	return res.readerMan.ByTicketNumber(r.Context(), _id)
}
func (res *Resource) getUser(r *http.Request) (*readerModel.User, error) {
	var u *ssoModel.User
	var err error
	if r.Context().Value(auth.UserIDKey) == nil {
		return nil, errors.New("user id is empty")
	}
	userId := int(r.Context().Value(auth.UserIDKey).(float64))
	u, err = res.userMan.ById(r.Context(), userId)
	if err != nil {
		return nil, err
	}
	return &readerModel.User{
		Id:       u.Id,
		Name:     u.Name,
		Username: u.Username,
		Email:    u.Email,
		State:    u.State,
	}, nil
}
