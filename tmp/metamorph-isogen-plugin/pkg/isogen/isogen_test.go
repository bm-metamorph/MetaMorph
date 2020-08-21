package isogen

import (
	"net/url"
	"testing"

	"github.com/bm-metamorph/MetaMorph/pkg/db/models/node"
	config "github.com/manojkva/metamorph-plugin/pkg/config"
	"github.com/stretchr/testify/assert"
	///   "errors"
)

func init() {
	config.SetLoggerConfig("logger.plugins.isogenpluginpath")

}

var tests_extractcopy = map[string]struct {
	destpath   string
	downldpath string
	err        error
}{
	"successfull download": {
		downldpath: "/tmp/alpine.iso",
		destpath:   "/tmp/isodest/alpine.iso",
		err:        nil},
	"invalid download": {
		destpath:   "/tmp/isodest/alpine2.iso",
		downldpath: "/tmp/alpine2.iso",
		err:        &url.Error{},
	},
}

func TestExtractAndCopyISO(t *testing.T) {
	for testName, test := range tests_extractcopy {
		t.Logf("Running test case %s", testName)
		err := ExtractAndCopyISO(test.downldpath, test.destpath)
		t.Logf("%v", err)
		assert.IsType(t, test.err, err)
	}
}

func TestValidateChecksum(t *testing.T) {
	checksumURL := "http://dl-cdn.alpinelinux.org/alpine/v3.11/releases/x86_64/alpine-standard-3.11.5-x86_64.iso.sha256"
	err := ValidateChecksum(checksumURL, "/tmp/alpine.iso")
	t.Logf("%v", err)
}

func TestCreateISO(t *testing.T) {
	bmhnode := &BMHNode{node.CreateTestNode()}
	err := bmhnode.CreateISO()
	t.Logf("%v", err)
	assert.IsType(t, err, nil)
}
