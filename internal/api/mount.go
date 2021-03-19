package api

import (
	"fmt"

	"github.com/jlowellwofford/imageapi/models"
)

// API operations on generic mountpoints
// This essentially acts as a switcher for rbd/overlay

// List all mounts
func ListMounts() (ret []*models.Mount) {
	ret = []*models.Mount{}
	strRbd := "rbd"
	for _, m := range MountsRbd.List() {
		ret = append(ret, &models.Mount{
			Kind:    &strRbd,
			MountID: m.ID,
			Rbd:     m,
		})
	}
	strOverlay := "overlay"
	for _, m := range MountsOverlay.List() {
		ret = append(ret, &models.Mount{
			Kind:    &strOverlay,
			MountID: m.ID,
			Overlay: m,
		})
	}
	return
}

// Mount based on a generic specification
func Mount(mnt *models.Mount) (ret *models.Mount, err error) {
	if mnt.MountID != 0 { // we can't specify an ID
		return nil, fmt.Errorf("disallowed mount_id was specified when trying to mount")
	}
	switch *mnt.Kind {
	case "rbd":
		if mnt.Rbd == nil {
			return nil, fmt.Errorf("rbd kind was requested, but no rbd specification was provided")
		}
		if mnt.Rbd, err = MountsRbd.Mount(mnt.Rbd); err != nil {
			return nil, err
		}
	case "overlay":
		if mnt.Overlay == nil {
			return nil, fmt.Errorf("overlay kind was requested, but no overlay specification was provided")
		}
		if mnt.Overlay, err = MountsOverlay.Mount(mnt.Overlay); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown mount kind: %s", *mnt.Kind)
	}
	ret = mnt
	return
}

// Unmount based on a generic specification
func Unmount(mnt *models.Mount) (ret *models.Mount, err error) {
	id := mnt.MountID
	switch *mnt.Kind {
	case "rbd":
		if id == 0 {
			if mnt.Rbd == nil || mnt.Rbd.ID == 0 {
				return nil, fmt.Errorf("no mount id specified")
			}
			id = mnt.Rbd.ID
		}
		if mnt.Rbd, err = MountsRbd.Unmount(id); err != nil {
			return nil, err
		}
	case "overlay":
		if id == 0 {
			if mnt.Overlay == nil || mnt.Overlay.ID == 0 {
				return nil, fmt.Errorf("no mount id specified")
			}
			id = mnt.Overlay.ID
		}
		if mnt.Overlay, err = MountsOverlay.Unmount(id); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown mount kind: %s", *mnt.Kind)
	}
	ret = mnt
	return
}

func MountGetMountpoint(mnt *models.Mount) (mntpt string, err error) {
	id := mnt.MountID
	switch *mnt.Kind {
	case "rbd":
		var rbd *models.MountRbd
		if id == 0 {
			if mnt.Rbd == nil || mnt.Rbd.ID == 0 {
				return "", fmt.Errorf("no mount id specified")
			}
			id = mnt.Rbd.ID
		}
		if rbd, err = MountsRbd.Get(id); err != nil {
			return "", err
		}
		mntpt = rbd.Mountpoint
	case "overlay":
		var overlay *models.MountOverlay
		if id == 0 {
			if mnt.Overlay == nil || mnt.Overlay.ID == 0 {
				return "", fmt.Errorf("no mount id specified")
			}
			id = mnt.Overlay.ID
		}
		if overlay, err = MountsOverlay.Get(id); err != nil {
			return "", err
		}
		mntpt = overlay.Mountpoint
	default:
		return "", fmt.Errorf("unknown mount kind: %s", *mnt.Kind)
	}
	return
}

func MountRefAdd(mnt *models.Mount, n int64) {
	switch *mnt.Kind {
	case "rbd":
		if mnt.MountID == 0 {
			mnt.MountID = mnt.Rbd.ID
		}
		MountsRbd.RefAdd(mnt.MountID, n)
	case "overlay":
		if mnt.MountID == 0 {
			mnt.MountID = mnt.Overlay.ID
		}
		MountsOverlay.RefAdd(mnt.MountID, n)
	}
}
