module github.com/kraken-hpc/imageapi

go 1.15

require (
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d // indirect
	github.com/bensallen/rbd v0.0.0-20210224155049-baf486eceefa
	github.com/go-openapi/analysis v0.20.1 // indirect
	github.com/go-openapi/errors v0.20.0
	github.com/go-openapi/jsonreference v0.19.6 // indirect
	github.com/go-openapi/loads v0.20.2
	github.com/go-openapi/runtime v0.19.29
	github.com/go-openapi/spec v0.20.3
	github.com/go-openapi/strfmt v0.20.1
	github.com/go-openapi/swag v0.19.15
	github.com/go-openapi/validate v0.20.2
	github.com/jessevdk/go-flags v1.5.0
	github.com/kraken-hpc/go-fork v0.1.1
	github.com/kraken-hpc/uinit v0.1.1
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/u-root/iscsinl v0.1.0
	github.com/u-root/u-root v7.0.0+incompatible
	go.mongodb.org/mongo-driver v1.7.0 // indirect
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2 // indirect
	golang.org/x/net v0.0.0-20210726213435-c6fcb2dbf985
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c
	golang.org/x/tools v0.1.5 // indirect
	google.golang.org/genproto v0.0.0-20210602131652-f16073e35f0c // indirect
)

replace github.com/u-root/u-root v7.0.0+incompatible => github.com/u-root/u-root v1.0.1-0.20201119150355-04f343dd1922

replace github.com/u-root/iscsinl v0.1.0 => github.com/jlowellwofford/iscsinl v0.1.1-0.20210831201708-527d3dbd2e6d
