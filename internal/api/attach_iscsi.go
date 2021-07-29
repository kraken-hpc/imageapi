package api

// API operations on rbd maps

import (
	"fmt"
	"net"
	"sync"

	"github.com/kraken-hpc/imageapi/models"
	"github.com/sirupsen/logrus"
	"github.com/u-root/iscsinl"
)

func init() {
	AttachDrivers[models.AttachKindIscsi] = &AttachDriverIscsi{}
}

type AttachDriverIscsi struct {
	log      *logrus.Entry
	sessions map[string]*iscsinl.IscsiTargetSession // devicefile -> session
	mutex    *sync.Mutex
}

func (a *AttachDriverIscsi) Init(log *logrus.Entry) {
	a.log = log
	a.log.Trace("initialized")
	a.sessions = map[string]*iscsinl.IscsiTargetSession{}
	a.mutex = &sync.Mutex{}
}

func (a *AttachDriverIscsi) Attach(att *Attach) (ret *Attach, err error) {
	// sanity check
	l := a.log.WithField("operation", "attach")
	if att.Iscsi == nil {
		l.Trace("attempted to attach iscsi with no iscsi definition")
		return nil, ErrInvalDat
	}
	// we can't really trust go-swagger to handle defaults, so we'll do that explictly
	// we do trust it for required properties though
	defaultMaxCommands := int64(128)
	defaultQueueDepth := int64(16)
	defaultScheduler := "mq-deadline"
	defaultPort := int64(3260)

	if att.Iscsi.MaxComands == nil {
		att.Iscsi.MaxComands = &defaultMaxCommands
	}
	if att.Iscsi.QueueDepth == nil {
		att.Iscsi.QueueDepth = &defaultQueueDepth
	}
	if att.Iscsi.Scheduler == nil {
		att.Iscsi.Scheduler = &defaultScheduler
	}
	if att.Iscsi.Port == nil {
		att.Iscsi.Port = &defaultPort
	}

	l = l.WithFields(logrus.Fields{
		"target":    *att.Iscsi.Target,
		"lun":       att.Iscsi.Lun,
		"initiator": *att.Iscsi.Initiator,
		"host":      *att.Iscsi.Host,
		"port":      *att.Iscsi.Port,
	})

	ips, err := net.LookupIP(*att.Iscsi.Host)
	if err != nil {
		l.WithError(err).Debug("host resolution failed")
		return nil, ErrInvalDat
	}
	ip := ips[0]
	opts := []iscsinl.Option{
		iscsinl.WithInitiator(*att.Iscsi.Initiator),
		iscsinl.WithTarget(fmt.Sprintf("%s:%d", ip.String(), *att.Iscsi.Port), *att.Iscsi.Target),
		iscsinl.WithCmdsMax(uint16(*att.Iscsi.MaxComands)),
		iscsinl.WithQueueDepth(uint16(*att.Iscsi.QueueDepth)),
		iscsinl.WithScheduler(*att.Iscsi.Scheduler),
	}
	// we can't use MountIscsi because it never tells us the session ID
	netlink, err := iscsinl.ConnectNetlink()
	if err != nil {
		l.WithError(err).Error("failed to connect to iscsi netlink socket")
		return nil, ErrSrv
	}
	session := iscsinl.NewSession(netlink, opts...)
	if err = session.Connect(); err != nil {
		l.WithError(err).Debug("iscsi connection failed")
		return nil, ErrFail
	}
	defer func() {
		if err != nil {
			session.TearDown()
		}
	}()
	if err := session.Login(); err != nil {
		l.WithError(err).Debug("iscsi login failed")
		return nil, ErrFail
	}
	if err := session.SetParams(); err != nil {
		l.WithError(err).Debug("iscsi setparams failed")
		return nil, ErrFail
	}
	if err := session.Start(); err != nil {
		l.WithError(err).Debug("iscsi start failed")
		return nil, ErrFail
	}
	devnames, err := session.ConfigureBlockDevs()
	if err != nil {
		l.WithError(err).Debug("iscsi failed to configure block devices")
		return nil, ErrFail
	}
	if att.Iscsi.Lun > int64(len(devnames)-1) {
		l.Debug("iscsi lun was not found")
	}
	att.DeviceFile = fmt.Sprintf("/dev/%s", devnames[att.Iscsi.Lun])
	if err := iscsinl.ReReadPartitionTable(att.DeviceFile); err != nil {
		l.WithField("devicefile", att.DeviceFile).WithError(err).Debug("failed to reread partition tables")
		return nil, ErrFail
	}
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.sessions[att.DeviceFile] = session
	return att, nil
}

func (a *AttachDriverIscsi) Detach(att *Attach) (ret *Attach, err error) {
	l := a.log.WithFields(logrus.Fields{
		"operation":  "detach",
		"target":     *att.Iscsi.Target,
		"lun":        att.Iscsi.Lun,
		"initiator":  *att.Iscsi.Initiator,
		"host":       *att.Iscsi.Host,
		"port":       *att.Iscsi.Port,
		"devicefile": att.DeviceFile,
	})
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if s, ok := a.sessions[att.DeviceFile]; ok {
		err = s.TearDown()
		if err != nil {
			l.WithError(err).Debug("detach failed")
			return ret, ErrFail
		}
		delete(a.sessions, att.DeviceFile)
		return att, nil
	}
	return nil, ErrNotFound
}
