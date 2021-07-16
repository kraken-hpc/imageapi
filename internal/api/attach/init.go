// attach interface and init
package attach

import (
	"github.com/kraken-hpc/imageapi/internal/api/types"
	"github.com/sirupsen/logrus"
)

var Rbds *RbdsType
var Log *logrus.Entry

var endpoints map[string]types.Endpoint

// Initialize attachment types
func Init(log *logrus.Logger) {
	Log = log.WithField("category", "attach")
}
