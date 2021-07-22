package api

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/kraken-hpc/imageapi/models"
	"github.com/kraken-hpc/imageapi/restapi/operations/attachments"
	"github.com/kraken-hpc/imageapi/restapi/operations/containers"
	"github.com/kraken-hpc/imageapi/restapi/operations/mounts"
)

func errorSanitize(err error) error {
	if _, ok := errorToHTTP[err]; !ok {
		return ERRSRV
	}
	return err
}

func errorPayload(err error) *models.Error {
	return &models.Error{Code: int64(errorToHTTP[err]), Message: swag.String(err.Error())}
}

// Attachments
var AttachmentsAttachHandler = attachments.AttachHandlerFunc(func(params attachments.AttachParams) middleware.Responder {
	var err error
	var a *Attach
	if a, err = API.Attachments.Attach((*Attach)(params.Attach)); err != nil {
		err = errorSanitize(err)
		return attachments.NewAttachDefault(errorToHTTP[err]).WithPayload(errorPayload(err))
	}
	return attachments.NewAttachCreated().WithPayload((*models.Attach)(a))
})

var AttachmentsDeleteAttachHandler = attachments.DeleteAttachHandlerFunc(func(params attachments.DeleteAttachParams) middleware.Responder {
	var err error
	var a *Attach
	if a, err = API.Attachments.Detach((*Attach)(params.Attach)); err != nil {
		err = errorSanitize(err)
		return attachments.NewDeleteAttachDefault(errorToHTTP[err]).WithPayload(errorPayload(err))
	}
	return attachments.NewDeleteAttachOK().WithPayload((*models.Attach)(a))
})

var AttachmentsListAttachmentsHandler = attachments.ListAttachmentsHandlerFunc(func(params attachments.ListAttachmentsParams) middleware.Responder {
	as := []*models.Attach{}
	if *params.ID != 0 {
		// Get
		a := API.Attachments.Get((models.ID)(*params.ID))
		if a == nil {
			return attachments.NewDeleteAttachDefault(errorToHTTP[ERRNOTFOUND]).WithPayload(errorPayload(ERRNOTFOUND))
		}
		as = append(as, (*models.Attach)(a))
	} else {
		// List
		ia := API.Attachments.List()
		for _, a := range ia {
			if *params.Kind != "" && a.Kind != *params.Kind {
				continue
			}
			as = append(as, (*models.Attach)(a))
		}
	}
	return attachments.NewListAttachmentsOK().WithPayload(as)
})

// Mounts
var MountsMountHandler = mounts.MountHandlerFunc(func(params mounts.MountParams) middleware.Responder {
	var err error
	var m *Mount
	if m, err = API.Mounts.Mount((*Mount)(params.Mount)); err != nil {
		err = errorSanitize(err)
		return mounts.NewMountDefault(errorToHTTP[err]).WithPayload(errorPayload(err))
	}
	return mounts.NewMountCreated().WithPayload((*models.Mount)(m))
})

var MountsDeleteMountHandler = mounts.DeleteMountHandlerFunc(func(params mounts.DeleteMountParams) middleware.Responder {
	var err error
	var m *Mount
	if m, err = API.Mounts.Unmount((*Mount)(params.Mount)); err != nil {
		err = errorSanitize(err)
		return mounts.NewDeleteMountDefault(errorToHTTP[err]).WithPayload(errorPayload(err))
	}
	return mounts.NewDeleteMountOK().WithPayload((*models.Mount)(m))
})

var MountsListMountsHandler = mounts.ListMountsHandlerFunc(func(params mounts.ListMountsParams) middleware.Responder {
	ms := []*models.Mount{}
	if *params.ID != 0 {
		// Get
		m := API.Mounts.Get((models.ID)(*params.ID))
		if m == nil {
			return mounts.NewListMountsDefault(errorToHTTP[ERRNOTFOUND]).WithPayload(errorPayload(ERRNOTFOUND))
		}
		ms = append(ms, (*models.Mount)(m))
	} else {
		// List
		im := API.Mounts.List()
		for _, m := range im {
			if *params.Kind != "" && m.Kind != *params.Kind {
				continue
			}
			ms = append(ms, (*models.Mount)(m))
		}
	}
	return mounts.NewListMountsOK().WithPayload(ms)
})

