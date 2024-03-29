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
		return ErrSrv
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
	if a, err = API.Attachments.Detach(&Attach{ID: (models.ID)(params.ID)}, *params.Force); err != nil {
		err = errorSanitize(err)
		return attachments.NewDeleteAttachDefault(errorToHTTP[err]).WithPayload(errorPayload(err))
	}
	return attachments.NewDeleteAttachOK().WithPayload((*models.Attach)(a))
})

var AttachmentsListAttachmentsHandler = attachments.ListAttachmentsHandlerFunc(func(params attachments.ListAttachmentsParams) middleware.Responder {
	as := []*models.Attach{}
	if params.ID != nil {
		// Get
		a := API.Attachments.Get((models.ID)(*params.ID))
		if a == nil {
			return attachments.NewDeleteAttachDefault(errorToHTTP[ErrNotFound]).WithPayload(errorPayload(ErrNotFound))
		}
		as = append(as, (*models.Attach)(a))
	} else {
		// List
		ia := API.Attachments.List()
		for _, a := range ia {
			if params.Kind != nil && a.Kind != *params.Kind {
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
	if m, err = API.Mounts.Unmount(&Mount{ID: (models.ID)(params.ID)}, *params.Force); err != nil {
		err = errorSanitize(err)
		return mounts.NewDeleteMountDefault(errorToHTTP[err]).WithPayload(errorPayload(err))
	}
	return mounts.NewDeleteMountOK().WithPayload((*models.Mount)(m))
})

var MountsListMountsHandler = mounts.ListMountsHandlerFunc(func(params mounts.ListMountsParams) middleware.Responder {
	ms := []*models.Mount{}
	if params.ID != nil {
		// Get
		m := API.Mounts.Get((models.ID)(*params.ID))
		if m == nil {
			return mounts.NewListMountsDefault(errorToHTTP[ErrNotFound]).WithPayload(errorPayload(ErrNotFound))
		}
		ms = append(ms, (*models.Mount)(m))
	} else {
		// List
		im := API.Mounts.List()
		for _, m := range im {
			if params.Kind != nil && m.Kind != *params.Kind {
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
	var id models.ID
	if *params.Force {
		return containers.NewDeleteContainerDefault(501).WithPayload(&models.Error{Code: 501, Message: swag.String("forced container deletion is not yet implemented")})
	}
	if params.ID == nil {
		if params.Name == nil {
			return containers.NewDeleteContainerDefault(errorToHTTP[ErrInvalDat]).WithPayload(&models.Error{Code: int64(errorToHTTP[ErrInvalDat]), Message: swag.String("either ID or Name must be provided")})
		}
		id = API.Containers.NameGetID((models.Name)(*params.Name))
	} else {
		id = models.ID(*params.ID)
	}
	if id < 1 {
		return containers.NewDeleteContainerDefault(errorToHTTP[ErrNotFound]).WithPayload(errorPayload(ErrNotFound))
	}
	if c, err = API.Containers.Delete(id); err != nil {
		err = errorSanitize(err)
		return containers.NewDeleteContainerDefault(errorToHTTP[err]).WithPayload(errorPayload(err))
	}
	return containers.NewDeleteContainerOK().WithPayload(c.Container)
})

var ContainersListContainersHandler = containers.ListContainersHandlerFunc(func(params containers.ListContainersParams) middleware.Responder {
	cs := []*models.Container{}
	if params.ID != nil || params.Name != nil {
		// Get
		var id models.ID
		if params.ID == nil {
			id = API.Containers.NameGetID((models.Name)(*params.Name))
		} else {
			id = models.ID(*params.ID)
		}
		if id < 1 {
			return containers.NewListContainersDefault(errorToHTTP[ErrNotFound]).WithPayload(errorPayload(ErrNotFound))
		}
		c := API.Containers.Get(id)
		if c == nil {
			return containers.NewListContainersDefault(errorToHTTP[ErrNotFound]).WithPayload(errorPayload(ErrNotFound))
		}
		cs = append(cs, c.Container)
	} else {
		// List
		is := API.Containers.List()
		for _, c := range is {
			if params.State != nil && c.Container.State != models.ContainerState(*params.State) {
				continue
			}
			cs = append(cs, c.Container)
		}
	}
	return containers.NewListContainersOK().WithPayload(cs)
})

var ContainersSetContainerStateHandler = containers.SetContainerStateHandlerFunc(func(params containers.SetContainerStateParams) middleware.Responder {
	var id models.ID
	if params.ID == nil {
		if params.Name == nil {
			return containers.NewSetContainerStateDefault(errorToHTTP[ErrInvalDat]).WithPayload(&models.Error{Code: int64(errorToHTTP[ErrInvalDat]), Message: swag.String("either ID or Name must be provided")})
		}
		id = API.Containers.NameGetID((models.Name)(*params.Name))
	} else {
		id = models.ID(*params.ID)
	}
	if id < 1 {
		return containers.NewSetContainerStateDefault(errorToHTTP[ErrNotFound]).WithPayload(errorPayload(ErrNotFound))
	}
	var c *Container
	var err error
	if c, err = API.Containers.SetState(id, models.ContainerState(params.State)); err != nil {
		err = errorSanitize(err)
		return containers.NewCreateContainerDefault(errorToHTTP[err]).WithPayload(errorPayload(err))
	}
	return containers.NewSetContainerStateOK().WithPayload(c.Container)
})
