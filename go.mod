module github.com/jlowellwofford/imageapi

go 1.15

replace github.com/bensallen/rbd v0.0.0-20201123204607-9b994a7784d0 => github.com/jlowellwofford/rbd v0.0.0-20210209155714-708f1b8f8fd7

require (
	github.com/bensallen/rbd v0.0.0-20201123204607-9b994a7784d0
	github.com/go-openapi/errors v0.20.0
	github.com/go-openapi/loads v0.20.2
	github.com/go-openapi/runtime v0.19.26
	github.com/go-openapi/spec v0.20.3
	github.com/go-openapi/strfmt v0.20.0
	github.com/go-openapi/swag v0.19.14
	github.com/go-openapi/validate v0.20.2
	github.com/jessevdk/go-flags v1.4.0
	github.com/jlowellwofford/go-fork v0.0.0-20210209212206-08fa7acf04af
	github.com/jlowellwofford/uinit v0.0.0-20210105192521-a5a9b4f1edb7
	golang.org/x/net v0.0.0-20210119194325-5f4716e94777
	golang.org/x/sys v0.0.0-20201119102817-f84b799fce68
)
