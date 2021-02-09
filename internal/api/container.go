package api

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"sync"
	"time"

	"github.com/jlowellwofford/imageapi/models"
)

type containerStateChange struct {
	id    int64
	state models.ContainerState
}

type container struct {
	ctn    *models.Container
	log    *log.Logger
	mnt    string
	cancel context.CancelFunc
}

type ContainersType struct {
	next  int64
	ctns  map[int64]*container
	mutex *sync.Mutex
}

func (c *ContainersType) Init() {
	c.next = 0
	c.ctns = make(map[int64]*container)
	c.mutex = &sync.Mutex{}
}

func (c *ContainersType) List() (r []*models.Container) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for _, ctn := range c.ctns {
		r = append(r, ctn.ctn)
	}
	return
}

func (c *ContainersType) Get(id int64) (*models.Container, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if r, ok := c.ctns[id]; ok {
		return r.ctn, nil
	}
	return nil, fmt.Errorf("no container by id %d", id)
}

func (c *ContainersType) Create(ctn *models.Container) (r *models.Container, err error) {
	// This creates a container in our list, and activates its initial state
	// find the mount
	n := &container{}
	switch *ctn.Mount.Kind {
	case models.MountKindOverlay:
		mnt, e := MountsOverlay.Get(*ctn.Mount.ID)
		if e != nil {
			return nil, fmt.Errorf("failed to get mount for container: %v", e)
		}
		n.mnt = mnt.Mountpoint
		MountsOverlay.RefAdd(mnt.ID, 1)
	case models.MountKindRbd:
		mnt, e := MountsRbd.Get(*ctn.Mount.ID)
		if e != nil {
			return nil, fmt.Errorf("failed to get mount for container: %v", e)
		}
		n.mnt = mnt.Mountpoint
		MountsRbd.RefAdd(*mnt.ID, 1)
	}
	// ok, we've got a valid mountpoint
	c.mutex.Lock()
	defer c.mutex.Unlock()
	ctn.ID = c.next

	// set up logger
	ctn.Logfile = path.Join(logDir, fmt.Sprintf("%d-%d.log", ctn.ID, time.Now().Unix()))
	f, e := os.Create(ctn.Logfile)
	if e != nil {
		return nil, fmt.Errorf("could not create log file: %v", e)
	}
	n.log = log.New(f, fmt.Sprintf("container(%d): ", ctn.ID), log.Ldate|log.Ltime|log.Lmsgprefix)
	n.log.Printf("container created")

	// handle initial state
	switch ctn.State {
	case models.ContainerStateRunning:
		// run it
	case models.ContainerStateRestarting,
		models.ContainerStatePaused,
		models.ContainerStateExited,
		models.ContainerStateDead:
		return nil, fmt.Errorf("requested invalid initial container state: %s.  valid initial states: [ %s, %s ]", ctn.State, models.ContainerStateCreated, models.ContainerStateRunning)
	case models.ContainerStateCreated:
		fallthrough
	default: // wasn't specified
		ctn.State = models.ContainerStateCreated
	}

	// container is ready to be entered
	c.ctns[ctn.ID] = n
	c.next++

	return ctn, nil
}

func (c *ContainersType) SetState(id int64, state models.ContainerState) (err error) {
	var ctn *container
	var ok bool
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if ctn, ok = c.ctns[id]; !ok {
		return fmt.Errorf("invalid container id: %d", id)
	}
	// handle state request
	switch state {
	case models.ContainerStateRunning:
		if ctn.ctn.State == state {
			return
		}
		// run it
	case models.ContainerStateExited:
		if ctn.ctn.State == state {
			return
		}
		// stop it
	case models.ContainerStatePaused,
		models.ContainerStateRestarting:
		return fmt.Errorf("container state is not yet implemented: %s", state)
	default: // something not valid
		return fmt.Errorf("can't set state to: %s.  valid initial states: [ %s, %s, %s, %s ]", state,
			models.ContainerStateRunning,
			models.ContainerStateExited,
			models.ContainerStateRestarting,
			models.ContainerStatePaused)
	}
	return
}

func (c *ContainersType) Delete(id int64) (err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	var ctn *container
	var ok bool
	if ctn, ok = c.ctns[id]; !ok {
		return fmt.Errorf("invalid container id: %d", id)
	}
	switch ctn.ctn.State {
	//case models.ContainerStatePaused:
	//case models.ContainerStateRestarting:
	case models.ContainerStateRunning:
		// stop it
	}
	ctn.log.Printf("container deleted")
	ctn.log.Writer().(io.WriteCloser).Close()
	delete(c.ctns, id)
	return
}

func (c *ContainersType) run(ctn *models.Container) {

}
