package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hashicorp/logutils"
	"golang.org/x/sync/errgroup"

	ebookCli "github.com/web-rabis/ebook-client"
	ebookClient "github.com/web-rabis/ebook-client/client"
	ebookModel "github.com/web-rabis/ebook-client/model"
	orderCli "github.com/web-rabis/order-client"
	orderClient "github.com/web-rabis/order-client/client"
	orderModel "github.com/web-rabis/order-client/model"
	readerCli "github.com/web-rabis/reader-client"
	readerClient "github.com/web-rabis/reader-client/client"
	readerModel "github.com/web-rabis/reader-client/model"
	"github.com/web-rabis/reader-workstation-api/internal/config"
	"github.com/web-rabis/reader-workstation-api/internal/manager/auth"
	"github.com/web-rabis/reader-workstation-api/internal/manager/control"
	"github.com/web-rabis/reader-workstation-api/internal/manager/dictionary"
	"github.com/web-rabis/reader-workstation-api/internal/manager/ebook"
	"github.com/web-rabis/reader-workstation-api/internal/manager/eorder"
	"github.com/web-rabis/reader-workstation-api/internal/manager/reader"
	"github.com/web-rabis/reader-workstation-api/internal/manager/user"
	"github.com/web-rabis/reader-workstation-api/internal/server/http"
	ssoCli "github.com/web-rabis/sso-client"
	ssoClient "github.com/web-rabis/sso-client/client"
	ssoModel "github.com/web-rabis/sso-client/model"
)

var (
	version = "unknown"
)

func main() {
	fmt.Printf("reader-workstation-api %s\n", version)

	opts := config.Parse(new(config.APIServer)).(*config.APIServer)

	setupLogUtils(opts.Dbg)

	appCtx, appCtxCancel := context.WithCancel(context.Background())
	defer appCtxCancel()

	go catchForTermination(appCtxCancel, os.Interrupt, syscall.SIGTERM)

	orderGrpcClient, ebookGrpcClient, readerGrpcClient, ssoGrpcClient, err := setupGrpcClients(opts)
	if err != nil {
		panic(err)
	}
	defer closeGrpcClients(orderGrpcClient, ebookGrpcClient, readerGrpcClient, ssoGrpcClient)

	authMan, userMan, dictMan, readerMan, orderMan, ebookMan, controlMan := setupManagers(opts, orderGrpcClient, ebookGrpcClient, readerGrpcClient, ssoGrpcClient)
	servers, serversCtx := errgroup.WithContext(appCtx)

	servers.Go(func() error {
		return http.Run(serversCtx, opts, authMan, dictMan, readerMan, orderMan, ebookMan, controlMan, userMan, version)
	})

	if err := servers.Wait(); err != nil {
		log.Printf("[INFO] process terminated, %s", err)
		return
	}

}

func setupLogUtils(inDebugMode bool) {
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel("INFO"),
		Writer:   os.Stdout,
	}

	log.SetFlags(log.Ldate | log.Ltime)

	if inDebugMode {
		log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
		filter.MinLevel = "DEBUG"
	}

	log.SetOutput(filter)
}

func catchForTermination(cancel context.CancelFunc, signals ...os.Signal) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, signals...)
	<-stop
	log.Print("[WARN] interrupt signal")
	cancel()
}
func setupGrpcClients(opts *config.APIServer) (orderClient.Base, ebookClient.Base, readerClient.Base, ssoClient.Base, error) {
	orderGrpcClient, err := orderCli.NewOrderClient(&orderModel.ConnectionConfig{
		Address:  opts.OrderConfig.GrpcAddress,
		Protocol: "grpc",
		Insecure: true,
	})
	if err != nil {
		return nil, nil, nil, nil, err
	}
	log.Printf("[INFO] connected to order grpc-server(%s)", opts.OrderConfig.GrpcAddress)
	ebookGrpcClient, err := ebookCli.NewEbookClient(&ebookModel.ConnectionConfig{
		Address:  opts.EbookConfig.GrpcAddress,
		Protocol: "grpc",
		Insecure: true,
	})
	if err != nil {
		return nil, nil, nil, nil, err
	}
	log.Printf("[INFO] connected to ebook grpc-server(%s)", opts.EbookConfig.GrpcAddress)
	readerGrpcClient, err := readerCli.NewReaderClient(&readerModel.ConnectionConfig{
		Address:  opts.ReaderConfig.GrpcAddress,
		Protocol: "grpc",
		Insecure: true,
	})
	if err != nil {
		return nil, nil, nil, nil, err
	}
	log.Printf("[INFO] connected to sso grpc-server(%s)", opts.SsoConfig.GrpcAddress)
	ssoGrpcClient, err := ssoCli.NewSsoClient(&ssoModel.ConnectionConfig{
		Address:  opts.SsoConfig.GrpcAddress,
		Protocol: "grpc",
		Insecure: true,
	})
	if err != nil {
		return nil, nil, nil, nil, err
	}
	return orderGrpcClient, ebookGrpcClient, readerGrpcClient, ssoGrpcClient, nil
}
func closeGrpcClients(orderGrpcClient orderClient.Base, ebookGrpcClient ebookClient.Base, readerGrpcClient readerClient.Base, ssoGrpcClient ssoClient.Base) {
	if err := orderGrpcClient.Close(); err != nil {
		log.Printf("[WARN] failed to close order grpc client: %v", err)
	}
	if err := ebookGrpcClient.Close(); err != nil {
		log.Printf("[WARN] failed to close ebook grpc client: %v", err)
	}
	if err := readerGrpcClient.Close(); err != nil {
		log.Printf("[WARN] failed to close reader grpc client: %v", err)
	}
	if err := ssoGrpcClient.Close(); err != nil {
		log.Printf("[WARN] failed to close sso grpc client: %v", err)
	}
}
func setupManagers(opts *config.APIServer, orderGrpcClient orderClient.Base, ebookGrpcClient ebookClient.Base, readerGrpcClient readerClient.Base, ssoGrpcClient ssoClient.Base) (*auth.Manager, *user.Manager, dictionary.IManager, reader.IManager, eorder.IManager, ebook.IManager, control.IManager) {
	return auth.NewAuthManager(
			[]byte(opts.JWTKey),
			time.Minute*time.Duration(opts.AccessTokenTTL),
			time.Hour*time.Duration(opts.RefreshTokenTTL)),
		user.NewUserManager(ssoGrpcClient.User()),
		dictionary.NewManager(readerGrpcClient.DictionarySvc()),
		reader.NewManager(readerGrpcClient.ReaderSvc()),

		eorder.NewManager(orderGrpcClient.Order()),
		ebook.NewManager(ebookGrpcClient.EbookSvc()),
		control.NewManager(readerGrpcClient.ControlSvc())

}
