package api

import (
	"fmt"
	"sync"

	"github.com/bensallen/rbd/models"
)

type containerStateChange struct {
	id    int64
	state models.ContainerState
}

type ContainersType struct {
	next       int64
	ctns       map[int64]*models.Container
	mutex      *sync.Mutex
	sWatchChan chan containerStateChange
}

func (c *ContainersType) Init() {
	c.next = 0
	c.ctns = make(map[int64]*models.Container)
	c.mutex = &sync.Mutex{}
	c.sWatchChan = make(chan containerStateChange)
}

func (c *ContainersType) List() (r []*models.Container) {
	for _, ctn := range c.ctns {
		r = append(r, ctn)
	}
	return
}

func (c *ContainersType) Get(id int64) (*models.Container, error) {
	if r, ok := c.ctns[id]; ok {
		return r, nil
	}
	return nil, fmt.Errorf("no container by id %d", id)
}

func (c *ContainersType) Create(ctn *models.Container) (r *models.Container, err error) {
	return
}

func (c *ContainersType) SetState(id int64, state string) (err error) {
	return
}

func (c *ContainersType) Delete(id int64) (err error) {
	return
}
