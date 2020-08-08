package redfish

import (
	config "github.com/manojkva/metamorph-plugin/pkg/config"
	"github.com/bm-metamorph/MetaMorph/pkg/db/models/node"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUpgradeFirmware(t *testing.T) {
	config.SetLoggerConfig("logger.apipath")
	bmhnode := &BMHNode{node.CreateTestNode()}
	res := bmhnode.UpgradeFirmware()
	assert.Equal(t, res, nil)
}
