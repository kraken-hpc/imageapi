// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/runtime/security"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/kraken-hpc/imageapi/restapi/operations/attach"
	"github.com/kraken-hpc/imageapi/restapi/operations/containers"
	"github.com/kraken-hpc/imageapi/restapi/operations/mounts"
)

// NewImageapiAPI creates a new Imageapi instance
func NewImageapiAPI(spec *loads.Document) *ImageapiAPI {
	return &ImageapiAPI{
		handlers:            make(map[string]map[string]http.Handler),
		formats:             strfmt.Default,
		defaultConsumes:     "application/json",
		defaultProduces:     "application/json",
		customConsumers:     make(map[string]runtime.Consumer),
		customProducers:     make(map[string]runtime.Producer),
		PreServerShutdown:   func() {},
		ServerShutdown:      func() {},
		spec:                spec,
		useSwaggerUI:        false,
		ServeError:          errors.ServeError,
		BasicAuthenticator:  security.BasicAuth,
		APIKeyAuthenticator: security.APIKeyAuth,
		BearerAuthenticator: security.BearerAuth,

		JSONConsumer: runtime.JSONConsumer(),

		JSONProducer: runtime.JSONProducer(),

		MountsDeleteMountHandler: mounts.DeleteMountHandlerFunc(func(params mounts.DeleteMountParams) middleware.Responder {
			return middleware.NotImplemented("operation mounts.DeleteMount has not yet been implemented")
		}),
		ContainersCreateContainerHandler: containers.CreateContainerHandlerFunc(func(params containers.CreateContainerParams) middleware.Responder {
			return middleware.NotImplemented("operation containers.CreateContainer has not yet been implemented")
		}),
		ContainersDeleteContainerHandler: containers.DeleteContainerHandlerFunc(func(params containers.DeleteContainerParams) middleware.Responder {
			return middleware.NotImplemented("operation containers.DeleteContainer has not yet been implemented")
		}),
		ContainersDeleteContainerBynameHandler: containers.DeleteContainerBynameHandlerFunc(func(params containers.DeleteContainerBynameParams) middleware.Responder {
			return middleware.NotImplemented("operation containers.DeleteContainerByname has not yet been implemented")
		}),
		ContainersGetContainerHandler: containers.GetContainerHandlerFunc(func(params containers.GetContainerParams) middleware.Responder {
			return middleware.NotImplemented("operation containers.GetContainer has not yet been implemented")
		}),
		ContainersGetContainerBynameHandler: containers.GetContainerBynameHandlerFunc(func(params containers.GetContainerBynameParams) middleware.Responder {
			return middleware.NotImplemented("operation containers.GetContainerByname has not yet been implemented")
		}),
		MountsGetMountOverlayHandler: mounts.GetMountOverlayHandlerFunc(func(params mounts.GetMountOverlayParams) middleware.Responder {
			return middleware.NotImplemented("operation mounts.GetMountOverlay has not yet been implemented")
		}),
		MountsGetMountRbdHandler: mounts.GetMountRbdHandlerFunc(func(params mounts.GetMountRbdParams) middleware.Responder {
			return middleware.NotImplemented("operation mounts.GetMountRbd has not yet been implemented")
		}),
		AttachGetRbdHandler: attach.GetRbdHandlerFunc(func(params attach.GetRbdParams) middleware.Responder {
			return middleware.NotImplemented("operation attach.GetRbd has not yet been implemented")
		}),
		ContainersListContainersHandler: containers.ListContainersHandlerFunc(func(params containers.ListContainersParams) middleware.Responder {
			return middleware.NotImplemented("operation containers.ListContainers has not yet been implemented")
		}),
		MountsListMountsHandler: mounts.ListMountsHandlerFunc(func(params mounts.ListMountsParams) middleware.Responder {
			return middleware.NotImplemented("operation mounts.ListMounts has not yet been implemented")
		}),
		MountsListMountsOverlayHandler: mounts.ListMountsOverlayHandlerFunc(func(params mounts.ListMountsOverlayParams) middleware.Responder {
			return middleware.NotImplemented("operation mounts.ListMountsOverlay has not yet been implemented")
		}),
		MountsListMountsRbdHandler: mounts.ListMountsRbdHandlerFunc(func(params mounts.ListMountsRbdParams) middleware.Responder {
			return middleware.NotImplemented("operation mounts.ListMountsRbd has not yet been implemented")
		}),
		AttachListRbdsHandler: attach.ListRbdsHandlerFunc(func(params attach.ListRbdsParams) middleware.Responder {
			return middleware.NotImplemented("operation attach.ListRbds has not yet been implemented")
		}),
		AttachMapRbdHandler: attach.MapRbdHandlerFunc(func(params attach.MapRbdParams) middleware.Responder {
			return middleware.NotImplemented("operation attach.MapRbd has not yet been implemented")
		}),
		MountsMountHandler: mounts.MountHandlerFunc(func(params mounts.MountParams) middleware.Responder {
			return middleware.NotImplemented("operation mounts.Mount has not yet been implemented")
		}),
		MountsMountOverlayHandler: mounts.MountOverlayHandlerFunc(func(params mounts.MountOverlayParams) middleware.Responder {
			return middleware.NotImplemented("operation mounts.MountOverlay has not yet been implemented")
		}),
		MountsMountRbdHandler: mounts.MountRbdHandlerFunc(func(params mounts.MountRbdParams) middleware.Responder {
			return middleware.NotImplemented("operation mounts.MountRbd has not yet been implemented")
		}),
		ContainersSetContainerStateHandler: containers.SetContainerStateHandlerFunc(func(params containers.SetContainerStateParams) middleware.Responder {
			return middleware.NotImplemented("operation containers.SetContainerState has not yet been implemented")
		}),
		ContainersSetContainerStateBynameHandler: containers.SetContainerStateBynameHandlerFunc(func(params containers.SetContainerStateBynameParams) middleware.Responder {
			return middleware.NotImplemented("operation containers.SetContainerStateByname has not yet been implemented")
		}),
		AttachUnmapRbdHandler: attach.UnmapRbdHandlerFunc(func(params attach.UnmapRbdParams) middleware.Responder {
			return middleware.NotImplemented("operation attach.UnmapRbd has not yet been implemented")
		}),
		MountsUnmountOverlayHandler: mounts.UnmountOverlayHandlerFunc(func(params mounts.UnmountOverlayParams) middleware.Responder {
			return middleware.NotImplemented("operation mounts.UnmountOverlay has not yet been implemented")
		}),
		MountsUnmountRbdHandler: mounts.UnmountRbdHandlerFunc(func(params mounts.UnmountRbdParams) middleware.Responder {
			return middleware.NotImplemented("operation mounts.UnmountRbd has not yet been implemented")
		}),
	}
}

