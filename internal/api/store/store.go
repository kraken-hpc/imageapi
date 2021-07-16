package store

import (
	"sync"

	"github.com/kraken-hpc/imageapi/models"
)

// The ObjectStore centralize EndpointObject storage
type ObjectStore struct {
	objs  map[models.ID]models.EndpointObject
	next  models.ID
	mutex *sync.Mutex
}

func (s *ObjectStore) Init() {
	s.next = 1
	s.objs = make(map[models.ID]models.EndpointObject)
	s.mutex = &sync.Mutex{}
}

func (s *ObjectStore) Register(o models.EndpointObject) models.EndpointObject {
	if o.GetID() != 0 {
		// refuse to register an object with non-zero id
		return nil
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	o.SetID(s.next)
	s.objs[s.next] = o
	s.next++
	return o
}

func (s *ObjectStore) Unregister(o models.EndpointObject) {
	if o.GetID() == 0 {
		return
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.objs, o.GetID())
}

func (s *ObjectStore) Update(o models.EndpointObject) {
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

func (s *ObjectStore) Get(id models.ID) models.EndpointObject {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if obj, ok := s.objs[id]; ok {
		return obj
	}
	return nil
}

func (s *ObjectStore) List() []models.EndpointObject {
	l := []models.EndpointObject{}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, o := range s.objs {
		l = append(l, o)
	}
	return l
}

func (s *ObjectStore) ListType(t models.EndpointObjectType) []models.EndpointObject {
	l := []models.EndpointObject{}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, o := range s.objs {
		if o.EndpointObjectType() == t {
			l = append(l, o)
		}
	}
	return l
}
