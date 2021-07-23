// Code generated by go-swagger; DO NOT EDIT.

package containers

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// New creates a new containers API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) ClientService {
	return &Client{transport: transport, formats: formats}
}

/*
Client for containers API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

// ClientOption is the option for Client methods
type ClientOption func(*runtime.ClientOperation)

// ClientService is the interface for Client methods
type ClientService interface {
	CreateContainer(params *CreateContainerParams, opts ...ClientOption) (*CreateContainerCreated, error)

	DeleteContainer(params *DeleteContainerParams, opts ...ClientOption) (*DeleteContainerOK, error)

	ListContainers(params *ListContainersParams, opts ...ClientOption) (*ListContainersOK, error)

	SetContainerState(params *SetContainerStateParams, opts ...ClientOption) (*SetContainerStateOK, error)

	SetTransport(transport runtime.ClientTransport)
}

/*
  CreateContainer Create a container
*/
func (a *Client) CreateContainer(params *CreateContainerParams, opts ...ClientOption) (*CreateContainerCreated, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewCreateContainerParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "create_container",
		Method:             "POST",
		PathPattern:        "/container",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &CreateContainerReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*CreateContainerCreated)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*CreateContainerDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
  DeleteContainer Delete a container defition.
Either `id` or `name` query parameter must be specified.

*/
func (a *Client) DeleteContainer(params *DeleteContainerParams, opts ...ClientOption) (*DeleteContainerOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewDeleteContainerParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "delete_container",
		Method:             "DELETE",
		PathPattern:        "/container",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &DeleteContainerReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*DeleteContainerOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*DeleteContainerDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
  ListContainers Get a list of containers
*/
func (a *Client) ListContainers(params *ListContainersParams, opts ...ClientOption) (*ListContainersOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewListContainersParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "list_containers",
		Method:             "GET",
		PathPattern:        "/container",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &ListContainersReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*ListContainersOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*ListContainersDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
  SetContainerState Request a (valid) state for a container.
Valid states to request include: `running`, `exited`, `paused` (paused is not yet implemented)

Either a valid Name or ID must be passed as a query parameter, along with a valid state parameter.

*/
func (a *Client) SetContainerState(params *SetContainerStateParams, opts ...ClientOption) (*SetContainerStateOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewSetContainerStateParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "set_container_state",
		Method:             "PATCH",
		PathPattern:        "/container",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &SetContainerStateReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*SetContainerStateOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*SetContainerStateDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
