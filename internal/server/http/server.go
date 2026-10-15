package http

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/web-rabis/reader-workstation-api/internal/manager/auth"
	"github.com/web-rabis/reader-workstation-api/internal/manager/control"
	"github.com/web-rabis/reader-workstation-api/internal/manager/dictionary"
	"github.com/web-rabis/reader-workstation-api/internal/manager/ebook"
	"github.com/web-rabis/reader-workstation-api/internal/manager/eorder"
	"github.com/web-rabis/reader-workstation-api/internal/manager/reader"
	"github.com/web-rabis/reader-workstation-api/internal/manager/user"
	"github.com/web-rabis/reader-workstation-api/internal/resource/http"
	v5 "github.com/web-rabis/reader-workstation-api/internal/resource/http/dictionary/v1"
	v4 "github.com/web-rabis/reader-workstation-api/internal/resource/http/reader/v1"

	"github.com/web-rabis/reader-workstation-api/internal/config"
	cherver "github.com/web-rabis/servers/http"
)

const (
	maxAge        = 300
	compressLevel = 5
)

func Run(serversCtx context.Context, opts *config.APIServer,
	authMan *auth.Manager,
	dictMan dictionary.IManager,
	readerMan reader.IManager,
	orderMan eorder.IManager,
	ebookMan ebook.IManager,
	controlMan control.IManager,
	userMan *user.Manager,
	version string) error {
	resources := []cherver.Resource{
		http.NewVersionResource("/version", version),
		http.NewFilesResource("/files", opts.ServerConfig.FilesDir),
		v4.NewAdminReaderResource("/api/v1/reader", authMan, readerMan, userMan, orderMan, ebookMan, controlMan),
		v5.NewAdminDictionaryResource("/api/v1/dictionary", authMan, dictMan),
		//swaggerV1.NewSwaggerResource("/swagger", opts.ServerConfig.BasePath, "/files"),
		//signV1.NewSignResource("/sign", signMan),
	}
	httpSrv := cherver.New(
		cherver.WithListenAddress(opts.ServerConfig.ListenAddr),
		cherver.WithCert(opts.ServerConfig.CertFile, opts.ServerConfig.KeyFile),
		cherver.WithResources(resources...),
		cherver.WithMiddlewares(middlewares(opts)...))

	return httpSrv.Run(serversCtx)

}

func middlewaresWithoutLogs(opts *config.APIServer) chi.Middlewares {
	return chi.Middlewares{
		middleware.NoCache,   // no-cache
		middleware.Recoverer, // управляемо обрабатывает паники и выдает stack trace при их возникновении
		middleware.RealIP,    // устанавливает RemoteAddr для каждого запроса с заголовками X-Forwarded-For или X-Real-IP
		middleware.NewCompressor(compressLevel).Handler,

		cors.Handler(cors.Options{
			AllowedOrigins:   allowedOrigins(opts.IsTesting),
			AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: false,
			MaxAge:           maxAge, // Maximum value not ignored by any of major browsers
		})}
}

func middlewares(opts *config.APIServer) chi.Middlewares {
	return append(middlewaresWithoutLogs(opts), middleware.Logger)
}
func allowedOrigins(testing bool) []string {
	if testing {
		return []string{"*"}
	}

	return []string{}
}
