package api

import (
	"fmt"
	"net"
	"strings"

	"github.com/bensallen/rbd/pkg/mount"
	"github.com/kraken-hpc/imageapi/models"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
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
	// Resolve host
	ips, err := net.LookupIP(*mnt.Nfs.Host)
	if err != nil {
		l.WithError(err).Debug("host resolution failed")
		return nil, ERRINVALDAT
	}
	ip := ips[0] // should we be smarter about this?
	// ok, we're good to attempt the mount
	// make a mountpoint
	flags := uintptr(0)
	if mnt.Nfs.Ro != nil && *mnt.Nfs.Ro {
		flags = unix.MS_RDONLY
	}
	version := "4.2"
	if mnt.Nfs.Version != nil {
		version = *mnt.Nfs.Version
	}
	mnt.Nfs.Options = append(mnt.Nfs.Options, fmt.Sprintf("addr=%s", ip.String()))
	mnt.Nfs.Options = append(mnt.Nfs.Options, fmt.Sprintf("clientaddr=%s", ip.String()))
	mnt.Nfs.Options = append(mnt.Nfs.Options, fmt.Sprintf("vers=%s", version))
	// this doesn't work because u-root Mount is broken
	//if err = mount.Mount(fmt.Sprintf("%s:%s", *mnt.Nfs.Host, *mnt.Nfs.Path), mnt.Mountpoint, "nfs", mnt.Nfs.MountOptions); err != nil {
	if err = unix.Mount(fmt.Sprintf("%s:%s", *mnt.Nfs.Host, *mnt.Nfs.Path), mnt.Mountpoint, "nfs", flags, strings.Join(mnt.Nfs.Options, ",")); err != nil {
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
		"host":      mnt.Nfs.Host,
		"path":      mnt.Nfs.Path,
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
