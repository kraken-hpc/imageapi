// mounts attachments
package mount

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/bensallen/rbd/pkg/mount"
	"github.com/kraken-hpc/imageapi/models"
	"github.com/sirupsen/logrus"
)

type MountsAttachType struct {
	next  models.ID
	mnts  map[models.ID]*models.MountRbd
	mutex *sync.Mutex
	log   *logrus.Entry
}

func (m *MountsRBDType) Init() {
	m.next = 1 // 0 == unspecified
	m.mnts = make(map[models.ID]*models.MountRbd)
	m.mutex = &sync.Mutex{}
	m.log = Log.WithField("subsys", "mount_rbd")
	m.log.Trace("initialized")
}

func (m *MountsRBDType) Mount(mnt *models.MountRbd) (ret *models.MountRbd, err error) {
	l := m.log.WithField("operation", "mount")
	m.mutex.Lock()
	defer m.mutex.Unlock()
	// does the mount already exist?
	if _, ok := m.mnts[mnt.ID]; ok {
		l.Debug("requested mount of already existing mount")
		return nil, fmt.Errorf("mount failure: moint already exists: %d", mnt.ID)
	}
	// make sure the dev exists, or attach it
	var rbd *models.Rbd
	if mnt.RbdID != 0 { // Rbd was specified by ID
		if rbd, err = Rbds.Get(mnt.RbdID); err != nil {
			l.WithError(err).Debug("rbd does not exist")
			return nil, ERRNOTFOUND
		}
		Rbds.RefAdd(rbd.ID, 1)
	} else if mnt.Rbd != nil { // try to attach it
		if rbd, err = Rbds.Map(mnt.Rbd); err != nil {
			l.WithError(err).Error("rbd map failed")
			return nil, fmt.Errorf("failed to attach underlying RBD image: %v", err)
		}
		mnt.RbdID = rbd.ID
	} else { // unspecified
		l.Error("no rbd specified")
		return nil, fmt.Errorf("no rbd specified")
	}
	defer func() {
		if err != nil {
			Rbds.RefAdd(rbd.ID, -1)
		}
	}()
	l.WithField("rbd", rbd.ID)
	// HERE!!!!
	// ok, we're good to attempt the mount
	// make a mountpoint
	if err = os.MkdirAll(MountDir, 0700); err != nil {
		l.WithError(err).Error("failed to make mount directory")
		return nil, fmt.Errorf("could not create base mount directory: %v", err)
	}
	if mnt.Mountpoint, err = ioutil.TempDir(MountDir, "mount_"); err != nil {
		l.WithError(err).Error("failed to make pointpoint")
		return nil, fmt.Errorf("could not create mountpoint: %v", err)
	}
	if err = mount.Mount(rbd.DeviceFile, mnt.Mountpoint, *mnt.FsType, mnt.MountOptions); err != nil {
		l.WithError(err).Error("failed to mount")
		return nil, fmt.Errorf("mount failure: %v", err)
	}
	mnt.ID = m.next
	mnt.Refs = 1
	m.next++
	m.mnts[mnt.ID] = mnt
	l.Info("successfully mounted")
	return mnt, nil
}

func (m *MountsRBDType) Unmount(id models.ID) (ret *models.MountRbd, err error) {
	l := m.log.WithFields(logrus.Fields{
		"operation": "unmount",
		"id":        id,
	})
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var mnt *models.MountRbd
	var ok bool

	if mnt, ok = m.mnts[id]; !ok {
		l.Debug("mount does not exists")
		return nil, ERRNOTFOUND
	}

	if mnt.Refs > 0 {
		l.WithField("refs", mnt.Refs).Debug("nonzero refcount")
		return nil, fmt.Errorf("unmount failure: mount is in use")
	}

	// always lazy unmount.  Good idea?
	if err = mount.Unmount(mnt.Mountpoint, false, true); err != nil {
		l.WithError(err).Error("unmount failed")
		return nil, fmt.Errorf("unmount failure: %v", err)
	}
	os.Remove(mnt.Mountpoint) // we shouldn't fail on this. Should we report it anyway?
	delete(m.mnts, id)
	Rbds.RefAdd(mnt.RbdID, -1)
	l.Info("successfully unmounted")
	// garbage collection should do our cleanup
	return
}

func (m *MountsRBDType) Get(id models.ID) (mnt *models.MountRbd, err error) {
	l := m.log.WithFields(logrus.Fields{
		"operation": "get",
		"id":        id,
	})
	m.mutex.Lock()
	defer m.mutex.Unlock()
	var ok bool
	if mnt, ok = m.mnts[id]; !ok {
		l.Trace("found")
		return nil, ERRNOTFOUND
	}
	l.Trace("not found")
	return
}

func (m *MountsRBDType) List() (mnts []*models.MountRbd) {
	l := m.log.WithField("operation", "list")
	for _, i := range m.mnts {
		mnts = append(mnts, i)
	}
	l.WithField("entries", len(mnts)).Trace("listing entries")
	return
}

func (m *MountsRBDType) RefAdd(id models.ID, n int64) {
	l := m.log.WithField("operation", "refadd")
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if mnt, ok := m.mnts[id]; ok {
		l.WithField("n", n).Trace("added")
		mnt.Refs += n
	} else {
		l.Debug("no such rbd")
	}
}

// Collect will run garbage collection on any RBDs with ref == 0
func (m *MountsRBDType) Collect() {
	l := m.log.WithField("operation", "collect")
	list := []models.ID{}
	m.mutex.Lock()
	for _, mnt := range m.mnts {
		if mnt.Refs == 0 {
			// let's collect
			list = append(list, mnt.ID)
		}
	}
	m.mutex.Unlock()
	if len(list) > 0 {
		l.WithField("collectIDs", list).Debug("collecting")
		for _, id := range list {
			m.Unmount(id)
		}
	}
}
