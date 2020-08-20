package util

import (

	"testing"
	"net/url"
	"github.com/stretchr/testify/assert"
	"github.com/manojkva/metamorph-plugin/pkg/config"
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

func TestDownloadURL(t *testing.T){
	config.SetLoggerConfig("logger.pluginpath")

    for testName, test := range tests_urldownload {
        t.Logf("Running test case %s", testName)
        err := DownloadUrl(test.downldpath, test.urlpath) 
        assert.IsType(t, test.err, err)
	}
}


