package api

import (
	"fmt"

	"github.com/kraken-hpc/imageapi/models"
	"github.com/sirupsen/logrus"
)

var AttachDrivers = map[string]AttachDriver{}

type Attach models.Attach

// Make sure Attach is an EndpointObject
var _ EndpointObject = (*Attach)(nil)

func (a *Attach) GetID() models.ID                       { return a.ID }
func (a *Attach) SetID(id models.ID)                     { a.ID = id }
func (a *Attach) GetRefs() int64                         { return a.Refs }
func (a *Attach) RefAdd(i int64)                         { a.Refs += i }
func (a *Attach) EndpointObjectType() EndpointObjectType { return EndpointObjectAttach }

type Attachments struct {
	log *logrus.Entry
}

// Init initializes the attachment subsystem
func (a *Attachments) Init(log *logrus.Entry) {
	a.log = log
	a.log.Info("initializing attachment drivers")
	for name, drv := range AttachDrivers {
		a.log.Debugf("initializing driver: %s", name)
		drv.Init(a.log.WithField("driver", name))
	}
	a.log.Info("attachment subsystem initialized")
}

// List lists all attachments
func (a *Attachments) List() (ret []*Attach) {
	ret = []*Attach{}
	for _, o := range API.Store.ListType(EndpointObjectAttach) {
		ret = append(ret, o.(*Attach))
	}
	return
}

// Get gets an attachment by ID
func (a *Attachments) Get(id models.ID) *Attach {
	if eo := API.Store.Get(id); eo != nil {
		if ret, ok := eo.(*Attach); ok {
			return ret
		}
	}
	return nil
}

// Attach attaches an attachment
func (a *Attachments) Attach(at *Attach) (ret *Attach, err error) {
	if at.ID != 0 {
		a.log.Errorf("requested an attachment with non-zero attachment ID")
		return nil, fmt.Errorf("requested an attachment with non-zero attachment ID")
	}
	if drv, ok := AttachDrivers[at.Kind]; ok {
		ret, err = drv.Attach(at)
		if err == nil {
			ret = API.Store.Register(ret).(*Attach)
		}
		return
	}
	return nil, fmt.Errorf("no driver found for attachment kind %s", at.Kind)
}

// GetOrAttach gets an attachment if it already exists, if it does not, it attempts to attach it
func (a *Attachments) GetOrAttach(at *Attach) (ret *Attach, err error) {
	if at.ID != 0 {
		ga := a.Get(at.ID)
		if ga != nil {
			return ga, nil
		}
		return nil, ERRNOTFOUND
	}
	return a.Attach(at)
}

// Detach detaches an attachment
func (a *Attachments) Detach(at *Attach, force bool) (ret *Attach, err error) {
	l := a.log.WithField("operation", "detach")
	if at.ID < 1 {
		l.Trace("detach called with ID 0")
		return nil, ERRNOTFOUND
	}
	eo := API.Store.Get(at.ID)
	if eo == nil {
		l.Tracef("detach called on non-existent attach ID: %d", at.ID)
		return nil, ERRNOTFOUND
	}
	defer func() {
		API.Store.RefAdd(eo.GetID(), -1)
	}()
	var ok bool
	if at, ok = eo.(*Attach); !ok {
		l.Trace("detach called on non-attach object")
		return nil, ERRNOTFOUND
	}
	l = l.WithField("id", at.ID)
	if at.Refs > 1 && !force { // we hold 1 from the Get above
		l.Debug("detach called on an attachment that is in use")
		return nil, ERRBUSY
	}
	if drv, ok := AttachDrivers[at.Kind]; ok {
		ret, err = drv.Detach(at)
		if err == nil {
			API.Store.Unregister(ret)
		}
		return ret, err
	}
	return nil, ERRNODRV
}
