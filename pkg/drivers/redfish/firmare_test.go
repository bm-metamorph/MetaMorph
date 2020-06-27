package redfish

import (
	config "bitbucket.com/metamorph/pkg/config"
	"bitbucket.com/metamorph/pkg/db/models/node"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUpgradeFirmware(t *testing.T) {
	config.SetLoggerConfig("logger.apipath")
	bmhnode := &BMHNode{node.CreateTestNode()}
	res := bmhnode.UpgradeFirmwareList()
	assert.Equal(t, res, true)
}
