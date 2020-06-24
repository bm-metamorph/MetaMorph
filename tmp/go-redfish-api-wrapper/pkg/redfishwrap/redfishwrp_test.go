package redfishwrap


import (
	"testing"
	"github.com/stretchr/testify/assert"
)
func TestCheckStatusCodeforGet(t *testing.T){
	var statustests =[]struct {
		n  int
		expected bool
	}{
		{200, true},
		{204, true},
		{202, true},
		{500, false},
		{400, false},

	}

	for _, tt := range statustests{
		res := checkStatusCodeforGet(tt.n)
		assert.Equal(t, res, tt.expected)
	}

}