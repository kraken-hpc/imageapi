package api

import (
	"io/ioutil"
	"os"

	"github.com/kraken-hpc/imageapi/models"
	"github.com/sirupsen/logrus"
)

var MountDrivers = map[string]MountDriver{}

type Mount models.Mount

// Make sure Mount is an EndpointObject
var _ EndpointObject = (*Mount)(nil)

func (m *Mount) GetID() models.ID                       { return m.ID }
func (m *Mount) SetID(id models.ID)                     { m.ID = id }
func (m *Mount) GetRefs() int64                         { return m.Refs }
func (m *Mount) RefAdd(i int64)                         { m.Refs += i }
func (m *Mount) EndpointObjectType() EndpointObjectType { return EndpointObjectMount }

type Mounts struct {
	log *logrus.Entry
}

// Init initializes the mounts subsystem
func (m *Mounts) Init(log *logrus.Entry) {
	// init driver
	m.log = log
	m.log.Info("initializing mount drivers")
	for name, drv := range MountDrivers {
		m.log.Debugf("initializing driver: %s", name)
		drv.Init(m.log.WithField("driver", name))
	}
	m.log.Info("mount subsystem initialized")
}

// List lists all mounts
func (m *Mounts) List() (ret []*Mount) {
	ret = []*Mount{}
	for _, o := range API.Store.ListType(EndpointObjectMount) {
		ret = append(ret, o.(*Mount))
	}
	return
}

// Get gets a mount by id
func (m *Mounts) Get(id models.ID) *Mount {
	if eo := API.Store.Get(id); eo != nil {
		if ret, ok := eo.(*Mount); ok {
			return ret
		}
	}
	return nil
}

// Mount based on a generic specification
func (m *Mounts) Mount(mnt *Mount) (ret *Mount, err error) {
	l := m.log.WithField("operation", "mount")
	if mnt.ID != 0 {
		l.Errorf("requested a mount with non-zero mount ID")
		return nil, ErrInvalDat
	}
	if drv, ok := MountDrivers[mnt.Kind]; ok {
		// we take responsibility for creating the mountpoint
		l = l.WithField("driver", drv)
		if err = os.MkdirAll(API.MountDir, 0700); err != nil {
			l.WithError(err).Error("could not create base mount directory")
			return nil, ErrSrv
		}
		if mnt.Mountpoint, err = ioutil.TempDir(API.MountDir, "mount_"); err != nil {
			l.WithError(err).Error("failed to create mountpoint")
			return nil, ErrSrv
		}
		defer func() {
			if err != nil {
				if rmerr := os.Remove(mnt.Mountpoint); rmerr != nil {
					// not fatal, but we should report it
					l.WithError(rmerr).Warn("failed to remove mountpoint after failed mount")
				}
			}
		}()
		l = l.WithField("mountpoint", mnt.Mountpoint)
		if chmoderr := os.Chmod(mnt.Mountpoint, os.FileMode(0755)); chmoderr != nil {
			l.WithError(chmoderr).Error("chmod of mountpoint failed")
		}
		ret, err = drv.Mount(mnt)
		if err == nil {
			ret = API.Store.Register(ret).(*Mount)
			l.WithField("id", ret.ID).Info("successfully mounted")
		}
		return
	}
	return nil, ErrNoDrv
}

// GetOrMount gets a mount if it already exists, if it does not, it attempts to mount
func (m *Mounts) GetOrMount(mnt *Mount) (ret *Mount, err error) {
	// existing mount?
	if mnt.ID != 0 {
		gm := m.Get(mnt.ID)
		if gm != nil {
			return gm, nil
		}
		return nil, ErrNotFound
	}
	// new mount?
	return m.Mount(mnt)
}

// Unmount based on a generic specification
func (m *Mounts) Unmount(mnt *Mount, force bool) (ret *Mount, err error) {
	l := m.log.WithField("operation", "unmount")
	if mnt.ID < 1 {
		l.Trace("unmount called with 0 ID")
		return nil, ErrNotFound
	}
	eo := API.Store.Get(mnt.ID)
	if eo == nil {
		l.Tracef("unmount called on non-existent mount ID: %d", mnt.ID)
		return nil, ErrNotFound
	}
	defer func() {
		API.Store.RefAdd(eo.GetID(), -1)
	}()
	var ok bool
	if mnt, ok = eo.(*Mount); !ok {
		l.Trace("unmount called on non-mount object")
		return nil, ErrNotFound
	}
	l = l.WithFields(logrus.Fields{
		"id":         mnt.ID,
		"mountpoint": mnt.Mountpoint,
	})
	if mnt.Refs > 1 && !force { // we hold 1 from our Get above
		l.Debug("unmount called on mount that is in use")
		return nil, ErrBusy
	}
	// two edge cases:
	// 1) mountpoint no longer exists.  Technically, any stat error (could be perms).
	if _, err := os.Stat(mnt.Mountpoint); err != nil {
		API.Store.Unregister(mnt)
		l.Warn("unmount called on a mountpoint that no longer exists, assuming it's already unmounted")
		// to the end user, we treat this as a successful unmount
		return ret, nil
	}
	// 2) mountpoint isn't a mount (maybe we got unmounted already?)
	if !isMountpoint(mnt.Mountpoint) {
		API.Store.Unregister(mnt)
		l.Warn("mountpoint is not mounted, unregistering")
		if rmerr := os.Remove(mnt.Mountpoint); rmerr != nil {
			l.WithError(rmerr).Warn("failed to remove mountpoint on unmount")
		}
		return ret, nil
	}
	if drv, ok := MountDrivers[mnt.Kind]; ok {
		l = l.WithField("driver", drv)
		ret, err = drv.Unmount(mnt)
		if err == nil {
			if rmerr := os.Remove(mnt.Mountpoint); rmerr != nil {
				l.WithError(rmerr).Warn("failed to remove mountpoint on unmount")
			}
			API.Store.Unregister(ret)
			l.Info("successfully unmounted")
		}
		return ret, err
	}
	return nil, ErrNoDrv
}
