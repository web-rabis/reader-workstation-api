module github.com/web-rabis/reader-workstation-api

go 1.26.0

require (
	github.com/disintegration/imaging v1.6.2
	github.com/go-chi/chi/v5 v5.0.10
	github.com/go-chi/cors v1.2.1
	github.com/go-chi/jwtauth v1.2.0
	github.com/go-chi/render v1.0.3
	github.com/hashicorp/logutils v1.0.0
	github.com/jessevdk/go-flags v1.5.0
	github.com/lestrrat-go/jwx v1.1.0
	github.com/swaggo/http-swagger v1.3.4
	github.com/swaggo/swag v1.16.4
	github.com/web-rabis/ebook-client v0.0.0-20260928182249-b1ffb5e36015
	github.com/web-rabis/httperrors v0.0.0-20230828182436-93be23854738
	github.com/web-rabis/order-client v0.0.0-20260921181821-4305380e0002
	github.com/web-rabis/reader-client v0.0.0-20261013081954-15ff016d105a
	github.com/web-rabis/servers v0.0.0-20230808095603-66994ba613c7
	github.com/web-rabis/sso-client v0.0.0-20260902072025-6f0894046888
	golang.org/x/sync v0.23.0
	google.golang.org/grpc v1.84.0
)

require (
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/ajg/form v1.5.1 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.20.0 // indirect
	github.com/go-openapi/spec v0.20.6 // indirect
	github.com/go-openapi/swag v0.19.15 // indirect
	github.com/goccy/go-json v0.3.5 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/lestrrat-go/backoff/v2 v2.0.7 // indirect
	github.com/lestrrat-go/httpcc v1.0.0 // indirect
	github.com/lestrrat-go/iter v1.0.0 // indirect
	github.com/lestrrat-go/option v1.0.0 // indirect
	github.com/mailru/easyjson v0.7.6 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/swaggo/files v0.0.0-20220610200504-28940afbdbfe // indirect
	golang.org/x/crypto v0.57.0 // indirect
	golang.org/x/image v0.0.0-20191009234506-e7c1f5e7dbb8 // indirect
	golang.org/x/net v0.60.0 // indirect
	golang.org/x/sys v0.48.0 // indirect
	golang.org/x/text v0.42.0 // indirect
	golang.org/x/tools v0.49.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20261005182115-fad411399dd8 // indirect
	google.golang.org/protobuf v1.36.12 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/web-rabis/reader-client => ../reader-client
