package api

var Rbds RbdsType
var MountsRbd MountsRBDType
var MountsOverlay MountsOverlayType

const mountDir = "/var/run/rbd-server/mounts"

func init() {
	Rbds = RbdsType{}
	Rbds.Init()
	MountsRbd = MountsRBDType{}
	MountsRbd.Init()
	MountsOverlay = MountsOverlayType{}
	MountsOverlay.Init()
}
