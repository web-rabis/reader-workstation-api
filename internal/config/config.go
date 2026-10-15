package config

import (
	"log"
	"os"

	"github.com/jessevdk/go-flags"
)

type ServerConfig struct {
	ListenAddr string `short:"l" long:"listen" env:"LISTEN" description:"Listen Address (format: :8080|127.0.0.1:8080)" required:"false" default:":9191"`
	BasePath   string `long:"base-path" env:"BASE_PATH" description:"base path of the host" required:"false" default:"/reader-workstation-api"`
	FilesDir   string `long:"files-directory" env:"FILES_DIR" description:"Directory where all static files are located" required:"false" default:"/usr/share/api"`
	CertFile   string `short:"c" long:"cert" env:"CERT_FILE" description:"Location of the SSL/TLS cert file" required:"false" default:""`
	KeyFile    string `short:"k" long:"key" env:"KEY_FILE" description:"Location of the SSL/TLS key file" required:"false" default:""`

	GrpcListenAddr    string `long:"grpc-listen" env:"GRPC_LISTEN" description:"Grpc Listen Address (format: :4000|127.0.0.1:4000)" required:"false" default:":4000"`
	MetricsListenAddr string `long:"metrics-listen" env:"METRICS_LISTEN" description:"Metrics Listen Address (format: :5000|127.0.0.1:5000)" required:"false" default:":4040"`
}
type AuthConfig struct {
	JWTKey          string `long:"jwt-key" env:"JWT_KEY" description:"JWT secret key" required:"false" default:"nlrk-secret"`
	AccessTokenTTL  int    `long:"access-token-ttl" env:"ACCESS_TOKEN_TTL" required:"false" default:"60"`   // Access token expiration in minutes
	RefreshTokenTTL int    `long:"refresh-token-ttl" env:"REFRESH_TOKEN_TTL" required:"false" default:"24"` // Refresh token expiration in hours
}
type GeneralConfig struct {
	Dbg       bool `long:"dbg" env:"DEBUG" description:"debug mode"`
	IsTesting bool `long:"testing" env:"APP_TESTING" description:"testing mode"`
}
type ReaderConfig struct {
	GrpcAddress string `long:"reader-grpc-address" env:"READER_GRPC_ADDRESS" description:"Reader Grpc Address (format: :4000|127.0.0.1:4000)" required:"false" default:":4000"`
}
type OrderConfig struct {
	GrpcAddress string `long:"order-grpc-address" env:"ORDER_GRPC_ADDRESS" description:"Order Grpc Address (format: :4000|127.0.0.1:4000)" required:"false" default:":4000"`
}
type EbookConfig struct {
	GrpcAddress string `long:"ebook-grpc-address" env:"EBOOK_GRPC_ADDRESS" description:"Ebook Grpc Address (format: :4000|127.0.0.1:4000)" required:"false" default:":4000"`
}
type SsoConfig struct {
	GrpcAddress string `long:"sso-grpc-address" env:"SSO_GRPC_ADDRESS" description:"Sso Grpc Address (format: :4000|127.0.0.1:4000)" required:"false" default:":4000"`
}

func Parse(c interface{}) interface{} {
	p := flags.NewParser(c, flags.Default)

	if _, err := p.Parse(); err != nil {
		log.Println("[ERROR] Error while parsing config options:", err)
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	return c
}
