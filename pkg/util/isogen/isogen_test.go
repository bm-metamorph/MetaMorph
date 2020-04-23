package isogen

import (
	"testing"
	"github.com/stretchr/testify/assert"
    "net/url"
    "bitbucket.com/metamorph/pkg/db/models/node"
 ///   "errors"
)

var tests_urldownload = map[string]struct {
    urlpath  string
    downldpath string
    err    error
}{
     "successfull download": { 
	   urlpath :"http://dl-cdn.alpinelinux.org/alpine/v3.11/releases/x86_64/alpine-standard-3.11.5-x86_64.iso",
	   downldpath : "/tmp/alpine.iso",
	   err:    nil, },
    "invalid download": {
        urlpath:  "http1://test.com",
        downldpath: "/tmp/alpine2.iso", 
        err:   &url.Error{},
    },
}

var tests_extractcopy = map[string]struct {
    destpath  string
    downldpath string
    err    error
}{
     "successfull download": { 
	   downldpath : "/tmp/alpine.iso",
     destpath : "/tmp/isodest/alpine.iso",
	   err:    nil, },
    "invalid download": {
        destpath: "/tmp/isodest/alpine2.iso",
        downldpath: "/tmp/alpine2.iso", 
        err:   &url.Error{},
    },
}

func TestDownloadURL(t *testing.T){
	
    for testName, test := range tests_urldownload {
        t.Logf("Running test case %s", testName)
        err := DownloadUrl(test.downldpath, test.urlpath) 
        assert.IsType(t, test.err, err)
	}
}

func TestExtractAndCopyISO(t *testing.T){
    for testName, test := range tests_extractcopy {
        t.Logf("Running test case %s", testName)
        err := ExtractAndCopyISO(test.downldpath, test.destpath)
        t.Logf("%v",err)
        assert.IsType(t, test.err, err)
    }
}

func TestValidateChecksum(t *testing.T){
    checksumURL := "http://dl-cdn.alpinelinux.org/alpine/v3.11/releases/x86_64/alpine-standard-3.11.5-x86_64.iso.sha256"
    err := ValidateChecksum(checksumURL, "/tmp/alpine.iso" )
    t.Logf("%v", err)
}

func TestPrepareISO(t *testing.T){
	bmhnode := &BMHNode{node.CreateTestNode()}
    err := bmhnode.PrepareISO()
    t.Logf("%v", err)
    assert.IsType(t, err,nil)
}