/*ImageapiAPI This API specification describes a service for attaching, mounting and preparing container images and manipulating those containers.

In general, higher level objects can either reference lower level objects (e.g. a mount referencing an attachment point) by a reference ID,
or, they can contain the full specification of those lower objects.

If an object references another by ID, deletion of that object does not effect the underlying object.

If an object defines a lower level object, that lower level object will automatically be deleted on deletion of the higher level object.

For instance, if a container contains all of the defintions for all mount points and attachments, deletion of the container will automatically unmount
and detach those lower objects.
*/
type ImageapiAPI struct {
	spec            *loads.Document
	context         *middleware.Context
	handlers        map[string]map[string]http.Handler
	formats         strfmt.Registry
	customConsumers map[string]runtime.Consumer
	customProducers map[string]runtime.Producer
	defaultConsumes string
	defaultProduces string
	Middleware      func(middleware.Builder) http.Handler
	useSwaggerUI    bool

	// BasicAuthenticator generates a runtime.Authenticator from the supplied basic auth function.
	// It has a default implementation in the security package, however you can replace it for your particular usage.
	BasicAuthenticator func(security.UserPassAuthentication) runtime.Authenticator
	// APIKeyAuthenticator generates a runtime.Authenticator from the supplied token auth function.
	// It has a default implementation in the security package, however you can replace it for your particular usage.
	APIKeyAuthenticator func(string, string, security.TokenAuthentication) runtime.Authenticator
	// BearerAuthenticator generates a runtime.Authenticator from the supplied bearer token auth function.
	// It has a default implementation in the security package, however you can replace it for your particular usage.
	BearerAuthenticator func(string, security.ScopedTokenAuthentication) runtime.Authenticator

	// JSONConsumer registers a consumer for the following mime types:
	//   - application/json
	JSONConsumer runtime.Consumer

	// JSONProducer registers a producer for the following mime types:
	//   - application/json
	JSONProducer runtime.Producer

	// MountsDeleteMountHandler sets the operation handler for the delete mount operation
	MountsDeleteMountHandler mounts.DeleteMountHandler
	// ContainersCreateContainerHandler sets the operation handler for the create container operation
	ContainersCreateContainerHandler containers.CreateContainerHandler
	// ContainersDeleteContainerHandler sets the operation handler for the delete container operation
	ContainersDeleteContainerHandler containers.DeleteContainerHandler
	// ContainersDeleteContainerBynameHandler sets the operation handler for the delete container byname operation
	ContainersDeleteContainerBynameHandler containers.DeleteContainerBynameHandler
	// ContainersGetContainerHandler sets the operation handler for the get container operation
	ContainersGetContainerHandler containers.GetContainerHandler
	// ContainersGetContainerBynameHandler sets the operation handler for the get container byname operation
	ContainersGetContainerBynameHandler containers.GetContainerBynameHandler
	// MountsGetMountOverlayHandler sets the operation handler for the get mount overlay operation
	MountsGetMountOverlayHandler mounts.GetMountOverlayHandler
	// MountsGetMountRbdHandler sets the operation handler for the get mount rbd operation
	MountsGetMountRbdHandler mounts.GetMountRbdHandler
	// AttachGetRbdHandler sets the operation handler for the get rbd operation
	AttachGetRbdHandler attach.GetRbdHandler
	// ContainersListContainersHandler sets the operation handler for the list containers operation
	ContainersListContainersHandler containers.ListContainersHandler
	// MountsListMountsHandler sets the operation handler for the list mounts operation
	MountsListMountsHandler mounts.ListMountsHandler
	// MountsListMountsOverlayHandler sets the operation handler for the list mounts overlay operation
	MountsListMountsOverlayHandler mounts.ListMountsOverlayHandler
	// MountsListMountsRbdHandler sets the operation handler for the list mounts rbd operation
	MountsListMountsRbdHandler mounts.ListMountsRbdHandler
	// AttachListRbdsHandler sets the operation handler for the list rbds operation
	AttachListRbdsHandler attach.ListRbdsHandler
	// AttachMapRbdHandler sets the operation handler for the map rbd operation
	AttachMapRbdHandler attach.MapRbdHandler
	// MountsMountHandler sets the operation handler for the mount operation
	MountsMountHandler mounts.MountHandler
	// MountsMountOverlayHandler sets the operation handler for the mount overlay operation
	MountsMountOverlayHandler mounts.MountOverlayHandler
	// MountsMountRbdHandler sets the operation handler for the mount rbd operation
	MountsMountRbdHandler mounts.MountRbdHandler
	// ContainersSetContainerStateHandler sets the operation handler for the set container state operation
	ContainersSetContainerStateHandler containers.SetContainerStateHandler
	// ContainersSetContainerStateBynameHandler sets the operation handler for the set container state byname operation
	ContainersSetContainerStateBynameHandler containers.SetContainerStateBynameHandler
	// AttachUnmapRbdHandler sets the operation handler for the unmap rbd operation
	AttachUnmapRbdHandler attach.UnmapRbdHandler
	// MountsUnmountOverlayHandler sets the operation handler for the unmount overlay operation
	MountsUnmountOverlayHandler mounts.UnmountOverlayHandler
	// MountsUnmountRbdHandler sets the operation handler for the unmount rbd operation
	MountsUnmountRbdHandler mounts.UnmountRbdHandler
	// ServeError is called when an error is received, there is a default handler
	// but you can set your own with this
	ServeError func(http.ResponseWriter, *http.Request, error)

	// PreServerShutdown is called before the HTTP(S) server is shutdown
	// This allows for custom functions to get executed before the HTTP(S) server stops accepting traffic
	PreServerShutdown func()

	// ServerShutdown is called when the HTTP(S) server is shut down and done
	// handling all active connections and does not accept connections any more
	ServerShutdown func()

	// Custom command line argument groups with their descriptions
	CommandLineOptionsGroups []swag.CommandLineOptionsGroup

	// User defined logger function.
	Logger func(string, ...interface{})
}

