package http

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ordermodel "github.com/web-rabis/order-client/model"
	readerModel "github.com/web-rabis/reader-client/model"
	"github.com/web-rabis/reader-workstation-api/internal/manager/ebook"
	"github.com/web-rabis/reader-workstation-api/internal/manager/eorder"
	"github.com/web-rabis/reader-workstation-api/internal/manager/reader"
	"github.com/web-rabis/reader-workstation-api/internal/manager/user"
)

const healthCheckTimeout = 3 * time.Second

// HealthResponse - ответ на проверку работоспособности сервиса.
type HealthResponse struct {
	Status       string            `json:"status"`
	Dependencies map[string]string `json:"dependencies"`
}

// HealthResource - проверка работоспособности сервиса: сам сервис и все его gRPC-соединения.
type HealthResource struct {
	path      string
	readerMan reader.IManager
	orderMan  eorder.IManager
	ebookMan  ebook.IManager
	userMan   *user.Manager
}

func NewHealthResource(path string, readerMan reader.IManager, orderMan eorder.IManager, ebookMan ebook.IManager, userMan *user.Manager) *HealthResource {
	return &HealthResource{
		path:      path,
		readerMan: readerMan,
		orderMan:  orderMan,
		ebookMan:  ebookMan,
		userMan:   userMan,
	}
}

func (hr *HealthResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", hr.get)

	return r
}

func (hr *HealthResource) Path() string {
	return hr.path
}

type dependencyCheck struct {
	name string
	err  error
}

// @Summary Health Check
// @Description Проверяет работоспособность сервиса и всех его gRPC-соединений (reader, order, ebook, sso)
// @Produce json
// @Tags Health
// @Success 200 {object} HealthResponse
// @Failure 503 {object} HealthResponse
// @Router /health [get]
func (hr *HealthResource) get(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), healthCheckTimeout)
	defer cancel()

	checks := map[string]func() error{
		"reader": func() error {
			_, _, err := hr.readerMan.List(ctx, &readerModel.ReaderFilter{}, &readerModel.Paging{Limit: 1})
			return err
		},
		"order": func() error {
			_, err := hr.orderMan.Count(ctx, &ordermodel.OrderFilters{})
			return err
		},
		"ebook": func() error {
			_, err := hr.ebookMan.EbookById(ctx, 0)
			return err
		},
		"sso": func() error {
			_, err := hr.userMan.ById(ctx, 0)
			return err
		},
	}

	results := make(chan dependencyCheck, len(checks))
	for name, check := range checks {
		go func(name string, check func() error) {
			results <- dependencyCheck{name: name, err: check()}
		}(name, check)
	}

	deps := make(map[string]string, len(checks))
	healthy := true
	for range checks {
		res := <-results
		if ok, msg := isDependencyAlive(res.err); ok {
			deps[res.name] = "ok"
		} else {
			deps[res.name] = msg
			healthy = false
		}
	}

	response := HealthResponse{Status: "ok", Dependencies: deps}
	if !healthy {
		response.Status = "degraded"
		render.Status(r, http.StatusServiceUnavailable)
	}
	render.JSON(w, r, response)
}

// isDependencyAlive считает соединение рабочим, если сервис ответил хоть
// каким-то gRPC-кодом: "не найдено" или бизнес-ошибка так же подтверждают,
// что соединение живое. Недоступность или таймаут соединения — единственные
// признаки того, что связь со службой разорвана.
func isDependencyAlive(err error) (bool, string) {
	if err == nil {
		return true, ""
	}
	switch status.Code(err) {
	case codes.Unavailable, codes.DeadlineExceeded:
		return false, err.Error()
	default:
		return true, ""
	}
}
