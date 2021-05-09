package api

import (
	"errors"

	"github.com/sirupsen/logrus"
)

var Rbds *RbdsType
var MountsRbd *MountsRBDType
var MountsOverlay *MountsOverlayType
var Containers *ContainersType
var Log *logrus.Logger

var MountDir string = "/var/run/imageapi/mounts"
var LogDir string = "/var/run/imageapi/logs"

var ERRNOTFOUND = errors.New("not found")

func GarbageCollect() {
	MountsOverlay.Collect()
	MountsRbd.Collect()
	Rbds.Collect()
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
