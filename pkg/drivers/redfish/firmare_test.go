package redfish

import (
	config "github.com/bm-metamorph/Metamorph/pkg/config"
	"github.com/bm-metamorph/Metamorph/pkg/db/models/node"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUpgradeFirmware(t *testing.T) {
	config.SetLoggerConfig("logger.apipath")
	bmhnode := &BMHNode{node.CreateTestNode()}
	res := bmhnode.UpgradeFirmwareList()
	assert.Equal(t, res, true)
}
