package api

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"sync"

	"github.com/bensallen/rbd/pkg/mount"
	"github.com/jlowellwofford/imageapi/models"
)

type MountsRBDType struct {
	mnts  map[int64]*models.MountRbd
	mutex *sync.Mutex
}

func (m *MountsRBDType) Init() {
	m.mnts = make(map[int64]*models.MountRbd)
	m.mutex = &sync.Mutex{}
}

func (m *MountsRBDType) Mount(mnt *models.MountRbd) (ret *models.MountRbd, err error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	// does the mount already exist?
	if _, ok := m.mnts[*mnt.ID]; ok {
		return nil, fmt.Errorf("mount failure: moint already exists for device %d", *mnt.ID)
	}
	// make sure the dev exists/is ours
	if _, err = Rbds.Get(*mnt.ID); err != nil {
		return nil, fmt.Errorf("mount failure: %v", err)
	}
	// ok, we're good to attempt the mount
	// make a mountpoint
	if err = os.MkdirAll(mountDir, 0700); err != nil {
		return nil, fmt.Errorf("could not create base mount directory: %v", err)
	}
	if mnt.Mountpoint, err = ioutil.TempDir(mountDir, "mount_"); err != nil {
		return nil, fmt.Errorf("could not create mountpoint: %v", err)
	}
	dev := "/dev/rbd" + strconv.FormatInt(*mnt.ID, 10)
	if err = mount.Mount(dev, mnt.Mountpoint, *mnt.FsType, mnt.MountOptions); err != nil {
		return nil, fmt.Errorf("mount failure: %v", err)
	}
	m.mnts[*mnt.ID] = mnt
	Rbds.RefAdd(*mnt.ID, 1)
	return mnt, nil
}

func (m *MountsRBDType) Unmount(id int64) (err error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var mnt *models.MountRbd
	var ok bool

	if mnt, ok = m.mnts[id]; !ok {
		return fmt.Errorf("unmount failure: no such device %d", id)
	}

	if mnt.Ref > 0 {
		return fmt.Errorf("unmount failure: mount is in use")
	}

	// always lazy unmount.  Good idea?
	if err = mount.Unmount(mnt.Mountpoint, false, true); err != nil {
		return fmt.Errorf("unmount failure: %v", err)
	}
	os.Remove(mnt.Mountpoint) // we shouldn't fail on this. Should we report it anyway?
	delete(m.mnts, id)
	Rbds.RefAdd(id, -1)
	return
}

func (m *MountsRBDType) Get(id int64) (mnt *models.MountRbd, err error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	var ok bool
	if mnt, ok = m.mnts[id]; !ok {
		return nil, fmt.Errorf("rbd mount does not exist: %d", id)
	}
	return
}

func (m *MountsRBDType) List() (mnts []*models.MountRbd) {
	for _, i := range m.mnts {
		mnts = append(mnts, i)
	}
	return
}

func (m *MountsRBDType) RefAdd(id, n int64) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if mnt, ok := m.mnts[id]; ok {
		mnt.Ref += n
	}
}
