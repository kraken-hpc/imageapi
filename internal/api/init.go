package api

import (
	"errors"
	"time"

	"github.com/sirupsen/logrus"
)

var Rbds RbdsType
var MountsRbd MountsRBDType
var MountsOverlay MountsOverlayType
var Containers ContainersType
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

func init() {
	Log = logrus.New()
	// fixme: read from os.Env?
	Log.Level = logrus.TraceLevel
	Log.Info("initializing imageapi-server")
	Rbds = RbdsType{}
	Rbds.Init()
	MountsRbd = MountsRBDType{}
	MountsRbd.Init()
	MountsOverlay = MountsOverlayType{}
	MountsOverlay.Init()
	Containers = ContainersType{}
	Containers.Init()
	Log.WithField("collectTime", collectTime).Debug("starting garbage collection")
	go garbageCollect()
}
