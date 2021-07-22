// initialization for internal data structures & processes
package api

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var API = &APIType{
	Store:           &ObjectStore{},
	Mounts:          &Mounts{},
	Attachments:     &Attachments{},
	Containers:      &Containers{},
	MountDir:        "/var/run/imageapi/mounts",
	LogDir:          "/var/run/imageapi/logs",
	CollectInterval: time.Second,
}

func Init() {
	loglevel := logrus.InfoLevel
	if llstr, ok := os.LookupEnv("IMAGEAPI_LOGLEVEL"); ok {
		if ll, ok := LogStringToLL[llstr]; ok {
			loglevel = ll
		}
	}
	l := logrus.New()
	l.SetLevel(loglevel)
	API.Log = l.WithField("application", "imageapi")
	API.Log.Infof("using log level: %s", loglevel)
	API.Log.Trace("initializing object store")
	API.Store.Init()
	API.Log.Trace("initializing mounts subsystem")
	API.Mounts.Init(API.Log.WithField("subsys", "mount"))
	API.Attachments.Init(API.Log.WithField("subsys", "attach"))
	API.Containers.Init(API.Log.WithField("subsys", "container"))
}
