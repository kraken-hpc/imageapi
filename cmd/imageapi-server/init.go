// initialize internal objects
// we keep this separated from main.go to minimize the amount of modification to generated code.

package main

import (
	"os"
	"time"

	internal "github.com/kraken-hpc/imageapi/internal/api"
	"github.com/sirupsen/logrus"
)

const collectTime = time.Second * 1

func initInternal() {
	internal.Log = logrus.New()
	// fixme: read from os.Env?
	logLevel := logrus.InfoLevel
	if val, ok := os.LookupEnv("IMAGEAPI_LOGLEVEL"); ok {
		if ll, ok := internal.LogStringToLL[val]; ok {
			logLevel = ll
		}
	}
	internal.MountDir = "/var/run/imageapi/mounts"
	internal.LogDir = "/var/run/imageapi/logs"
	if val, ok := os.LookupEnv("IMAGEAPI_MOUNTDIR"); ok {
		internal.MountDir = val
	}
	if val, ok := os.LookupEnv("IMAGEAPI_LOGDIR"); ok {
		internal.LogDir = val
	}

	internal.ForkInit()

	internal.Log.Level = logLevel
	internal.Log.Info("initializing imageapi-server")
	internal.Rbds = &internal.RbdsType{}
	internal.Rbds.Init()
	internal.MountsRbd = &internal.MountsRBDType{}
	internal.MountsRbd.Init()
	internal.MountsOverlay = &internal.MountsOverlayType{}
	internal.MountsOverlay.Init()
	internal.Containers = &internal.ContainersType{}
	internal.Containers.Init()
	internal.Log.WithField("collectTime", collectTime).Debug("starting garbage collection")

	go func() {
		for {
			time.Sleep(collectTime)
			internal.GarbageCollect()
		}
	}()
}
