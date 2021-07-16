package types

import (
	"errors"

	"github.com/kraken-hpc/imageapi/models"
	"github.com/sirupsen/logrus"
)

// Generic
type Collectable interface {
	Collect()
	RefAdd(models.ID, int64)
}

type Listable interface {
	List() []Model
	Get(models.ID) (Model, error)
}

type Endpoint interface {
	Init()
	Listable
	Collectable
}

type AttachTypeEnum uint8

// Any Model
type Model interface {
	GetID() models.ID
	GetRefs() int64
}

// Models for attachment objects
type AttachModel interface {
	Model
	GetDevice() string
}

type Attach interface {
	Attach(AttachModel) (AttachModel, error)
	Detach(models.ID) (AttachModel, error)
	Endpoint
}

type MountModel interface {
	GetMountpoint() string
	Model
}

type Mount interface {
	Mount(MountModel)
	Unmount(models.ID)
	Endpoint
}

type Container interface {
	Endpoint
}

var ERRNOTFOUND = errors.New("not found")
var ERRINVALDAT = errors.New("invalid data type")

var LogStringToLL = map[string]logrus.Level{
	"PANIC": logrus.PanicLevel,
	"FATAL": logrus.FatalLevel,
	"ERROR": logrus.ErrorLevel,
	"WARN":  logrus.WarnLevel,
	"INFO":  logrus.InfoLevel,
	"DEBUG": logrus.DebugLevel,
	"TRACE": logrus.TraceLevel,
}
