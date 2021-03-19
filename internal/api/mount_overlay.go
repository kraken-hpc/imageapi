package api

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/bensallen/rbd/pkg/mount"
	"github.com/jlowellwofford/imageapi/models"
	"github.com/sirupsen/logrus"
)

type MountsOverlayType struct {
	next  models.ID
	mnts  map[models.ID]*models.MountOverlay
	mutex *sync.Mutex
	log   *logrus.Entry
}

func (m *MountsOverlayType) Init() {
	m.next = 1 // 0 == unspecified
	m.mnts = make(map[models.ID]*models.MountOverlay)
	m.mutex = &sync.Mutex{}
	m.log = Log.WithField("subsys", "mount_overlay")
	m.log.Trace("initialized")
}

func (m *MountsOverlayType) List() (r []*models.MountOverlay) {
	l := m.log.WithField("operation", "list")
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for _, mnt := range m.mnts {
		r = append(r, mnt)
	}
	l.WithField("entries", len(r)).Trace("listing entries")
	return
}

func (m *MountsOverlayType) Get(id models.ID) (*models.MountOverlay, error) {
	l := m.log.WithFields(logrus.Fields{
		"operation": "get",
		"id":        id,
	})
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if r, ok := m.mnts[id]; ok {
		l.Trace("found")
		return r, nil
	}
	l.Trace("not found")
	return nil, ERRNOTFOUND
}

func (m *MountsOverlayType) Mount(mnt *models.MountOverlay) (r *models.MountOverlay, err error) {
	l := m.log.WithField("operation", "mount")
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// note: we don't check for existence, because we will happily make more than one overlay of the same thing

	// there most be at least one lower
	if len(mnt.Lower) == 0 {
		l.Debug("no lower mount(s) specified")
		return nil, fmt.Errorf("at least one lower mount must be specified")
	}

	// make sure lower mounts exits, or mount them if we need to
	// warning: there's a possible race here if someone removed these mounts while we're assembling
	//			we might need an extneral interface to lock them.
	lmnts := []string{}
	for i := range mnt.Lower {
		if mnt.Lower[i].MountID == 0 { // try to mount
			if mnt.Lower[i], err = Mount(mnt.Lower[i]); err != nil {
				// clear refs for mounts we already processed
				for j := 0; j < i; j++ {
					MountRefAdd(mnt.Lower[i], -1)
				}
				l.WithError(err).Error("lower mount failed")
				return nil, fmt.Errorf("failed to mount lower mount: %v", err)
			}
		} else {
			MountRefAdd(mnt.Lower[i], 1) // we allow this to silently fail
		}
		var mntpt string
		if mntpt, err = MountGetMountpoint(mnt.Lower[i]); err != nil {
			// this means that the mount doesn't exist...
			for j := 0; j < i; j++ {
				MountRefAdd(mnt.Lower[i], -1)
			}
			l.WithError(err).Error("failed to get mountpoint for lower mount")
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
		l.WithError(err).Error("could not create base directory")
		return nil, fmt.Errorf("could not create base mount directory: %v", err)
	}
	if mnt.Mountpoint, err = ioutil.TempDir(mountDir, "mount_"); err != nil {
		l.WithError(err).Error("could not create mountpoint")
		return nil, fmt.Errorf("could not create mountpoint: %v", err)
	}
	os.Chmod(mnt.Mountpoint, os.FileMode(0755))
	if mnt.Upperdir, err = ioutil.TempDir(mountDir, "upper_"); err != nil {
		l.WithError(err).Error("could not create upperdir")
		return nil, fmt.Errorf("could not create upperdir: %v", err)
	}
	os.Chmod(mnt.Upperdir, os.FileMode(0755))
	if mnt.Workdir, err = ioutil.TempDir(mountDir, "work_"); err != nil {
		l.WithError(err).Error("could not create workdir")
		return nil, fmt.Errorf("could not create workdir: %v", err)
	}
	os.Chmod(mnt.Workdir, os.FileMode(0755))

	// try the mounmt
	opts := []string{
		"lowerdir=" + strings.Join(lmnts, ":"),
		"upperdir=" + mnt.Upperdir,
		"workdir=" + mnt.Workdir,
	}
	l.WithField("opts", opts)
	if err = mount.Mount("overlay", mnt.Mountpoint, "overlay", opts); err != nil {
		l.WithError(err).Error("mount failed")
		return nil, fmt.Errorf("mount failure: %v", err)
	}

	// store
	mnt.ID = m.next
	mnt.Refs = 1
	m.next++
	m.mnts[mnt.ID] = mnt
	l.Info("successfully mounted")

	return mnt, nil
}

func (m *MountsOverlayType) Unmount(id models.ID) (mnt *models.MountOverlay, err error) {
	l := m.log.WithFields(logrus.Fields{
		"operation": "unmount",
		"id":        id,
	})
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var ok bool

	if mnt, ok = m.mnts[id]; !ok {
		l.Debug("mount does not exist")
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

	os.Remove(mnt.Mountpoint)  // we shouldn't fail on this. Should we report it anyway?
	os.RemoveAll(mnt.Workdir)  // option to leave behind?
	os.RemoveAll(mnt.Upperdir) // option to leave behind? Or store on RBD?
	delete(m.mnts, id)
	for _, l := range mnt.Lower {
		MountRefAdd(l, -1)
		// garbage collection will handle cleanup
	}
	l.Info("successfully unmounted")
	return
}

func (m *MountsOverlayType) RefAdd(id models.ID, n int64) {
	l := m.log.WithField("operation", "refadd")
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if r, ok := m.mnts[id]; ok {
		l.WithField("n", n).Trace("added")
		r.Refs += n
	} else {
		l.Debug("no such rbd")
	}
}

// Collect will run garbage collection on any RBDs with ref == 0
func (m *MountsOverlayType) Collect() {
	l := m.log.WithField("operation", "collect")
	list := []models.ID{}
	m.mutex.Lock()
	for _, mnt := range m.mnts {
		if mnt.Refs <= 0 {
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
