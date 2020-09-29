package main

import (
    "fmt"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestNodes(t *testing.T) {
    // Inject the StartServer method into a test server
    ts := httptest.NewServer(StartServer())
    defer ts.Close()

    // Make a request to your server with the {base url}/nodes
    resp, err := http.Get(fmt.Sprintf("%s/nodes", ts.URL))
    require.NoError(t, err, "Error querying nodes")
    assert.Equal(t, http.StatusOK, resp.StatusCode)
    // Do rest of your validations here
}
