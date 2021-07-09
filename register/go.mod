module github.com/wcgcyx/fc-retrieval/register

go 1.16

replace github.com/wcgcyx/fc-retrieval/common => ../common

require (
	github.com/go-openapi/errors v0.20.0
	github.com/go-openapi/loads v0.20.2
	github.com/go-openapi/runtime v0.19.29
	github.com/go-openapi/spec v0.20.3
	github.com/go-openapi/strfmt v0.20.1
	github.com/go-openapi/swag v0.19.15
	github.com/go-openapi/validate v0.20.2
	github.com/go-redis/redis/v8 v8.11.0
	github.com/jessevdk/go-flags v1.5.0
	github.com/rs/cors v1.8.0
	github.com/spf13/viper v1.8.1
	github.com/wcgcyx/fc-retrieval/common v0.0.0-00010101000000-000000000000
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e
)
