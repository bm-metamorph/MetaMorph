package api

import (
//	"github.com/manojkva/metamorph-plugin/pkg/logger"
//	"github.com/bm-metamorph/MetaMorph/proto"
	"fmt"
//	"github.com/gin-contrib/zap"
//	"github.com/gin-gonic/gin"
//	"go.uber.org/zap"
//	"google.golang.org/grpc"
//	"net/http"
//	"time"
//        "os"
	"io/ioutil"
          "encoding/json"
	  "testing"
	"reflect"
	"strings"
        "time"
        "net"
//	"bytes"
	"github.com/manojkva/metamorph-plugin/pkg/config"
//	"github.com/gin-gonic/gin"
	"net/http"
//	ctrlgRPCServer "github.com/bm-metamorph/MetaMorph/pkg/controller/grpc"
	"github.com/stretchr/testify/require"
        "net/http/httptest"
	"github.com/stretchr/testify/assert"
)

var w http.ResponseWriter
var node_id string
func init(){

	config.SetLoggerConfig("logger.apipath")
	port := "4040"
	ln, err := net.Listen("tcp", ":" + port)
	ln.Close()
        if err != nil {
           fmt.Println("Controller already running on port ",port,":", err)
        } else {
           fmt.Println("running controller on port ", port)
           go run_controller()
	   time.Sleep(5 * time.Second)
	}
}

func TestGrpcClient(t *testing.T) {
	config.SetLoggerConfig("logger.apipath")
	client,conn:= grpcClient()
	fmt.Println(reflect.TypeOf(client))
	fmt.Println(reflect.TypeOf(conn))
        if client==nil {
                 t.Error("got",client, "want", "proto.NodeServiceClient")
		 if conn==nil{
                 t.Error("got",conn , "want", "*grpc.ClientConn")
		}
        }

}

func TestCreateNode(t *testing.T) {
	config.SetLoggerConfig("logger.apipath")
	r := strings.NewReader(" { \"AllowFirwareUpgrade\": false }")
        router := Serve()
        w := httptest.NewRecorder()

        req, _ := http.NewRequest("POST", "/node/", r)
        router.ServeHTTP(w, req)
	type content struct {
                  Node_id   string     `json:"result"`
	}

        var resp content
        i:=w.Body.String()
        json.Unmarshal([]byte(i), &resp)
        fmt.Println(resp.Node_id)

        assert.Equal(t, 200, w.Code)
}

/* func TestGetUUID(t *testing.T){
	config.SetLoggerConfig("logger.apipath")
	r := strings.NewReader(" { \"AllowFirwareUpgrade\": false }")
        router := Serve()
        w := httptest.NewRecorder()

        req, _ := http.NewRequest("POST", "/uuid", r)
        router.ServeHTTP(w, req)

       fmt.Println(w.Body.String())
        assert.Equal(t, 200, w.Code)
}*/

func TestDescribeNode(t *testing.T) {
	config.SetLoggerConfig("logger.apipath")
// Creating a node and getting its node_id

        r := strings.NewReader(" { \"AllowFirwareUpgrade\": false }")
        router := Serve()
        w := httptest.NewRecorder()

        req, _ := http.NewRequest("POST", "/node/", r)
        router.ServeHTTP(w, req)
        type content struct {
                  Node_id   string     `json:"result"`
        }

        var resp content
        i:=w.Body.String()
        json.Unmarshal([]byte(i), &resp)
        fmt.Println(resp.Node_id)

// end of node creation
	ts := httptest.NewServer(Serve())
	defer ts.Close()
	res, err := http.Get(fmt.Sprintf("%s/node/%s", ts.URL,resp.Node_id))
       if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }

    if res.StatusCode != 200 {
        t.Fatalf("Expected status code 200, got %v", res.StatusCode)
    }
}

func TestDeployNode(t *testing.T) {
	config.SetLoggerConfig("logger.apipath")
// Creating a node and getting its node_id

        r := strings.NewReader(" { \"AllowFirwareUpgrade\": false }")
        router := Serve()
        w := httptest.NewRecorder()

        req, _ := http.NewRequest("POST", "/node/", r)
        router.ServeHTTP(w, req)
        type content struct {
                  Node_id   string     `json:"result"`
        }

        var resp content
        i:=w.Body.String()
        json.Unmarshal([]byte(i), &resp)
        fmt.Println(resp.Node_id)

// end of node creation
        w = httptest.NewRecorder()
        req, _ = http.NewRequest("POST", fmt.Sprintf("/node/deploy/%s", resp.Node_id), nil)
        router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

}

