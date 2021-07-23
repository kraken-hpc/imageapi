package api

// API operations on rbd maps

import (
	"github.com/bensallen/rbd/pkg/krbd"
	"github.com/kraken-hpc/imageapi/models"
	"github.com/sirupsen/logrus"
)

func init() {
	AttachDrivers[models.AttachKindRbd] = &AttachDriverRbd{}
}

type AttachDriverRbd struct {
	log *logrus.Entry
}

func (a *AttachDriverRbd) Init(log *logrus.Entry) {
	a.log = log
	a.log.Trace("initialized")
}

func (a *AttachDriverRbd) Attach(att *Attach) (ret *Attach, err error) {
	// sanity check
	l := a.log.WithField("operation", "attach")
	if att.Rbd == nil {
		l.Trace("attempted to attach rbd with no rbd definition")
		return nil, ERRINVALDAT
	}

	if len(att.Rbd.Monitors) == 0 || *att.Rbd.Pool == "" || *att.Rbd.Image == "" || att.Rbd.Options.Name == "" || att.Rbd.Options.Secret == "" {
		a.log.Debug("incorrect options")
		return nil, ERRINVALDAT
	}
	l = l.WithFields(logrus.Fields{
		"image":    *att.Rbd.Image,
		"pool":     *att.Rbd.Pool,
		"name":     att.Rbd.Options.Name,
		"monitors": att.Rbd.Monitors,
	})
	w, err := krbd.RBDBusAddWriter()
	if err != nil {
		l.WithError(err).Error("failed to get krbd bus writer")
		return nil, ERRSRV
	}
	defer w.Close()

	// We allow this because we get free IPv4 format checking
	mons := []string{}
	for _, m := range att.Rbd.Monitors {
		mons = append(mons, m.String())
	}

	i := krbd.Image{
		Monitors: mons,
		Pool:     *att.Rbd.Pool,
		Image:    *att.Rbd.Image,
		Snapshot: att.Rbd.Snapshot,
		Options: &krbd.Options{
			ReadOnly:  att.Rbd.Options.Ro,
			Name:      att.Rbd.Options.Name,
			Secret:    att.Rbd.Options.Secret,
			Namespace: att.Rbd.Options.Namespace,
		},
	}
	// make sure ID doesn't already exist
	dev := krbd.Device{Image: i.Image, Pool: i.Pool, Namespace: i.Options.Namespace, Snapshot: i.Snapshot}

	if err := dev.Find(); err == nil {
		l.Debug("tried to map device that already exists")
		return nil, ERRBUSY
	}
	// map the rbd
	if err := i.Map(w); err != nil {
		l.WithError(err).Error("map failed")
		return nil, ERRFAIL
	}

	// now go find our ID
	if err := dev.Find(); err != nil {
		l.WithError(err).Error("mapped device was not found")
		return nil, ERRSRV
	}

	att.Rbd.DeviceID = dev.ID
	att.DeviceFile = dev.DevPath()
	l.Info("successfully mapped")
	return att, err
}

func (a *AttachDriverRbd) Detach(att *Attach) (ret *Attach, err error) {
	l := a.log.WithFields(logrus.Fields{
		"operation": "unmap",
		"id":        att.ID,
		"rbd_id":    att.Rbd.DeviceID,
		"image":     *att.Rbd.Image,
		"pool":      *att.Rbd.Pool,
		"name":      att.Rbd.Options.Name,
		"monitors":  att.Rbd.Monitors,
	})

	wc, err := krbd.RBDBusRemoveWriter()
	if err != nil {
		l.WithError(err).Error("couldn't get remove writer")
		return nil, ERRSRV
	}
	defer wc.Close()

	i := krbd.Image{
		DevID: int(att.Rbd.DeviceID),
		Options: &krbd.Options{
			Force: att.Rbd.Options.Force,
		},
	}

	if err := i.Unmap(wc); err != nil {
		l.WithError(err).Error("unmap failed")
		return nil, ERRFAIL
	}
	// remove from our map
	l.Info("successfully unmapped")

	return att, nil
}