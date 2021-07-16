package models

type EndpointObjectType uint8

// Allows inspection without reflection
const (
	EndpointObjectAttach    EndpointObjectType = iota
	EndpointObjectMount     EndpointObjectType = iota
	EndpointObjectContainer EndpointObjectType = iota
)

type EndpointObject interface {
	SetID(ID)
	GetID() ID
	GetRefs() int64
	RefAdd(int64)
	EndpointObjectType() EndpointObjectType
}

// Make sure Attach is an EndpointObject
var _ EndpointObject = (*Attach)(nil)

func (a *Attach) GetID() ID                              { return a.ID }
func (a *Attach) SetID(id ID)                            { a.ID = id }
func (a *Attach) GetRefs() int64                         { return a.Refs }
func (a *Attach) RefAdd(i int64)                         { a.Refs += i }
func (a *Attach) EndpointObjectType() EndpointObjectType { return EndpointObjectAttach }

// Make sure Mount is an EndpointObject
var _ EndpointObject = (*Mount)(nil)

func (m *Mount) GetID() ID                              { return m.ID }
func (m *Mount) SetID(id ID)                            { m.ID = id }
func (m *Mount) GetRefs() int64                         { return m.Refs }
func (m *Mount) RefAdd(i int64)                         { m.Refs += i }
func (m *Mount) EndpointObjectType() EndpointObjectType { return EndpointObjectMount }

// Make sure Container is an EndpointObject
var _ EndpointObject = (*Container)(nil)

func (c *Container) GetID() ID                              { return c.ID }
func (c *Container) SetID(id ID)                            { c.ID = id }
func (c *Container) GetRefs() int64                         { return c.Refs }
func (c *Container) RefAdd(i int64)                         { c.Refs += i }
func (c *Container) EndpointObjectType() EndpointObjectType { return EndpointObjectContainer }
