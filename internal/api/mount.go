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
		return nil, ERRINVALDAT
	}
	if drv, ok := MountDrivers[mnt.Kind]; ok {
		// we take responsibility for creating the mountpoint
		l = l.WithField("driver", drv)
		if err = os.MkdirAll(API.MountDir, 0700); err != nil {
			l.WithError(err).Error("could not create base mount directory")
			return nil, ERRSRV
		}
		if mnt.Mountpoint, err = ioutil.TempDir(API.MountDir, "mount_"); err != nil {
			l.WithError(err).Error("failed to create mountpoint")
			return nil, ERRSRV
		}
		if chmoderr := os.Chmod(mnt.Mountpoint, os.FileMode(0755)); chmoderr != nil {
			l.WithError(chmoderr).Error("chmod of mountpoint failed")
		}
		ret, err = drv.Mount(mnt)
		if err == nil {
			ret = API.Store.Register(ret).(*Mount)
			l.Infof("successfully mounted: %s", ret.Mountpoint)
		} else {
			// cleanup mountpoint
			if rmerr := os.Remove(mnt.Mountpoint); rmerr != nil {
				// not fatal, but we should report it
				l.WithError(rmerr).Warn("failed to remove mountpoint after failed mount")
			}
		}
		return
	}
	return nil, ERRNODRV
}

// GetOrMount gets a mount if it already exists, if it does not, it attempts to mount
func (m *Mounts) GetOrMount(mnt *Mount) (ret *Mount, err error) {
	// existing mount?
	if mnt.ID != 0 {
		gm := m.Get(mnt.ID)
		if gm != nil {
			return gm, nil
		}
		return nil, ERRNOTFOUND
	}
	// new mount?
	return m.Mount(mnt)
}

// Unmount based on a generic specification
func (m *Mounts) Unmount(mnt *Mount, force bool) (ret *Mount, err error) {
	l := m.log.WithField("operation", "unmount")
	if mnt.ID < 1 {
		l.Trace("unmount called with 0 ID")
		return nil, ERRNOTFOUND
	}
	eo := API.Store.Get(mnt.ID)
	if eo == nil {
		l.Tracef("unmount called on non-existent mount ID: %d", mnt.ID)
		return nil, ERRNOTFOUND
	}
	defer func() {
		API.Store.RefAdd(eo.GetID(), -1)
	}()
	var ok bool
	if mnt, ok = eo.(*Mount); !ok {
		l.Trace("unmount called on non-mount object")
		return nil, ERRNOTFOUND
	}
	l = l.WithField("id", mnt.ID)
	if mnt.Refs > 1 && !force { // we hold 1 from our Get above
		l.Debug("unmount called on mount that is in use")
		return nil, ERRBUSY
	}
	if drv, ok := MountDrivers[mnt.Kind]; ok {
		l = l.WithField("driver", drv)
		ret, err = drv.Unmount(mnt)
		if err == nil {
			if rmerr := os.Remove(mnt.Mountpoint); rmerr != nil {
				l.WithError(rmerr).Warn("failed to remove mountpoint on unmount")
			}
			API.Store.Unregister(ret)
			l.Infof("successfully unmounted: %s", ret.Mountpoint)
		}
		return ret, err
	}
	return nil, ERRNODRV
}
