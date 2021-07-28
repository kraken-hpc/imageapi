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

var ErrNotFound = errors.New("not found")
var ErrInvalDat = errors.New("invalid data type")
var ErrBusy = errors.New("object is busy")
var ErrNoDrv = errors.New("no driver found for this object type")
var ErrSrv = errors.New("internal server error")
var ErrFail = errors.New("operation failed")

var errorToHTTP = map[error]int{
	ErrNotFound: 404,
	ErrInvalDat: 400,
	ErrBusy:     409,
	ErrNoDrv:    501,
	ErrSrv:      500,
	ErrFail:     500,
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
