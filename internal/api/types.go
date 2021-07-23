package api

import (
	"errors"
	"time"

	"github.com/kraken-hpc/imageapi/models"
	"github.com/sirupsen/logrus"
)

type EndpointObjectType uint8

// Allows inspection without reflection
const (
	EndpointObjectAttach    EndpointObjectType = iota
	EndpointObjectMount     EndpointObjectType = iota
	EndpointObjectContainer EndpointObjectType = iota
)

type EndpointObject interface {
	SetID(models.ID)
	GetID() models.ID
	GetRefs() int64
	RefAdd(int64)
	EndpointObjectType() EndpointObjectType
}

type AttachDriver interface {
	Init(*logrus.Entry)
	Attach(*Attach) (*Attach, error)
	Detach(*Attach) (*Attach, error)
}

type MountDriver interface {
	Init(*logrus.Entry)
	Mount(*Mount) (*Mount, error)
	Unmount(*Mount) (*Mount, error)
}

type APIType struct {
	Store            *ObjectStore
	Mounts           *Mounts
	Attachments      *Attachments
	Containers       *Containers
	Log              *logrus.Entry
	MountDir, LogDir string
	CollectInterval  time.Duration
}

var ERRNOTFOUND = errors.New("not found")
var ERRINVALDAT = errors.New("invalid data type")
var ERRBUSY = errors.New("object is busy")
var ERRNODRV = errors.New("no driver found for this object type")
var ERRSRV = errors.New("internal server error")
var ERRFAIL = errors.New("operation failed")

var errorToHTTP = map[error]int{
	ERRNOTFOUND: 404,
	ERRINVALDAT: 400,
	ERRBUSY:     409,
	ERRNODRV:    501,
	ERRSRV:      500,
	ERRFAIL:     500,
}

var LogStringToLL = map[string]logrus.Level{
	"PANIC": logrus.PanicLevel,
	"FATAL": logrus.FatalLevel,
	"ERROR": logrus.ErrorLevel,
	"WARN":  logrus.WarnLevel,
	"INFO":  logrus.InfoLevel,
	"DEBUG": logrus.DebugLevel,
	"TRACE": logrus.TraceLevel,
}