// UseRedoc for documentation at /docs
func (o *ImageapiAPI) UseRedoc() {
	o.useSwaggerUI = false
}

// UseSwaggerUI for documentation at /docs
func (o *ImageapiAPI) UseSwaggerUI() {
	o.useSwaggerUI = true
}

// SetDefaultProduces sets the default produces media type
func (o *ImageapiAPI) SetDefaultProduces(mediaType string) {
	o.defaultProduces = mediaType
}

// SetDefaultConsumes returns the default consumes media type
func (o *ImageapiAPI) SetDefaultConsumes(mediaType string) {
	o.defaultConsumes = mediaType
}

// SetSpec sets a spec that will be served for the clients.
func (o *ImageapiAPI) SetSpec(spec *loads.Document) {
	o.spec = spec
}

// DefaultProduces returns the default produces media type
func (o *ImageapiAPI) DefaultProduces() string {
	return o.defaultProduces
}

// DefaultConsumes returns the default consumes media type
func (o *ImageapiAPI) DefaultConsumes() string {
	return o.defaultConsumes
}

// Formats returns the registered string formats
func (o *ImageapiAPI) Formats() strfmt.Registry {
	return o.formats
}

// RegisterFormat registers a custom format validator
func (o *ImageapiAPI) RegisterFormat(name string, format strfmt.Format, validator strfmt.Validator) {
	o.formats.Add(name, format, validator)
}

