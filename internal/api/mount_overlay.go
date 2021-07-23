package api

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/bensallen/rbd/pkg/mount"
	"github.com/kraken-hpc/imageapi/models"
	"github.com/sirupsen/logrus"
)

func init() {
	MountDrivers[models.MountKindOverlay] = &MountDriverOverlay{}
}

type MountDriverOverlay struct {
	log *logrus.Entry
}

func (m *MountDriverOverlay) Init(log *logrus.Entry) {
	m.log = log
	m.log.Trace("initialized")
}

func (m *MountDriverOverlay) Mount(mnt *Mount) (r *Mount, err error) {
	l := m.log.WithField("operation", "mount")
	if mnt.Overlay == nil {
		l.Trace("attempted to overlay mount without overlay definition")
		return nil, ERRINVALDAT
	}

	// there most be at least one lower
	if len(mnt.Overlay.Lower) == 0 {
		l.Debug("no lower mount(s) specified")
		return nil, ERRINVALDAT
	}

	// make sure lower mounts exits, or mount them if we need to
	// warning: there's a possible race here if someone removed these mounts while we're assembling
	//			we might need an extneral interface to lock them.
	lmnts := []string{}
	refs := []models.ID{} // keep track of refs we hold
	defer func() {
		if err != nil { // cleanup on error
			for _, r := range refs {
				API.Store.RefAdd(r, -1)
			}
			if mnt.Overlay.Upperdir != "" {
				os.Remove(mnt.Overlay.Upperdir)
			}
			if mnt.Overlay.Workdir != "" {
				os.Remove(mnt.Overlay.Workdir)
			}
		}
	}()

	// acquire lower mounts
	for i := range mnt.Overlay.Lower {
		lmnt, err := API.Mounts.GetOrMount((*Mount)(mnt.Overlay.Lower[i]))
		if err != nil {
			l.WithError(err).Debug("lower mount GetOrMount failed for overlay mount")
			return nil, err
		}
		mnt.Overlay.Lower[i] = (*models.Mount)(lmnt)
		refs = append(refs, lmnt.ID)
		lmnts = append(lmnts, lmnt.Mountpoint)
	}

	// ok, we're good to attempt the mount
	// make a upperdir/workdir
	if mnt.Overlay.Upperdir, err = ioutil.TempDir(API.MountDir, "upper_"); err != nil {
		l.WithError(err).Error("could not create upperdir")
		return nil, ERRSRV
	}
	if chmoderr := os.Chmod(mnt.Overlay.Upperdir, os.FileMode(0755)); chmoderr != nil {
		l.WithError(chmoderr).Error("failed to chmod upperdir")
	}
	if mnt.Overlay.Workdir, err = ioutil.TempDir(API.MountDir, "work_"); err != nil {
		l.WithError(err).Error("could not create workdir")
		return nil, ERRSRV
	}
	if chmoderr := os.Chmod(mnt.Overlay.Workdir, os.FileMode(0755)); chmoderr != nil {
		l.WithError(chmoderr).Error("failed to chmod workdir")
	}

	// try the mounmt
	opts := []string{
		"lowerdir=" + strings.Join(lmnts, ":"),
		"upperdir=" + mnt.Overlay.Upperdir,
		"workdir=" + mnt.Overlay.Workdir,
	}
	l.WithField("opts", opts)
	if err = mount.Mount("overlay", mnt.Mountpoint, "overlay", opts); err != nil {
		l.WithError(err).Error("overlay mount failed")
		return nil, ERRFAIL
	}
	l.Info("successfully mounted")
	return mnt, nil
}

func (m *MountDriverOverlay) Unmount(mnt *Mount) (ret *Mount, err error) {
	l := m.log.WithFields(logrus.Fields{
		"operation": "unmount",
		"id":        mnt.ID,
	})

	// always lazy unmount.  Good idea?
	if err = mount.Unmount(mnt.Mountpoint, false, true); err != nil {
		l.WithError(err).Error("unmount failed")
		return nil, ERRFAIL
	}

	os.RemoveAll(mnt.Overlay.Workdir)  // option to leave behind?
	os.RemoveAll(mnt.Overlay.Upperdir) // option to leave behind? Or store on RBD?
	for _, l := range mnt.Overlay.Lower {
		API.Store.RefAdd(l.ID, -1)
	}
	l.Info("successfully unmounted")
	return mnt, nil
}
