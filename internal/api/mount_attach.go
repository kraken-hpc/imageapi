package api

import (
	"github.com/bensallen/rbd/pkg/mount"
	"github.com/kraken-hpc/imageapi/models"
	"github.com/sirupsen/logrus"
)

func init() {
	MountDrivers[models.MountKindAttach] = &MountDriverAttach{}
}

type MountDriverAttach struct {
	log *logrus.Entry //
}

func (m *MountDriverAttach) Init(log *logrus.Entry) {
	m.log = log
	m.log.Trace("initialized")
}

func (m *MountDriverAttach) Mount(mnt *Mount) (ret *Mount, err error) {
	l := m.log.WithField("operation", "mount")
	if mnt.Attach == nil {
		l.Trace("attempted attach mount with no mount definition")
		return nil, ErrInvalDat
	}
	// go-swagger should handle other validation that we need
	var ma *Attach
	if ma, err = API.Attachments.GetOrAttach((*Attach)(mnt.Attach.Attach)); err != nil {
		l.WithError(err).Debug("GetOrAttach failed")
		return nil, err
	}
	mnt.Attach.Attach = (*models.Attach)(ma)
	defer func() {
		if err != nil {
			API.Store.RefAdd(mnt.Attach.Attach.ID, -1)
		}
	}()
	l.WithField("attach_id", mnt.Attach.Attach.ID)
	// ok, we're good to attempt the mount
	// make a mountpoint
	if err = mount.Mount(mnt.Attach.Attach.DeviceFile, mnt.Mountpoint, *mnt.Attach.FsType, mnt.Attach.MountOptions); err != nil {
		l.WithError(err).Error("failed to mount")
		return nil, ErrFail
	}
	return mnt, nil
}

func (m *MountDriverAttach) Unmount(mnt *Mount) (ret *Mount, err error) {
	l := m.log.WithFields(logrus.Fields{
		"operation": "unmount",
		"id":        mnt.ID,
		"attach_id": mnt.Attach.Attach.ID,
	})

	// always lazy unmount.  Good idea?
	if err = mount.Unmount(mnt.Mountpoint, false, true); err != nil {
		l.WithError(err).Error("unmount failed")
		return nil, ErrFail
	}
	API.Store.RefAdd(mnt.Attach.Attach.ID, -1)
	// garbage collection should do our cleanup
	return mnt, nil
}
