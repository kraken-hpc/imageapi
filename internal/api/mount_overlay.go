package api

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/bensallen/rbd/pkg/mount"
	"github.com/jlowellwofford/imageapi/models"
)

type MountsOverlayType struct {
	next  models.ID
	mnts  map[models.ID]*models.MountOverlay
	mutex *sync.Mutex
}

func (m *MountsOverlayType) Init() {
	m.next = 1 // 0 == unspecified
	m.mnts = make(map[models.ID]*models.MountOverlay)
	m.mutex = &sync.Mutex{}
}

func (m *MountsOverlayType) List() (r []*models.MountOverlay) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for _, mnt := range m.mnts {
		r = append(r, mnt)
	}
	return
}

func (m *MountsOverlayType) Get(id models.ID) (*models.MountOverlay, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if r, ok := m.mnts[id]; ok {
		return r, nil
	}
	return nil, ERRNOTFOUND
}

func (m *MountsOverlayType) Mount(mnt *models.MountOverlay) (r *models.MountOverlay, err error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// note: we don't check for existence, because we will happily make more than one overlay of the same thing

	// there most be at least one lower
	if len(mnt.Lower) == 0 {
		return nil, fmt.Errorf("at least one lower mount must be specified")
	}

	// make sure lower mounts exits, or mount them if we need to
	// warning: there's a possible race here if someone removed these mounts while we're assembling
	//			we might need an extneral interface to lock them.
	lmnts := []string{}
	for i := range mnt.Lower {
		if mnt.Lower[i].MountID == 0 { // try to mount
			if mnt.Lower[i], err = Mount(mnt.Lower[i]); err != nil {
				return nil, fmt.Errorf("failed to mount lower mount: %v", err)
			}
		} else {
			MountRefAdd(mnt.Lower[i], 1) // we allow this to silently fail
		}
		var mntpt string
		if mntpt, err = MountGetMountpoint(mnt.Lower[i]); err != nil {
			return nil, fmt.Errorf("failed to get mountpoint for lower mount: %v", err)
		}
		lmnts = append(lmnts, mntpt)
	}
	defer func() {
		if err != nil {
			for _, m := range mnt.Lower {
				MountRefAdd(m, -1)
			}
		}
	}()

	// ok, we're good to attempt the mount
	// make a mountpoint/upperdir/workdir
	if err = os.MkdirAll(mountDir, 0700); err != nil {
		return nil, fmt.Errorf("could not create base mount directory: %v", err)
	}
	if mnt.Mountpoint, err = ioutil.TempDir(mountDir, "mount_"); err != nil {
		return nil, fmt.Errorf("could not create mountpoint: %v", err)
	}
	os.Chmod(mnt.Mountpoint, os.FileMode(0755))
	if mnt.Upperdir, err = ioutil.TempDir(mountDir, "upper_"); err != nil {
		return nil, fmt.Errorf("could not create upperdir: %v", err)
	}
	os.Chmod(mnt.Upperdir, os.FileMode(0755))
	if mnt.Workdir, err = ioutil.TempDir(mountDir, "work_"); err != nil {
		return nil, fmt.Errorf("could not create workdir: %v", err)
	}
	os.Chmod(mnt.Workdir, os.FileMode(0755))

	// try the mounmt
	opts := []string{
		"lowerdir=" + strings.Join(lmnts, ":"),
		"upperdir=" + mnt.Upperdir,
		"workdir=" + mnt.Workdir,
	}
	if err = mount.Mount("overlay", mnt.Mountpoint, "overlay", opts); err != nil {
		return nil, fmt.Errorf("mount failure: %v", err)
	}

	// store
	mnt.ID = m.next
	mnt.Refs = 1
	m.next++
	m.mnts[mnt.ID] = mnt

	return mnt, nil
}

func (m *MountsOverlayType) Unmount(id models.ID) (mnt *models.MountOverlay, err error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var ok bool

	if mnt, ok = m.mnts[id]; !ok {
		return nil, ERRNOTFOUND
	}

	if mnt.Refs > 1 {
		return nil, fmt.Errorf("unmount failure: mount is in use")
	}

	// always lazy unmount.  Good idea?
	if err = mount.Unmount(mnt.Mountpoint, false, true); err != nil {
		return nil, fmt.Errorf("unmount failure: %v", err)
	}

	os.Remove(mnt.Mountpoint)  // we shouldn't fail on this. Should we report it anyway?
	os.RemoveAll(mnt.Workdir)  // option to leave behind?
	os.RemoveAll(mnt.Upperdir) // option to leave behind? Or store on RBD?
	delete(m.mnts, id)
	for _, l := range mnt.Lower {
		MountRefAdd(l, -1)
		// garbage collection will handle cleanup
	}
	return
}

func (m *MountsOverlayType) RefAdd(id models.ID, n int64) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if r, ok := m.mnts[id]; ok {
		r.Refs += n
	}
}

// Collect will run garbage collection on any RBDs with ref == 0
func (m *MountsOverlayType) Collect() {
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
