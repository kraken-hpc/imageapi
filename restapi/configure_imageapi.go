// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"

	internal "github.com/jlowellwofford/imageapi/internal/api"
	"github.com/jlowellwofford/imageapi/models"
	"github.com/jlowellwofford/imageapi/restapi/operations"
	"github.com/jlowellwofford/imageapi/restapi/operations/attach"
	"github.com/jlowellwofford/imageapi/restapi/operations/containers"
	"github.com/jlowellwofford/imageapi/restapi/operations/mounts"
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

	////////////////////////////////
	// Begin: Non-generated block /
	//////////////////////////////

	api.AttachListRbdsHandler = attach.ListRbdsHandlerFunc(func(params attach.ListRbdsParams) middleware.Responder {
		return attach.NewListRbdsOK().WithPayload(internal.Rbds.List())
	})

	api.AttachMapRbdHandler = attach.MapRbdHandlerFunc(func(params attach.MapRbdParams) middleware.Responder {
		var err error
		var r *models.Rbd
		if r, err = internal.Rbds.Map(params.Rbd); err != nil {
			return attach.NewMapRbdDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
		}
		return attach.NewMapRbdCreated().WithPayload(r)
	})

	api.AttachGetRbdHandler = attach.GetRbdHandlerFunc(func(params attach.GetRbdParams) middleware.Responder {
		var err error
		var r *models.Rbd
		if r, err = internal.Rbds.Get(models.ID(params.ID)); err != nil {
			return attach.NewGetRbdDefault(404).WithPayload(&models.Error{Code: 404, Message: swag.String("rbd not found")})
		}
		return attach.NewGetRbdOK().WithPayload(r)
	})

	api.AttachUnmapRbdHandler = attach.UnmapRbdHandlerFunc(func(params attach.UnmapRbdParams) middleware.Responder {
		var err error
		var r *models.Rbd
		if r, err = internal.Rbds.Unmap(models.ID(params.ID)); err != nil {
			if err == internal.ERRNOTFOUND {
				return attach.NewUnmapRbdDefault(404).WithPayload(&models.Error{Code: 404, Message: swag.String("rbd not found")})
			}
			return attach.NewUnmapRbdDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
		}
		return attach.NewUnmapRbdOK().WithPayload(r)
	})

	// MountsRbd

	api.MountsListMountsRbdHandler = mounts.ListMountsRbdHandlerFunc(func(params mounts.ListMountsRbdParams) middleware.Responder {
		return mounts.NewListMountsRbdOK().WithPayload(internal.MountsRbd.List())
	})

	api.MountsMountRbdHandler = mounts.MountRbdHandlerFunc(func(params mounts.MountRbdParams) middleware.Responder {
		var err error
		var r *models.MountRbd
		if r, err = internal.MountsRbd.Mount(params.Mount); err != nil {
			return mounts.NewMountRbdDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
		}
		return mounts.NewMountRbdCreated().WithPayload(r)
	})

	api.MountsGetMountRbdHandler = mounts.GetMountRbdHandlerFunc(func(params mounts.GetMountRbdParams) middleware.Responder {
		var err error
		var r *models.MountRbd
		if r, err = internal.MountsRbd.Get(models.ID(params.ID)); err != nil {
			return mounts.NewGetMountRbdDefault(404).WithPayload(&models.Error{Code: 404, Message: swag.String(err.Error())})
		}
		return mounts.NewGetMountRbdOK().WithPayload(r)
	})

	api.MountsUnmountRbdHandler = mounts.UnmountRbdHandlerFunc(func(params mounts.UnmountRbdParams) middleware.Responder {
		var err error
		var r *models.MountRbd
		if r, err = internal.MountsRbd.Unmount(models.ID(params.ID)); err != nil {
			if err == internal.ERRNOTFOUND {
				return mounts.NewUnmountRbdDefault(404).WithPayload(&models.Error{Code: 404, Message: swag.String("mount not found")})
			}
			return mounts.NewUnmountRbdDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
		}
		return mounts.NewUnmountRbdOK().WithPayload(r)
	})

	// MountsOverlay

	api.MountsListMountsOverlayHandler = mounts.ListMountsOverlayHandlerFunc(func(params mounts.ListMountsOverlayParams) middleware.Responder {
		return mounts.NewListMountsOverlayOK().WithPayload(internal.MountsOverlay.List())
	})

	api.MountsMountOverlayHandler = mounts.MountOverlayHandlerFunc(func(params mounts.MountOverlayParams) middleware.Responder {
		var err error
		var r *models.MountOverlay
		if r, err = internal.MountsOverlay.Mount(params.Mount); err != nil {
			return mounts.NewMountOverlayDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
		}
		return mounts.NewMountOverlayCreated().WithPayload(r)
	})

	api.MountsGetMountOverlayHandler = mounts.GetMountOverlayHandlerFunc(func(params mounts.GetMountOverlayParams) middleware.Responder {
		var err error
		var r *models.MountOverlay
		if r, err = internal.MountsOverlay.Get(models.ID(params.ID)); err != nil {
			return mounts.NewGetMountOverlayDefault(404).WithPayload(&models.Error{Code: 404, Message: swag.String("mount not found")})
		}
		return mounts.NewGetMountOverlayOK().WithPayload(r)
	})

	api.MountsUnmountOverlayHandler = mounts.UnmountOverlayHandlerFunc(func(params mounts.UnmountOverlayParams) middleware.Responder {
		var err error
		var r *models.MountOverlay
		if r, err = internal.MountsOverlay.Unmount(models.ID(params.ID)); err != nil {
			if err == internal.ERRNOTFOUND {
				return mounts.NewUnmountOverlayDefault(404).WithPayload(&models.Error{Code: 404, Message: swag.String("mount not found")})
			}
			return mounts.NewUnmountOverlayDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
		}
		return mounts.NewUnmountOverlayOK().WithPayload(r)
	})

	// Containers
	api.ContainersCreateContainerHandler = containers.CreateContainerHandlerFunc(func(params containers.CreateContainerParams) middleware.Responder {
		var ctn *models.Container
		var err error
		if ctn, err = internal.Containers.Create(params.Container); err != nil {
			return containers.NewCreateContainerDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
		}
		return containers.NewCreateContainerCreated().WithPayload(ctn)
	})

	api.ContainersDeleteContainerHandler = containers.DeleteContainerHandlerFunc(func(params containers.DeleteContainerParams) middleware.Responder {
		var ctn *models.Container
		var err error
		if ctn, err = internal.Containers.Delete(models.ID(params.ID)); err != nil {
			if err == internal.ERRNOTFOUND {
				return containers.NewDeleteContainerDefault(404).WithPayload(&models.Error{Code: 404, Message: swag.String("container not found")})
			}
			return containers.NewDeleteContainerDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
		}
		return containers.NewDeleteContainerOK().WithPayload(ctn)
	})

	api.ContainersGetContainerHandler = containers.GetContainerHandlerFunc(func(params containers.GetContainerParams) middleware.Responder {
		var ctn *models.Container
		var err error
		if ctn, err = internal.Containers.Get(models.ID(params.ID)); err != nil {
			return containers.NewGetContainerDefault(404).WithPayload(&models.Error{Code: 404, Message: swag.String("container not found")})
		}
		return containers.NewGetContainerOK().WithPayload(ctn)
	})

	api.ContainersListContainersHandler = containers.ListContainersHandlerFunc(func(params containers.ListContainersParams) middleware.Responder {
		return containers.NewListContainersOK().WithPayload(internal.Containers.List())
	})

	api.ContainersSetContainerStateHandler = containers.SetContainerStateHandlerFunc(func(params containers.SetContainerStateParams) middleware.Responder {
		if err := internal.Containers.SetState(models.ID(params.ID), models.ContainerState(params.State)); err != nil {
			if err == internal.ERRNOTFOUND {
				return containers.NewSetContainerStateDefault(404).WithPayload(&models.Error{Code: 404, Message: swag.String("container not found")})
			}
			return containers.NewSetContainerStateDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
		}
		ctn, _ := internal.Containers.Get(models.ID(params.ID))
		return containers.NewSetContainerStateOK().WithPayload(ctn)
	})

	// containers byname
	api.ContainersGetContainerBynameHandler = containers.GetContainerBynameHandlerFunc(func(params containers.GetContainerBynameParams) middleware.Responder {
		id := internal.Containers.NameGetID(models.Name(params.Name))
		if id < 0 {
			return containers.NewGetContainerBynameDefault(404).WithPayload(&models.Error{Code: 404, Message: swag.String("no container by name: " + params.Name)})
		}
		return api.ContainersGetContainerHandler.Handle(containers.GetContainerParams{HTTPRequest: params.HTTPRequest, ID: int64(id)})
	})

	api.ContainersDeleteContainerBynameHandler = containers.DeleteContainerBynameHandlerFunc(func(params containers.DeleteContainerBynameParams) middleware.Responder {
		id := internal.Containers.NameGetID(models.Name(params.Name))
		if id < 0 {
			return containers.NewDeleteContainerBynameDefault(404).WithPayload(&models.Error{Code: 404, Message: swag.String("no container by name: " + params.Name)})
		}
		return api.ContainersDeleteContainerHandler.Handle(containers.DeleteContainerParams{HTTPRequest: params.HTTPRequest, ID: int64(id)})
	})

	api.ContainersSetContainerStateBynameHandler = containers.SetContainerStateBynameHandlerFunc(func(params containers.SetContainerStateBynameParams) middleware.Responder {
		id := internal.Containers.NameGetID(models.Name(params.Name))
		if id < 0 {
			return containers.NewSetContainerStateBynameDefault(404).WithPayload(&models.Error{Code: 404, Message: swag.String("no container by name: " + params.Name)})
		}
		return api.ContainersSetContainerStateHandler.Handle(containers.SetContainerStateParams{HTTPRequest: params.HTTPRequest, State: params.State, ID: int64(id)})
	})

	// generic mounts

	api.MountsListMountsHandler = mounts.ListMountsHandlerFunc(func(params mounts.ListMountsParams) middleware.Responder {
		return mounts.NewListMountsOK().WithPayload(internal.ListMounts())
	})

	api.MountsMountHandler = mounts.MountHandlerFunc(func(params mounts.MountParams) middleware.Responder {
		var mnt *models.Mount
		var err error
		if mnt, err = internal.Mount(params.Mount); err != nil {
			return mounts.NewMountDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
		}
		return mounts.NewMountCreated().WithPayload(mnt)
	})

	api.MountsDeleteMountHandler = mounts.DeleteMountHandlerFunc(func(params mounts.DeleteMountParams) middleware.Responder {
		var mnt *models.Mount
		var err error
		if mnt, err = internal.Unmount(params.Mount); err != nil {
			return mounts.NewDeleteMountDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
		}
		return mounts.NewDeleteMountOK().WithPayload(mnt)
	})

	//////////////////////////////
	// End: Non-generated block /
	////////////////////////////

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
