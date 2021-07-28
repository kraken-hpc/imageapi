package api

import (
	"path"

	"github.com/bensallen/rbd/pkg/mount"
	"github.com/kraken-hpc/imageapi/models"
	"github.com/sirupsen/logrus"
)

func init() {
	MountDrivers[models.MountKindBind] = &MountDriverBind{}
}

type MountDriverBind struct {
	log *logrus.Entry //
}

func (m *MountDriverBind) Init(log *logrus.Entry) {
	m.log = log
	m.log.Trace("initialized")
}

func (m *MountDriverBind) Mount(mnt *Mount) (ret *Mount, err error) {
	l := m.log.WithField("operation", "mount")
	if mnt.Bind == nil {
		l.Trace("attempted bind mount with no mount definition")
		return nil, ErrInvalDat
	}
	// go-swagger should handle other validation that we need
	base := "/"
	if *mnt.Bind.Base == models.MountBindBaseMount {
		if mnt.Bind.Mount == nil {
			l.Debug("bind mount called with mount base but no mount definition")
			return nil, ErrInvalDat
		}
		m, err := API.Mounts.GetOrMount((*Mount)(mnt.Bind.Mount))
		if err != nil {
			l.Error("base mount failed to GetOrMount")
			return nil, ErrFail
		}
		mnt.Bind.Mount = (*models.Mount)(m)
		defer func() {
			if err != nil {
				API.Store.RefAdd(mnt.Bind.Mount.ID, -1)
			}
		}()
		base = mnt.Bind.Mount.Mountpoint
	}
	fullPath := path.Join(base, *mnt.Bind.Path)
	l.WithField("fullPath", fullPath)
	// ok, we're good to attempt the mount
	// make a mountpoint
	options := []string{"bind"}
	if mnt.Bind.Recursive != nil && *mnt.Bind.Recursive {
		options = append(options, "rec")
	}
	if mnt.Bind.Ro != nil && *mnt.Bind.Ro {
		options = append(options, "ro")
	}
	if err = mount.Mount(fullPath, mnt.Mountpoint, "bind", options); err != nil {
		l.WithError(err).Error("failed to mount")
		return nil, ErrFail
	}
	return mnt, nil
}

func (m *MountDriverBind) Unmount(mnt *Mount) (ret *Mount, err error) {
	l := m.log.WithFields(logrus.Fields{
		"operation": "unmount",
		"id":        mnt.ID,
	})

	// always lazy unmount.  Good idea?
	if err = mount.Unmount(mnt.Mountpoint, false, true); err != nil {
		l.WithError(err).Error("unmount failed")
		return nil, ErrFail
	}
	if *mnt.Bind.Base == models.MountBindBaseMount {
		API.Store.RefAdd(mnt.Bind.Mount.ID, -1)
	}
	// garbage collection should do our cleanup
	return mnt, nil
}
