package attach

// API operations on rbd maps

import (
	"fmt"
	"sync"

	"github.com/bensallen/rbd/pkg/krbd"
	"github.com/kraken-hpc/imageapi/internal/api/types"
	"github.com/kraken-hpc/imageapi/models"
	"github.com/sirupsen/logrus"
)

var _ types.AttachModel = (*RbdModel)(nil)

type RbdModel models.Rbd

// ID for the attachment
func (m *RbdModel) GetID() models.ID {
	return m.ID
}

// Refs for the attachment
func (m *RbdModel) GetRefs() int64 {
	return m.Refs
}

// Device file for the attachment
func (m *RbdModel) GetDevice() string {
	return m.DeviceFile
}

var _ types.Attach = (*RbdsType)(nil)

type RbdsType struct {
	next  models.ID
	rbds  map[models.ID]*RbdModel
	mutex *sync.Mutex
	log   *logrus.Entry
}

func (r *RbdsType) Init() {
	r.next = 1 // starting from 1 means 0 == unspecified
	r.rbds = make(map[models.ID]*RbdModel)
	r.mutex = &sync.Mutex{}
	r.log = Log.WithField("subsys", "rbd")
	r.log.Trace("initialized")
}

func (r *RbdsType) List() (result []types.Model) {
	l := r.log.WithField("operation", "list")
	r.mutex.Lock()
	defer r.mutex.Unlock()
	for _, m := range r.rbds {
		result = append(result, (*RbdModel)(m))
	}
	l.WithField("entries", len(result)).Trace("listing entries")
	return
}

func (r *RbdsType) Attach(rbdm types.AttachModel) (m types.AttachModel, err error) {
	rbd, ok := rbdm.(*RbdModel)
	if !ok {
		return nil, types.ERRINVALDAT
	}
	r.mutex.Lock()
	defer r.mutex.Unlock()
	// sanity check
	l := r.log.WithField("operation", "map")
	if len(rbd.Monitors) == 0 || *rbd.Pool == "" || *rbd.Image == "" || rbd.Options.Name == "" || rbd.Options.Secret == "" {
		r.log.Debug("incorrect options")
		return nil, fmt.Errorf("the following are required: 1 or more monitors, pool, image, options/name, options/secret")
	}
	l = l.WithFields(logrus.Fields{
		"image":    *rbd.Image,
		"pool":     *rbd.Pool,
		"name":     rbd.Options.Name,
		"monitors": rbd.Monitors,
	})
	w, err := krbd.RBDBusAddWriter()
	if err != nil {
		l.WithError(err).Error("failed to get krbd bus writer")
		return nil, fmt.Errorf("krbd error: %v", err)
	}
	defer w.Close()

	// We allow this because we get free IPv4 format checking
	mons := []string{}
	for _, m := range rbd.Monitors {
		mons = append(mons, m.String())
	}

	i := krbd.Image{
		Monitors: mons,
		Pool:     *rbd.Pool,
		Image:    *rbd.Image,
		Snapshot: rbd.Snapshot,
		Options: &krbd.Options{
			ReadOnly:  rbd.Options.Ro,
			Name:      rbd.Options.Name,
			Secret:    rbd.Options.Secret,
			Namespace: rbd.Options.Namespace,
		},
	}
	// make sure ID doesn't already exist
	dev := krbd.Device{Image: i.Image, Pool: i.Pool, Namespace: i.Options.Namespace, Snapshot: i.Snapshot}

	if err := dev.Find(); err == nil {
		l.Debug("tried to map device that already exists")
		return nil, fmt.Errorf("rbd device already exists")
	}
	// map the rbd
	if err := i.Map(w); err != nil {
		l.WithError(err).Error("map failed")
		return nil, fmt.Errorf("krbd error: %v", err)
	}

	// now go find our ID
	if err := dev.Find(); err != nil {
		l.WithError(err).Error("mapped device was not found")
		return nil, fmt.Errorf("could not find device ID: %v", err)
	}
	rbd.DeviceID = dev.ID
	rbd.DeviceFile = dev.DevPath()
	rbd.ID = r.next
	rbd.Refs = 1
	r.next++
	r.rbds[rbd.ID] = rbd
	l.Info("successfully mapped")
	return rbd, err
}

func (r *RbdsType) Get(id models.ID) (m types.Model, err error) {
	l := r.log.WithFields(logrus.Fields{
		"operation": "get",
		"id":        id,
	})
	r.mutex.Lock()
	defer r.mutex.Unlock()
	var ok bool
	if m, ok = r.rbds[id]; ok {
		l.Trace("found")
		return
	}
	l.Trace("not found")
	return nil, types.ERRNOTFOUND
}

func (r *RbdsType) Detach(id models.ID) (m types.AttachModel, err error) {
	l := r.log.WithFields(logrus.Fields{
		"operation": "unmap",
		"id":        id,
	})
	r.mutex.Lock()
	defer r.mutex.Unlock()

	var rbd *RbdModel
	var ok bool

	if rbd, ok = r.rbds[id]; !ok {
		l.Debug("not found")
		return nil, types.ERRNOTFOUND
	}
	l = l.WithFields(logrus.Fields{
		"image":    *rbd.Image,
		"pool":     *rbd.Pool,
		"name":     rbd.Options.Name,
		"monitors": rbd.Monitors,
	})

	// should we be able to force this?
	if rbd.Refs > 0 {
		l.WithField("refs", rbd.Refs).Debug("nonzero refcount")
		return nil, fmt.Errorf("device %d is in use, cannot unmap", id)
	}

	wc, err := krbd.RBDBusRemoveWriter()
	if err != nil {
		l.WithError(err).Error("couldn't get remove writer")
		return nil, err
	}
	defer wc.Close()

	i := krbd.Image{
		DevID: int(rbd.DeviceID),
		Options: &krbd.Options{
			Force: rbd.Options.Force,
		},
	}

	if err := i.Unmap(wc); err != nil {
		l.WithError(err).Error("unmap failed")
		return nil, fmt.Errorf("krbd error: %v", err)
	}
	// remove from our map
	delete(r.rbds, id)
	l.Info("successfully unmapped")

	return rbd, nil
}

// add/subtract from ref counter
// silently fails if id doesn't exist
func (r *RbdsType) RefAdd(id models.ID, n int64) {
	l := r.log.WithField("operation", "refadd")
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if rbd, ok := r.rbds[id]; ok {
		l.WithField("n", n).Trace("added")
		rbd.Refs += n
	} else {
		l.Debug("no such rbd")
	}
}

// Collect will run garbage collection on any RBDs with ref == 0
func (r *RbdsType) Collect() {
	l := r.log.WithField("operation", "collect")
	list := []models.ID{}
	r.mutex.Lock()
	for _, rbd := range r.rbds {
		if rbd.Refs <= 0 {
			// let's collect
			list = append(list, rbd.ID)
		}
	}
	r.mutex.Unlock()
	if len(list) > 0 {
		l.WithField("collectIDs", list).Debug("collecting")
		for _, id := range list {
			r.Detach(id)
		}
	}
}