// Containers
var ContainersCreateContainerHandler = containers.CreateContainerHandlerFunc(func(params containers.CreateContainerParams) middleware.Responder {
	var err error
	var c *Container
	if c, err = API.Containers.Create(&Container{Container: params.Container}); err != nil {
		err = errorSanitize(err)
		return containers.NewCreateContainerDefault(errorToHTTP[err]).WithPayload(errorPayload(err))
	}
	return containers.NewCreateContainerCreated().WithPayload(c.Container)
})

var ContainersDeleteContainerHandler = containers.DeleteContainerHandlerFunc(func(params containers.DeleteContainerParams) middleware.Responder {
	var err error
	var c *Container
	if *params.ID == 0 {
		if *params.Name == "" {
			return containers.NewDeleteContainerDefault(errorToHTTP[ERRINVALDAT]).WithPayload(&models.Error{Code: int64(errorToHTTP[ERRINVALDAT]), Message: swag.String("either ID or Name must be provided")})
		}
		*params.ID = int64(API.Containers.NameGetID((models.Name)(*params.Name)))
	}
	if *params.ID < 1 {
		return containers.NewDeleteContainerDefault(errorToHTTP[ERRNOTFOUND]).WithPayload(errorPayload(ERRNOTFOUND))
	}
	if c, err = API.Containers.Delete((models.ID)(*params.ID)); err != nil {
		err = errorSanitize(err)
		return containers.NewDeleteContainerDefault(errorToHTTP[err]).WithPayload(errorPayload(err))
	}
	return containers.NewDeleteContainerOK().WithPayload(c.Container)
})

var ContainersListContainersHandler = containers.ListContainersHandlerFunc(func(params containers.ListContainersParams) middleware.Responder {
	cs := []*models.Container{}
	if *params.ID != 0 || *params.Name != "" {
		// Get
		if *params.ID == 0 {
			*params.ID = int64(API.Containers.NameGetID((models.Name)(*params.Name)))
		}
		if *params.ID < 1 {
			return containers.NewListContainersDefault(errorToHTTP[ERRNOTFOUND]).WithPayload(errorPayload(ERRNOTFOUND))
		}
		c := API.Containers.Get((models.ID)((*params.ID)))
		if c == nil {
			return containers.NewListContainersDefault(errorToHTTP[ERRNOTFOUND]).WithPayload(errorPayload(ERRNOTFOUND))
		}
		cs = append(cs, c.Container)
	} else {
		// List
		is := API.Containers.List()
		for _, c := range is {
			if *params.State != "" && c.Container.State != models.ContainerState(*params.State) {
				continue
			}
			cs = append(cs, c.Container)
		}
	}
	return containers.NewListContainersOK().WithPayload(cs)
})

var ContainersSetContainerStateHandler = containers.SetContainerStateHandlerFunc(func(params containers.SetContainerStateParams) middleware.Responder {
	if *params.ID == 0 {
		if *params.Name == "" {
			return containers.NewSetContainerStateDefault(errorToHTTP[ERRINVALDAT]).WithPayload(&models.Error{Code: int64(errorToHTTP[ERRINVALDAT]), Message: swag.String("either ID or Name must be provided")})
		}
		*params.ID = int64(API.Containers.NameGetID((models.Name)(*params.Name)))
	}
	if *params.ID < 1 {
		return containers.NewSetContainerStateDefault(errorToHTTP[ERRNOTFOUND]).WithPayload(errorPayload(ERRNOTFOUND))
	}
	var c *Container
	var err error
	if c, err = API.Containers.SetState((models.ID)(*params.ID), models.ContainerState(params.State)); err != nil {
		err = errorSanitize(err)
		return containers.NewCreateContainerDefault(errorToHTTP[err]).WithPayload(errorPayload(err))
	}
	return containers.NewSetContainerStateOK().WithPayload(c.Container)
})
