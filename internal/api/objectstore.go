package api

import (
	"sync"
	"time"

	"github.com/kraken-hpc/imageapi/models"
	"github.com/sirupsen/logrus"
)

// The ObjectStore centralize EndpointObject storage
type ObjectStore struct {
	objs  map[models.ID]EndpointObject
	next  models.ID
	mutex *sync.Mutex
}

func (s *ObjectStore) Init() {
	s.next = 1
	s.objs = make(map[models.ID]EndpointObject)
	s.mutex = &sync.Mutex{}
	// start the garbage collector
	go func() {
		for {
			time.Sleep(API.CollectInterval)
			s.mutex.Lock()
			for _, o := range s.objs {
				if o.GetID() == 0 {
					go s.collect(o)
				}
			}
			s.mutex.Unlock()
		}
	}()
}

func (s *ObjectStore) Register(o EndpointObject) EndpointObject {
	if o.GetID() != 0 {
		// refuse to register an object with non-zero id
		return nil
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	o.SetID(s.next)
	s.objs[s.next] = o
	s.refAdd(o.GetID(), 1)
	s.next++
	return o
}

func (s *ObjectStore) Unregister(o EndpointObject) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.unregister(o)
}

// internal, non-locking version
func (s *ObjectStore) unregister(o EndpointObject) {
	if o.GetID() == 0 {
		return
	}
	delete(s.objs, o.GetID())
}

func (s *ObjectStore) Update(o EndpointObject) {
	if o.GetID() == 0 {
		return
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if old, ok := s.objs[o.GetID()]; ok {
		if old.EndpointObjectType() != o.EndpointObjectType() {
			// don't update unlike types
			return
		}
		s.objs[o.GetID()] = o
	}
	// we don't update things that don't exist
}

func (s *ObjectStore) Get(id models.ID) EndpointObject {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if obj, ok := s.objs[id]; ok {
		s.refAdd(id, 1)
		return obj
	}
	return nil
}

func (s *ObjectStore) List() []EndpointObject {
	l := []EndpointObject{}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, o := range s.objs {
		l = append(l, o)
	}
	return l
}

func (s *ObjectStore) ListType(t EndpointObjectType) []EndpointObject {
	l := []EndpointObject{}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, o := range s.objs {
		if o.EndpointObjectType() == t {
			l = append(l, o)
		}
	}
	return l
}

func (s *ObjectStore) RefAdd(id models.ID, i int64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.refAdd(id, i)
}

// non-locking refadd for internal use
func (s *ObjectStore) refAdd(id models.ID, i int64) {
	if obj, ok := s.objs[id]; ok {
		obj.RefAdd(i)
		if obj.GetRefs() == 0 {
			go s.collect(obj)
		}
	}
}

// attempt to collect an object.  This is unfortunately hardcoded.
// in the future it may make sense to abstract this out
// collect should only be called when a lock is already heald
func (s *ObjectStore) collect(eo EndpointObject) {
	l := API.Log.WithFields(logrus.Fields{
		"subsys":    "store",
		"operation": "collect",
		"id":        eo.GetID(),
	})
	switch eo.EndpointObjectType() {
	case EndpointObjectAttach:
		if _, err := API.Attachments.Detach(eo.(*Attach), false); err == nil {
			s.Unregister(eo)
			l.Trace("successfully collected attach")
			return
		}
	case EndpointObjectMount:
		if _, err := API.Mounts.Unmount(eo.(*Mount), false); err == nil {
			s.Unregister(eo)
			l.Trace("successfully collected mount")
			return
		}
	case EndpointObjectContainer:
		// we don't currently garbage collect these
	}
}