// Validate validates the registrations in the ImageapiAPI
func (o *ImageapiAPI) Validate() error {
	var unregistered []string

	if o.JSONConsumer == nil {
		unregistered = append(unregistered, "JSONConsumer")
	}

	if o.JSONProducer == nil {
		unregistered = append(unregistered, "JSONProducer")
	}

	if o.MountsDeleteMountHandler == nil {
		unregistered = append(unregistered, "mounts.DeleteMountHandler")
	}
	if o.ContainersCreateContainerHandler == nil {
		unregistered = append(unregistered, "containers.CreateContainerHandler")
	}
	if o.ContainersDeleteContainerHandler == nil {
		unregistered = append(unregistered, "containers.DeleteContainerHandler")
	}
	if o.ContainersDeleteContainerBynameHandler == nil {
		unregistered = append(unregistered, "containers.DeleteContainerBynameHandler")
	}
	if o.ContainersGetContainerHandler == nil {
		unregistered = append(unregistered, "containers.GetContainerHandler")
	}
	if o.ContainersGetContainerBynameHandler == nil {
		unregistered = append(unregistered, "containers.GetContainerBynameHandler")
	}
	if o.MountsGetMountOverlayHandler == nil {
		unregistered = append(unregistered, "mounts.GetMountOverlayHandler")
	}
	if o.MountsGetMountRbdHandler == nil {
		unregistered = append(unregistered, "mounts.GetMountRbdHandler")
	}
	if o.AttachGetRbdHandler == nil {
		unregistered = append(unregistered, "attach.GetRbdHandler")
	}
	if o.ContainersListContainersHandler == nil {
		unregistered = append(unregistered, "containers.ListContainersHandler")
	}
	if o.MountsListMountsHandler == nil {
		unregistered = append(unregistered, "mounts.ListMountsHandler")
	}
	if o.MountsListMountsOverlayHandler == nil {
		unregistered = append(unregistered, "mounts.ListMountsOverlayHandler")
	}
	if o.MountsListMountsRbdHandler == nil {
		unregistered = append(unregistered, "mounts.ListMountsRbdHandler")
	}
	if o.AttachListRbdsHandler == nil {
		unregistered = append(unregistered, "attach.ListRbdsHandler")
	}
	if o.AttachMapRbdHandler == nil {
		unregistered = append(unregistered, "attach.MapRbdHandler")
	}
	if o.MountsMountHandler == nil {
		unregistered = append(unregistered, "mounts.MountHandler")
	}
	if o.MountsMountOverlayHandler == nil {
		unregistered = append(unregistered, "mounts.MountOverlayHandler")
	}
	if o.MountsMountRbdHandler == nil {
		unregistered = append(unregistered, "mounts.MountRbdHandler")
	}
	if o.ContainersSetContainerStateHandler == nil {
		unregistered = append(unregistered, "containers.SetContainerStateHandler")
	}
	if o.ContainersSetContainerStateBynameHandler == nil {
		unregistered = append(unregistered, "containers.SetContainerStateBynameHandler")
	}
	if o.AttachUnmapRbdHandler == nil {
		unregistered = append(unregistered, "attach.UnmapRbdHandler")
	}
	if o.MountsUnmountOverlayHandler == nil {
		unregistered = append(unregistered, "mounts.UnmountOverlayHandler")
	}
	if o.MountsUnmountRbdHandler == nil {
		unregistered = append(unregistered, "mounts.UnmountRbdHandler")
	}

	if len(unregistered) > 0 {
		return fmt.Errorf("missing registration: %s", strings.Join(unregistered, ", "))
	}

	return nil
}

