package api

import (
	"errors"
	"time"
)

var Rbds RbdsType
var MountsRbd MountsRBDType
var MountsOverlay MountsOverlayType
var Containers ContainersType

const mountDir = "/var/run/imageapi/mounts"
const logDir = "/var/run/imageapi/logs"
const collectTime = time.Second * 2

var ERRNOTFOUND = errors.New("not found")

func garbageCollect() {
	for {
		time.Sleep(collectTime)
		MountsOverlay.Collect()
		MountsRbd.Collect()
		Rbds.Collect()
	}
}

func init() {
	Rbds = RbdsType{}
	Rbds.Init()
	MountsRbd = MountsRBDType{}
	MountsRbd.Init()
	MountsOverlay = MountsOverlayType{}
	MountsOverlay.Init()
	Containers = ContainersType{}
	Containers.Init()
	go garbageCollect()
}
