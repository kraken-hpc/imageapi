package api

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/bensallen/rbd/pkg/mount"
	"github.com/jlowellwofford/imageapi/models"
)

type MountsRBDType struct {
	next  models.ID
	mnts  map[models.ID]*models.MountRbd
	mutex *sync.Mutex
}

func (m *MountsRBDType) Init() {
	m.next = 1 // 0 == unspecified
	m.mnts = make(map[models.ID]*models.MountRbd)
	m.mutex = &sync.Mutex{}
}

func (m *MountsRBDType) Mount(mnt *models.MountRbd) (ret *models.MountRbd, err error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	// does the mount already exist?
	if _, ok := m.mnts[mnt.ID]; ok {
		return nil, fmt.Errorf("mount failure: moint already exists: %d", mnt.ID)
	}
	// make sure the dev exists, or attach it
	var rbd *models.Rbd
	if mnt.RbdID != 0 { // Rbd was specified by ID
		if rbd, err = Rbds.Get(mnt.RbdID); err != nil {
			return nil, ERRNOTFOUND
		}
		Rbds.RefAdd(rbd.ID, 1)
	} else if mnt.Rbd != nil { // try to attach it
		if rbd, err = Rbds.Map(mnt.Rbd); err != nil {
			return nil, fmt.Errorf("failed to attach underlying RBD image: %v", err)
		}
		mnt.RbdID = rbd.ID
	} else { // unspecified
		return nil, fmt.Errorf("no rbd specified")
	}
	defer func() {
		if err != nil {
			Rbds.RefAdd(rbd.ID, -1)
		}
	}()
	// ok, we're good to attempt the mount
	// make a mountpoint
	if err = os.MkdirAll(mountDir, 0700); err != nil {
		return nil, fmt.Errorf("could not create base mount directory: %v", err)
	}
	if mnt.Mountpoint, err = ioutil.TempDir(mountDir, "mount_"); err != nil {
		return nil, fmt.Errorf("could not create mountpoint: %v", err)
	}
	if err = mount.Mount(rbd.DeviceFile, mnt.Mountpoint, *mnt.FsType, mnt.MountOptions); err != nil {
		return nil, fmt.Errorf("mount failure: %v", err)
	}
	mnt.ID = m.next
	mnt.Refs = 1
	m.next++
	m.mnts[mnt.ID] = mnt
	return mnt, nil
}

func (m *MountsRBDType) Unmount(id models.ID) (ret *models.MountRbd, err error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var mnt *models.MountRbd
	var ok bool

	if mnt, ok = m.mnts[id]; !ok {
		return nil, ERRNOTFOUND
	}

	if mnt.Refs > 0 {
		return nil, fmt.Errorf("unmount failure: mount is in use")
	}

	// always lazy unmount.  Good idea?
	if err = mount.Unmount(mnt.Mountpoint, false, true); err != nil {
		return nil, fmt.Errorf("unmount failure: %v", err)
	}
	os.Remove(mnt.Mountpoint) // we shouldn't fail on this. Should we report it anyway?
	delete(m.mnts, id)
	Rbds.RefAdd(mnt.RbdID, -1)
	// garbage collection should do our cleanup
	return
}

func (m *MountsRBDType) Get(id models.ID) (mnt *models.MountRbd, err error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	var ok bool
	if mnt, ok = m.mnts[id]; !ok {
		return nil, ERRNOTFOUND
	}
	return
}

func (m *MountsRBDType) List() (mnts []*models.MountRbd) {
	for _, i := range m.mnts {
		mnts = append(mnts, i)
	}
	return
}

func (m *MountsRBDType) RefAdd(id models.ID, n int64) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if mnt, ok := m.mnts[id]; ok {
		mnt.Refs += n
	}
}

// Collect will run garbage collection on any RBDs with ref == 0
func (m *MountsRBDType) Collect() {
	list := []models.ID{}
	m.mutex.Lock()
	for _, mnt := range m.mnts {
		if mnt.Refs == 0 {
			// let's collect
			list = append(list, mnt.ID)
		}
	}
	m.mutex.Unlock()
	for _, id := range list {
		m.Unmount(id)
	}
}
