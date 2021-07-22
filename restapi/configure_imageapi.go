// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/kraken-hpc/imageapi/internal/api"
	"github.com/kraken-hpc/imageapi/restapi/operations"
	"github.com/kraken-hpc/imageapi/restapi/operations/attachments"
	"github.com/kraken-hpc/imageapi/restapi/operations/containers"
	"github.com/kraken-hpc/imageapi/restapi/operations/mounts"
)

//go:generate swagger generate server --target ../../imageapi --name Imageapi --spec ../swagger.yaml --principal interface{}

func configureFlags(api *operations.ImageapiAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.ImageapiAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	if api.AttachmentsDeleteAttachHandler == nil {
		api.AttachmentsDeleteAttachHandler = attachments.DeleteAttachHandlerFunc(func(params attachments.DeleteAttachParams) middleware.Responder {
			return middleware.NotImplemented("operation attachments.DeleteAttach has not yet been implemented")
		})
	}
	if api.MountsDeleteMountHandler == nil {
		api.MountsDeleteMountHandler = mounts.DeleteMountHandlerFunc(func(params mounts.DeleteMountParams) middleware.Responder {
			return middleware.NotImplemented("operation mounts.DeleteMount has not yet been implemented")
		})
	}
	if api.AttachmentsAttachHandler == nil {
		api.AttachmentsAttachHandler = attachments.AttachHandlerFunc(func(params attachments.AttachParams) middleware.Responder {
			return middleware.NotImplemented("operation attachments.Attach has not yet been implemented")
		})
	}
	if api.ContainersCreateContainerHandler == nil {
		api.ContainersCreateContainerHandler = containers.CreateContainerHandlerFunc(func(params containers.CreateContainerParams) middleware.Responder {
			return middleware.NotImplemented("operation containers.CreateContainer has not yet been implemented")
		})
	}
	if api.ContainersDeleteContainerHandler == nil {
		api.ContainersDeleteContainerHandler = containers.DeleteContainerHandlerFunc(func(params containers.DeleteContainerParams) middleware.Responder {
			return middleware.NotImplemented("operation containers.DeleteContainer has not yet been implemented")
		})
	}
	if api.AttachmentsListAttachmentsHandler == nil {
		api.AttachmentsListAttachmentsHandler = attachments.ListAttachmentsHandlerFunc(func(params attachments.ListAttachmentsParams) middleware.Responder {
			return middleware.NotImplemented("operation attachments.ListAttachments has not yet been implemented")
		})
	}
	if api.ContainersListContainersHandler == nil {
		api.ContainersListContainersHandler = containers.ListContainersHandlerFunc(func(params containers.ListContainersParams) middleware.Responder {
			return middleware.NotImplemented("operation containers.ListContainers has not yet been implemented")
		})
	}
	if api.MountsListMountsHandler == nil {
		api.MountsListMountsHandler = mounts.ListMountsHandlerFunc(func(params mounts.ListMountsParams) middleware.Responder {
			return middleware.NotImplemented("operation mounts.ListMounts has not yet been implemented")
		})
	}
	if api.MountsMountHandler == nil {
		api.MountsMountHandler = mounts.MountHandlerFunc(func(params mounts.MountParams) middleware.Responder {
			return middleware.NotImplemented("operation mounts.Mount has not yet been implemented")
		})
	}
	if api.ContainersSetContainerStateHandler == nil {
		api.ContainersSetContainerStateHandler = containers.SetContainerStateHandlerFunc(func(params containers.SetContainerStateParams) middleware.Responder {
			return middleware.NotImplemented("operation containers.SetContainerState has not yet been implemented")
		})
	}

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
	api.Init()
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
