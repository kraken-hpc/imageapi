package api

// API operations on rbd maps

import (
	"os"
	"path"

	"github.com/kraken-hpc/imageapi/models"
	"github.com/sirupsen/logrus"
	"github.com/u-root/u-root/pkg/mount/loop"
)

func init() {
	AttachDrivers[models.AttachKindLoopback] = &AttachDriverLoopback{}
}

type AttachDriverLoopback struct {
	log *logrus.Entry
}

func (a *AttachDriverLoopback) Init(log *logrus.Entry) {
	a.log = log
	a.log.Trace("initialized")
}

func (a *AttachDriverLoopback) Attach(att *Attach) (ret *Attach, err error) {
	// sanity check
	l := a.log.WithField("operation", "attach")
	if att.Loopback == nil {
		l.Trace("attempted to attach loopback with no loopback definition")
		return nil, ERRINVALDAT
	}
	l = l.WithFields(logrus.Fields{
		"path": att.Loopback.Path,
		"base": att.Loopback.Base,
	})

	base := "/"
	if *att.Loopback.Base == models.MountBindBaseMount {
		if att.Loopback.Mount == nil {
			l.Debug("bind mount called with mount base but no mount definition")
			return nil, ERRINVALDAT
		}
		m, err := API.Mounts.GetOrMount((*Mount)(att.Loopback.Mount))
		if err != nil {
			l.Error("base mount failed to GetOrMount")
			return nil, ERRFAIL
		}
		att.Loopback.Mount = (*models.Mount)(m)
		defer func() {
			if err != nil {
				API.Store.RefAdd(att.Loopback.Mount.ID, -1)
			}
		}()
		base = att.Loopback.Mount.Mountpoint
	}
	fullPath := path.Join(base, *att.Loopback.Path)
	l.WithField("fullPath", fullPath)
	info, err := os.Stat(fullPath)
	if err != nil {
		l.Debug("loopback attach called on file that doesn't exist")
		return nil, ERRINVALDAT
	}
	if !info.Mode().IsRegular() {
		l.Debug("loopback attach call on a non-regular file")
		return nil, ERRINVALDAT
	}

	att.DeviceFile, err = loop.FindDevice()
	if err != nil {
		l.WithError(err).Error("failed to acquire loopback device")
		return nil, ERRFAIL
	}
	if err = loop.SetFile(att.DeviceFile, fullPath); err != nil {
		l.WithError(err).Debug("failed to assign file to loopback device")
		return nil, ERRFAIL
	}

	l.Info("successfully mapped")
	return att, err
}

func (a *AttachDriverLoopback) Detach(att *Attach) (ret *Attach, err error) {
	l := a.log.WithFields(logrus.Fields{
		"operation": "unmap",
		"id":        att.ID,
		"path":      att.Loopback.Path,
		"base":      att.Loopback.Base,
	})
	if err = loop.ClearFile(att.DeviceFile); err != nil {
		l.Debug("failed to clear loopback association")
		return nil, ERRFAIL
	}
	return att, nil
}