func TestDeleteNode(t *testing.T) {
	config.SetLoggerConfig("logger.apipath")
// Creating a node and getting its node_id

        r := strings.NewReader(" { \"AllowFirwareUpgrade\": false }")
        router := Serve()
        w := httptest.NewRecorder()

        req, _ := http.NewRequest("POST", "/node/", r)
        router.ServeHTTP(w, req)
        type content struct {
                  Node_id   string     `json:"result"`
        }

        var resp content
        i:=w.Body.String()
        json.Unmarshal([]byte(i), &resp)
        fmt.Println(resp.Node_id)

// end of node creation
// Deleting the node
       ts := httptest.NewServer(Serve())
        defer ts.Close()

	client := &http.Client{}

    // Create request
    request, err := http.NewRequest("DELETE", fmt.Sprintf("%s/node/%s", ts.URL,resp.Node_id), nil)
    if err != nil {
        fmt.Println(err)
        return
    }

    // Fetch Request
    response, err := client.Do(request)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer response.Body.Close()

    // Read Response Body
    respBody, err := ioutil.ReadAll(response.Body)
    if err != nil {
        fmt.Println(err)
        return
    }

    // Display Results
    fmt.Println("response Status : ", response.Status)
    fmt.Println("response Headers : ", response.Header)
    fmt.Println("response Body : ", string(respBody))
        if response.StatusCode != 200 {
        t.Fatalf("Expected status code 200, got %v", response.StatusCode)
    }


}

func TestUpdateNode(t *testing.T) {
	config.SetLoggerConfig("logger.apipath")
// Creating a node and getting its node_id

        r := strings.NewReader(" { \"AllowFirwareUpgrade\": false }")
        router := Serve()
        w := httptest.NewRecorder()

        req, _ := http.NewRequest("POST", "/node/", r)
        router.ServeHTTP(w, req)
        type content struct {
                  Node_id   string     `json:"result"`
        }

        var resp content
        i:=w.Body.String()
        json.Unmarshal([]byte(i), &resp)
        fmt.Println(resp.Node_id)

// end of node creation
// Updating the node
       ts := httptest.NewServer(Serve())
        defer ts.Close()

	client := &http.Client{}
	id := fmt.Sprintf("%s",resp.Node_id)

    // Create request
        //r = strings.NewReader(" { \"AllowFirwareUpgrade\": true,\"NodeID\": id }")
        //r = strings.NewReader(fmt.Sprintf("{\"AllowFirwareUpgrade\": true,\"NodeID\": %s}",id))
	update_req := fmt.Sprintf("{\"AllowFirwareUpgrade\": true, \"NodeID\": \"%s\"}",id)
	fmt.Println(update_req)
        r = strings.NewReader(update_req)
        router = Serve()
    request, err := http.NewRequest("PUT", fmt.Sprintf("%s/node/%s", ts.URL,resp.Node_id), r)
    if err != nil {
        fmt.Println(err)
        return
    }

    // Fetch Request
    response, err := client.Do(request)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer response.Body.Close()

    // Read Response Body
    respBody, err := ioutil.ReadAll(response.Body)
    if err != nil {
        fmt.Println(err)
        return
    }

    // Display Results
    fmt.Println("response Status : ", response.Status)
    fmt.Println("response Headers : ", response.Header)
    fmt.Println("response Body : ", string(respBody))
        if response.StatusCode != 200 {
        t.Fatalf("Expected status code 200, got %v", response.StatusCode)
    }


}
//func TestGetNodeHWStatus(t *testing.T) {
//	config.SetLoggerConfig("logger.apipath")
//	ts := httptest.NewServer(Serve())
//	defer ts.Close()
//	resp, err := http.Get(fmt.Sprintf("%s/hwstatus/8c3a4c25-b352-4b82-9347-84257a3e5341", ts.URL))
  //     if err != nil {
    //    t.Fatalf("Expected no error, got %v", err)
    //}

    //if resp.StatusCode != 200 {
      //  t.Fatalf("Expected status code 200, got %v", resp.StatusCode)
   // }
//}

func TestListNodes(t *testing.T) {
        config.SetLoggerConfig("logger.apipath")
        ts := httptest.NewServer(Serve())
        defer ts.Close()
	fmt.Println(ts.URL)
        res, err := http.Get(fmt.Sprintf("%s/nodes", ts.URL))
	fmt.Printf("%T",res)
        fmt.Println("####List of nodes####")
	fmt.Println(*res)
	/*
         fmt.Printf("%T",res.Body)
	   type content struct {
                  Node_id   string     `json:"result"`
        }

        var resp content

        i:=res.Body.String()
        json.Unmarshal([]byte(i), &resp)
        fmt.Println(resp)
        */
       if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }

    if res.StatusCode != 200 {
        t.Fatalf("Expected status code 200, got %v", res.StatusCode)
    }

}

func TestServe(t *testing.T) {
	config.SetLoggerConfig("logger.apipath")
	    // Inject the StartServer method into a test server
    ts := httptest.NewServer(Serve())
    defer ts.Close()

    // Make a request to your server with the {base url}/nodes
    resp, err := http.Get(fmt.Sprintf("%s/nodes", ts.URL))
    require.NoError(t, err, "Error querying nodes")
    assert.Equal(t, http.StatusOK, resp.StatusCode)
}


