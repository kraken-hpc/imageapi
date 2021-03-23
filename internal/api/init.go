package api

import (
	"errors"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var Rbds *RbdsType
var MountsRbd *MountsRBDType
var MountsOverlay *MountsOverlayType
var Containers *ContainersType
var Log *logrus.Logger

const mountDir = "/var/run/imageapi/mounts"
const logDir = "/var/run/imageapi/logs"
const collectTime = time.Second * 1

var ERRNOTFOUND = errors.New("not found")

func garbageCollect() {
	for {
		time.Sleep(collectTime)
		MountsOverlay.Collect()
		MountsRbd.Collect()
		Rbds.Collect()
	}
}

var llStringToLL = map[string]logrus.Level{
	"PANIC": logrus.PanicLevel,
	"FATAL": logrus.FatalLevel,
	"ERROR": logrus.ErrorLevel,
	"WARN":  logrus.WarnLevel,
	"INFO":  logrus.InfoLevel,
	"DEBUG": logrus.DebugLevel,
	"TRACE": logrus.TraceLevel,
}

func init() {
	Log = logrus.New()
	// fixme: read from os.Env?
	logLevel := logrus.InfoLevel
	if val, ok := os.LookupEnv("IMAGEAPI_LOGLEVEL"); ok {
		if ll, ok := llStringToLL[val]; ok {
			logLevel = ll
		}
	}
	Log.Level = logLevel
	Log.Info("initializing imageapi-server")
	Rbds = &RbdsType{}
	Rbds.Init()
	MountsRbd = &MountsRBDType{}
	MountsRbd.Init()
	MountsOverlay = &MountsOverlayType{}
	MountsOverlay.Init()
	Containers = &ContainersType{}
	Containers.Init()
	Log.WithField("collectTime", collectTime).Debug("starting garbage collection")
	go garbageCollect()
}