// ServeErrorFor gets a error handler for a given operation id
func (o *ImageapiAPI) ServeErrorFor(operationID string) func(http.ResponseWriter, *http.Request, error) {
	return o.ServeError
}

// AuthenticatorsFor gets the authenticators for the specified security schemes
func (o *ImageapiAPI) AuthenticatorsFor(schemes map[string]spec.SecurityScheme) map[string]runtime.Authenticator {
	return nil
}

// Authorizer returns the registered authorizer
func (o *ImageapiAPI) Authorizer() runtime.Authorizer {
	return nil
}

// ConsumersFor gets the consumers for the specified media types.
// MIME type parameters are ignored here.
func (o *ImageapiAPI) ConsumersFor(mediaTypes []string) map[string]runtime.Consumer {
	result := make(map[string]runtime.Consumer, len(mediaTypes))
	for _, mt := range mediaTypes {
		switch mt {
		case "application/json":
			result["application/json"] = o.JSONConsumer
		}

		if c, ok := o.customConsumers[mt]; ok {
			result[mt] = c
		}
	}
	return result
}

// ProducersFor gets the producers for the specified media types.
// MIME type parameters are ignored here.
func (o *ImageapiAPI) ProducersFor(mediaTypes []string) map[string]runtime.Producer {
	result := make(map[string]runtime.Producer, len(mediaTypes))
	for _, mt := range mediaTypes {
		switch mt {
		case "application/json":
			result["application/json"] = o.JSONProducer
		}

		if p, ok := o.customProducers[mt]; ok {
			result[mt] = p
		}
	}
	return result
}

// HandlerFor gets a http.Handler for the provided operation method and path
func (o *ImageapiAPI) HandlerFor(method, path string) (http.Handler, bool) {
	if o.handlers == nil {
		return nil, false
	}
	um := strings.ToUpper(method)
	if _, ok := o.handlers[um]; !ok {
		return nil, false
	}
	if path == "/" {
		path = ""
	}
	h, ok := o.handlers[um][path]
	return h, ok
}

// Context returns the middleware context for the imageapi API
func (o *ImageapiAPI) Context() *middleware.Context {
	if o.context == nil {
		o.context = middleware.NewRoutableContext(o.spec, o, nil)
	}

	return o.context
}

