module github.com/web-rabis/reader-workstation-api

go 1.26.0

require (
	github.com/disintegration/imaging v1.6.2
	github.com/go-chi/chi/v5 v5.0.10
	github.com/go-chi/cors v1.2.1
	github.com/go-chi/jwtauth v1.2.0
	github.com/go-chi/render v1.0.3
	github.com/hashicorp/logutils v1.0.0
	github.com/jackc/pgx/v5 v5.10.0
	github.com/jessevdk/go-flags v1.5.0
	github.com/lestrrat-go/jwx v1.1.0
	github.com/web-rabis/db v0.0.0-20261006094859-c55bf2bdbb49
	github.com/web-rabis/ebook-client v0.0.0-20260928182249-b1ffb5e36015
	github.com/web-rabis/httperrors v0.0.0-20230828182436-93be23854738
	github.com/web-rabis/order-client v0.0.0-20260921181821-4305380e0002
	github.com/web-rabis/reader-client v0.0.0-20261013081954-15ff016d105a
	github.com/web-rabis/servers v0.0.0-20230808095603-66994ba613c7
	golang.org/x/sync v0.23.0
	gorm.io/gorm v1.25.4
)

require (
	github.com/ajg/form v1.5.1 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/goccy/go-json v0.3.5 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/lestrrat-go/backoff/v2 v2.0.7 // indirect
	github.com/lestrrat-go/httpcc v1.0.0 // indirect
	github.com/lestrrat-go/iter v1.0.0 // indirect
	github.com/lestrrat-go/option v1.0.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_golang v1.24.1 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.70.1 // indirect
	github.com/prometheus/procfs v0.21.1 // indirect
	github.com/web-rabis/sso-client v0.0.0-20260902072025-6f0894046888 // indirect
	golang.org/x/crypto v0.57.0 // indirect
	golang.org/x/image v0.0.0-20191009234506-e7c1f5e7dbb8 // indirect
	golang.org/x/net v0.60.0 // indirect
	golang.org/x/sys v0.48.0 // indirect
	golang.org/x/text v0.42.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20261005182115-fad411399dd8 // indirect
	google.golang.org/grpc v1.84.0 // indirect
	google.golang.org/protobuf v1.36.12 // indirect
)

replace github.com/web-rabis/reader-client => ../reader-client
