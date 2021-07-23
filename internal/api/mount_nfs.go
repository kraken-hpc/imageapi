package api

import (
	"fmt"

	"github.com/bensallen/rbd/pkg/mount"
	"github.com/kraken-hpc/imageapi/models"
	"github.com/sirupsen/logrus"
)

func init() {
	MountDrivers[models.MountKindNfs] = &MountDriverNFS{}
}

type MountDriverNFS struct {
	log *logrus.Entry //
}

func (m *MountDriverNFS) Init(log *logrus.Entry) {
	m.log = log
	m.log.Trace("initialized")
}

func (m *MountDriverNFS) Mount(mnt *Mount) (ret *Mount, err error) {
	l := m.log.WithField("operation", "mount")
	if mnt.Nfs == nil {
		l.Trace("attempted nfs mount with no mount definition")
		return nil, ERRINVALDAT
	}
	l = l.WithFields(logrus.Fields{
		"host": *mnt.Nfs.Host,
		"path": *mnt.Nfs.Path,
	})
	// ok, we're good to attempt the mount
	// make a mountpoint
	if err = mount.Mount(fmt.Sprintf("%s:%s", *mnt.Nfs.Host, *mnt.Nfs.Path), mnt.Mountpoint, "nfs", mnt.Nfs.MountOptions); err != nil {
		l.WithError(err).Error("failed to mount")
		return nil, ERRFAIL
	}
	l.Info("successfully mounted")
	return mnt, nil
}

func (m *MountDriverNFS) Unmount(mnt *Mount) (ret *Mount, err error) {
	l := m.log.WithFields(logrus.Fields{
		"operation": "unmount",
		"id":        mnt.ID,
		"attach_id": mnt.Attach.Attach.ID,
	})

	// always lazy unmount.  Good idea?
	if err = mount.Unmount(mnt.Mountpoint, false, true); err != nil {
		l.WithError(err).Error("unmount failed")
		return nil, ERRFAIL
	}
	API.Store.RefAdd(mnt.Attach.Attach.ID, -1)
	l.Info("successfully unmounted")
	// garbage collection should do our cleanup
	return mnt, nil
}