func (o *ImageapiAPI) initHandlerCache() {
	o.Context() // don't care about the result, just that the initialization happened
	if o.handlers == nil {
		o.handlers = make(map[string]map[string]http.Handler)
	}

	if o.handlers["DELETE"] == nil {
		o.handlers["DELETE"] = make(map[string]http.Handler)
	}
	o.handlers["DELETE"]["/mount"] = mounts.NewDeleteMount(o.context, o.MountsDeleteMountHandler)
	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/container"] = containers.NewCreateContainer(o.context, o.ContainersCreateContainerHandler)
	if o.handlers["DELETE"] == nil {
		o.handlers["DELETE"] = make(map[string]http.Handler)
	}
	o.handlers["DELETE"]["/container/{id}"] = containers.NewDeleteContainer(o.context, o.ContainersDeleteContainerHandler)
	if o.handlers["DELETE"] == nil {
		o.handlers["DELETE"] = make(map[string]http.Handler)
	}
	o.handlers["DELETE"]["/container/byname/{name}"] = containers.NewDeleteContainerByname(o.context, o.ContainersDeleteContainerBynameHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/container/{id}"] = containers.NewGetContainer(o.context, o.ContainersGetContainerHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/container/byname/{name}"] = containers.NewGetContainerByname(o.context, o.ContainersGetContainerBynameHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/mount/overlay/{id}"] = mounts.NewGetMountOverlay(o.context, o.MountsGetMountOverlayHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/mount/rbd/{id}"] = mounts.NewGetMountRbd(o.context, o.MountsGetMountRbdHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/attach/rbd/{id}"] = attach.NewGetRbd(o.context, o.AttachGetRbdHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/container"] = containers.NewListContainers(o.context, o.ContainersListContainersHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/mount"] = mounts.NewListMounts(o.context, o.MountsListMountsHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/mount/overlay"] = mounts.NewListMountsOverlay(o.context, o.MountsListMountsOverlayHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/mount/rbd"] = mounts.NewListMountsRbd(o.context, o.MountsListMountsRbdHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/attach/rbd"] = attach.NewListRbds(o.context, o.AttachListRbdsHandler)
	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/attach/rbd"] = attach.NewMapRbd(o.context, o.AttachMapRbdHandler)
	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/mount"] = mounts.NewMount(o.context, o.MountsMountHandler)
	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/mount/overlay"] = mounts.NewMountOverlay(o.context, o.MountsMountOverlayHandler)
	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/mount/rbd"] = mounts.NewMountRbd(o.context, o.MountsMountRbdHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/container/{id}/{state}"] = containers.NewSetContainerState(o.context, o.ContainersSetContainerStateHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/container/byname/{name}/{state}"] = containers.NewSetContainerStateByname(o.context, o.ContainersSetContainerStateBynameHandler)
	if o.handlers["DELETE"] == nil {
		o.handlers["DELETE"] = make(map[string]http.Handler)
	}
	o.handlers["DELETE"]["/attach/rbd/{id}"] = attach.NewUnmapRbd(o.context, o.AttachUnmapRbdHandler)
	if o.handlers["DELETE"] == nil {
		o.handlers["DELETE"] = make(map[string]http.Handler)
	}
	o.handlers["DELETE"]["/mount/overlay/{id}"] = mounts.NewUnmountOverlay(o.context, o.MountsUnmountOverlayHandler)
	if o.handlers["DELETE"] == nil {
		o.handlers["DELETE"] = make(map[string]http.Handler)
	}
	o.handlers["DELETE"]["/mount/rbd/{id}"] = mounts.NewUnmountRbd(o.context, o.MountsUnmountRbdHandler)
}

// Serve creates a http handler to serve the API over HTTP
// can be used directly in http.ListenAndServe(":8000", api.Serve(nil))
func (o *ImageapiAPI) Serve(builder middleware.Builder) http.Handler {
	o.Init()

	if o.Middleware != nil {
		return o.Middleware(builder)
	}
	if o.useSwaggerUI {
		return o.context.APIHandlerSwaggerUI(builder)
	}
	return o.context.APIHandler(builder)
}

// Init allows you to just initialize the handler cache, you can then recompose the middleware as you see fit
func (o *ImageapiAPI) Init() {
	if len(o.handlers) == 0 {
		o.initHandlerCache()
	}
}

// RegisterConsumer allows you to add (or override) a consumer for a media type.
func (o *ImageapiAPI) RegisterConsumer(mediaType string, consumer runtime.Consumer) {
	o.customConsumers[mediaType] = consumer
}

// RegisterProducer allows you to add (or override) a producer for a media type.
func (o *ImageapiAPI) RegisterProducer(mediaType string, producer runtime.Producer) {
	o.customProducers[mediaType] = producer
}

// AddMiddlewareFor adds a http middleware to existing handler
func (o *ImageapiAPI) AddMiddlewareFor(method, path string, builder middleware.Builder) {
	um := strings.ToUpper(method)
	if path == "/" {
		path = ""
	}
	o.Init()
	if h, ok := o.handlers[um][path]; ok {
		o.handlers[method][path] = builder(h)
	}
}
