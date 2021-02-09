package api

var Rbds RbdsType
var MountsRbd MountsRBDType
var MountsOverlay MountsOverlayType
var Containers ContainersType

const mountDir = "/var/run/imageapi/mounts"
const logDir = "/var/run/imageapi/logs"

func init() {
	Rbds = RbdsType{}
	Rbds.Init()
	MountsRbd = MountsRBDType{}
	MountsRbd.Init()
	MountsOverlay = MountsOverlayType{}
	MountsOverlay.Init()
	Containers = ContainersType{}
	Containers.Init()
}
