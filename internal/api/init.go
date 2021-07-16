// initialization for internal data structures & processes
package api

import (
	"github.com/kraken-hpc/imageapi/internal/api/types"
	"github.com/sirupsen/logrus"
)

var Rbds *RbdsType
var MountsRbd *MountsRBDType
var MountsOverlay *MountsOverlayType
var Containers *ContainersType
var Log *logrus.Logger

var MountDir string = "/var/run/imageapi/mounts"
var LogDir string = "/var/run/imageapi/logs"

var Collections = []types.Collectable{}

func GarbageCollect() {
	MountsOverlay.Collect()
	MountsRbd.Collect()
	Rbds.Collect()
}
