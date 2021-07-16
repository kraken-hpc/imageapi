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
	List() []models.EndpointObject
	Get(models.ID) (models.EndpointObject, error)
}

type Endpoint interface {
	Init()
	Listable
	Collectable
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
