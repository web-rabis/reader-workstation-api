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
	ebookCli "github.com/web-rabis/ebook-client"
	ebookModel "github.com/web-rabis/ebook-client/model"
	orderCli "github.com/web-rabis/order-client"
	orderModel "github.com/web-rabis/order-client/model"
	readerCli "github.com/web-rabis/reader-client"
	readerModel "github.com/web-rabis/reader-client/model"
	"github.com/web-rabis/reader-workstation-api/internal/manager/auth"
	"github.com/web-rabis/reader-workstation-api/internal/manager/control"
	"github.com/web-rabis/reader-workstation-api/internal/manager/dictionary"
	"github.com/web-rabis/reader-workstation-api/internal/manager/ebook"
	"github.com/web-rabis/reader-workstation-api/internal/manager/eorder"
	"github.com/web-rabis/reader-workstation-api/internal/manager/reader"
	"github.com/web-rabis/reader-workstation-api/internal/manager/user"
	"github.com/web-rabis/reader-workstation-api/internal/server/http"
	ssoCli "github.com/web-rabis/sso-client"
	ssoModel "github.com/web-rabis/sso-client/model"
	"golang.org/x/sync/errgroup"

	"github.com/web-rabis/reader-workstation-api/internal/config"
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

	orderGrpcClient, err := orderCli.NewOrderClient(&orderModel.ConnectionConfig{
		Address:  opts.OrderConfig.GrpcAddress,
		Protocol: "grpc",
		Insecure: true,
	})
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := orderGrpcClient.Close(); err != nil {
			log.Printf("[WARN] failed to close order grpc client: %v", err)
		}
	}()
	ebookGrpcClient, err := ebookCli.NewEbookClient(&ebookModel.ConnectionConfig{
		Address:  opts.EbookConfig.GrpcAddress,
		Protocol: "grpc",
		Insecure: true,
	})
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := ebookGrpcClient.Close(); err != nil {
			log.Printf("[WARN] failed to close ebook grpc client: %v", err)
		}
	}()
	readerGrpcClient, err := readerCli.NewReaderClient(&readerModel.ConnectionConfig{
		Address:  opts.ReaderConfig.GrpcAddress,
		Protocol: "grpc",
		Insecure: true,
	})
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := readerGrpcClient.Close(); err != nil {
			log.Printf("[WARN] failed to close reader grpc client: %v", err)
		}
	}()
	ssoGrpcClient, err := ssoCli.NewSsoClient(&ssoModel.ConnectionConfig{
		Address:  opts.SsoConfig.GrpcAddress,
		Protocol: "grpc",
		Insecure: true,
	})
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := ssoGrpcClient.Close(); err != nil {
			log.Printf("[WARN] failed to close sso grpc client: %v", err)
		}
	}()
	authMan := auth.NewAuthManager(
		[]byte(opts.JWTKey),
		time.Minute*time.Duration(opts.AccessTokenTTL),
		time.Hour*time.Duration(opts.RefreshTokenTTL))
	userMan := user.NewUserManager(ssoGrpcClient.User())
	dictMan := dictionary.NewManager(readerGrpcClient.DictionarySvc())
	readerMan := reader.NewManager(readerGrpcClient.ReaderSvc())

	orderMan := eorder.NewManager(orderGrpcClient.Order())
	ebookMan := ebook.NewManager(ebookGrpcClient.EbookSvc())
	controlMan := control.NewManager(readerGrpcClient.ControlSvc())
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
