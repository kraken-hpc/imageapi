package api

import (
	"os"

	"github.com/kraken-hpc/imageapi/models"
	"github.com/sirupsen/logrus"
)

func init() {
	AttachDrivers[models.AttachKindLoopback] = &AttachDriverLocal{}
}

type AttachDriverLocal struct {
	log *logrus.Entry
}

func (a *AttachDriverLocal) Init(log *logrus.Entry) {
	a.log = log
	a.log.Trace("initialized")
}

func (a *AttachDriverLocal) Attach(att *Attach) (ret *Attach, err error) {
	// sanity check
	l := a.log.WithField("operation", "attach")
	if att.Local == nil {
		l.Trace("attempted to attach local with no local definition")
		return nil, ERRINVALDAT
	}
	l = l.WithField("path", att.Local.Path)

	finfo, err := os.Stat(*att.Local.Path)
	if err != nil {
		l.WithError(err).Debug("failed to stat device file")
		return nil, ERRFAIL
	}

	if finfo.Mode()&os.ModeDevice == 0 {
		l.Trace("path is not a device file")
		return nil, ERRINVALDAT
	}

	if finfo.Mode()&os.ModeCharDevice != 0 {
		l.Trace("path points to character device")
		return nil, ERRINVALDAT
	}
	att.DeviceFile = *att.Local.Path

	l.Info("successfully mapped")
	return att, nil
}

func (a *AttachDriverLocal) Detach(att *Attach) (ret *Attach, err error) {
	// this is a dummy operation
	return att, nil
}
