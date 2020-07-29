package controller

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/bm-metamorph/MetaMorph/pkg/db/models/node"
	"github.com/bm-metamorph/MetaMorph/pkg/logger"
	//	"github.com/bm-metamorph/MetaMorph/pkg/util/isogen"
	"github.com/bm-metamorph/MetaMorph/pkg/plugin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type BMNode struct {
	*node.Node
}

type nodedb interface {
	GetNodes() ([]node.Node, error)
}

type DBHandler struct {
	db nodedb
}

func (db *DBHandler) GetNodes() ([]node.Node, error) {
	return node.GetNodes()
}

const (
	NEW            = "new"
	READY          = "ready"
	READYWAIT      = "readywait"
	SETUPREADY     = "setupready"
	SETUPREADYWAIT = "setupreadywait"
	DEPLOYED       = "deployed"
	DEPLOYING      = "deploying"
	FAILED         = "failed"
	INTRANSITION   = "in-transition"
	USERDATALOADED = "userdataloaded"
)

type NodeStatus struct {
	NodeUUID uuid.UUID
	Status   bool
}

func StartMetamorphFSM(runOnce bool) {
	logger.Log.Info("Starting Metamorph FSM")
	dbHandler := new(DBHandler)
	dbHandler.startFSM(runOnce)

}

func (h *DBHandler) startFSM(runOnce bool) {

	logger.Log.Info("Starting Metamorph FSM")

	var wg sync.WaitGroup

	requestsChan := make(chan BMNode)
	nodeStatusChan := make(chan NodeStatus) //NodeUUID of nodes that did not make successful update to db
	//TODO : ensure sync.Oncce for these two go routines
	//TODO : Waitgroup for goroutines
	go checkFailedNodes(nodeStatusChan, &wg)
	go serviceRequest(requestsChan, nodeStatusChan, &wg)
	wg.Add(2)

	logger.Log.Debug("Number of Go Routines at the Start", zap.Int("Number ", runtime.NumGoroutine()))

	for {

		nodelist, err := node.GetNodes()
		//nodelist, err := h.db.GetNodes()

		if err != nil {
			logger.Log.Debug("No Nodes to process...")

		} else {

			for _, bmnode := range nodelist {
				// What about nodes that are already in transistions.. should there be a transition state.
				fmt.Printf("[%v] - Starting Processing\n", bmnode.Name)
				requestsChan <- BMNode{&bmnode}

			}
		}
		// set the array to nil for the next cycle
		nodelist = nil

		time.Sleep(10 * time.Second) // sleep for 10 ms before the start of the next cycle
		if runOnce == true {         // for testing purpose only.
			break
		}
	}

	logger.Log.Debug("Number of Go Routines Before start of close", zap.Int("Number ", runtime.NumGoroutine()))

	close(requestsChan)
	close(nodeStatusChan)
	wg.Wait()
	logger.Log.Debug("Number of Go Routines at the End", zap.Int("Number ", runtime.NumGoroutine()))

}

func checkFailedNodes(nodeStatusChan chan NodeStatus, wg *sync.WaitGroup) {

	for nodestatus := range nodeStatusChan {

		if nodestatus.Status == false {
			logger.Log.Warn("Failed Node ", zap.String("NodeUUID", nodestatus.NodeUUID.String()))
			//try update of the db
		}
	}
	logger.Log.Info("Closing checkFailedNodes Goroutine")
	wg.Done()
}

func serviceRequest(requestsChan chan BMNode, nodeStatusChan chan<- NodeStatus, wg *sync.WaitGroup) {

	for bmnode := range requestsChan {

		switch bmnode.State {
		case NEW:
			logger.Log.Info("Node Changing States", zap.String("Node Name", bmnode.Name), zap.String("From", NEW), zap.String("To", bmnode.State))
			wg.Add(1)
			go ReadystateHandler(bmnode, nodeStatusChan, wg)
		case READY:
			logger.Log.Info("Node Changing States", zap.String("Node Name", bmnode.Name), zap.String("From", READY), zap.String("To", bmnode.State))
			wg.Add(1)
			go SetupreadyHandler(bmnode, nodeStatusChan, wg)
		case SETUPREADY:
			logger.Log.Info("Node Changing States", zap.String("Node Name", bmnode.Name), zap.String("From", SETUPREADY), zap.String("To", bmnode.State))
			wg.Add(1)
			go DeployedHandler(bmnode, nodeStatusChan, wg)
		default:
			logger.Log.Warn("State not handled...", zap.String("Node Name", bmnode.Name))
			//nodestatus := NodeStatus { NodeUUID: node.NodeUUID, Status: false }
			//nodeStatusChan <- nodestatus
		}

		logger.Log.Debug("Number of Go Routines During request handling", zap.Int("Number ", runtime.NumGoroutine()))

	}
	logger.Log.Info("Closing serviceRequest Goroutine")
	wg.Done()

}

func ReadystateHandler(bmnode BMNode, nodeStatusChan chan<- NodeStatus, wg *sync.WaitGroup) {
	var err error
	var state string
	var nodestatus NodeStatus
	var node_uuid uuid.UUID // to be removed once UUID = Server UUID is planned.
	var nodeuuidStringFromServer string
	var redfishManagerID string 
	var redfishSystemID string
	var redfishVersion string
	var maphwInventory  map[string]string
	var hwInventory  interface{}
	redfishClient := &plugin.BMHNode{bmnode.Node}
	logger.Log.Info("ReadystateHandler()", zap.String("Node Name", bmnode.Name), zap.String("Node UUID", bmnode.NodeUUID.String()))
	//Update the DB Now
	err = node.Update(&node.Node{State: INTRANSITION, NodeUUID: bmnode.NodeUUID})

	if err != nil {
		logger.Log.Error(" Failed to update DB.Setting Node to FAILED State", zap.String("Node Name", bmnode.Name))
		goto End

	}


	// Check if we could extract UUID from the Node using Redfish
	//This call should satisfy the following requirements
	// - Network Connectivity
	// - working credentials
	// - Redfish API availability(though ver is not compared yet)
        err = redfishClient.ReadConfigFile()
	if  err  != nil{
		logger.Log.Error("Failed to Read Plugin info from configuration file. Setting Node to FAILED State", zap.String("NodeName", bmnode.Name))
		goto End
	}
	hwInventory, err = redfishClient.DispenseClientRequest("gethwinventory")
	maphwInventory = hwInventory.(map[string]string)

	if err != nil {
		logger.Log.Error("Failed to retrieve HW Inventory using Redfish Protocol Setting Node to FAILED State", zap.String("IPMIIP", bmnode.Node.IPMIIP), zap.String("IPMIUser", bmnode.IPMIUser))
		goto End
	}

	redfishManagerID = maphwInventory["RedfishManagerID"]
	redfishSystemID = maphwInventory["RedfishSystemID"]
	redfishVersion = maphwInventory["RedfishVersion"]

	if redfishSystemID != "" {
		uuidasIntf, err := redfishClient.DispenseClientRequest("getguuid")
		if err != nil {
			logger.Log.Error("Failed to retrieve Node GUUID using Redfish Protocol Setting Node to FAILED State", zap.String("IPMIIP", bmnode.Node.IPMIIP), zap.String("IPMIUser", bmnode.IPMIUser))
			goto End
		}
		nodeuuidStringFromServer = uuidasIntf.(string)
	}

	node_uuid, err = uuid.Parse(nodeuuidStringFromServer)

	node_uuid = bmnode.NodeUUID // to be removed once UUID = Server UUID is planned.

	state = READYWAIT
	nodestatus = NodeStatus{NodeUUID: bmnode.NodeUUID, Status: true}
	//Update the DB Now
	err = node.Update(&node.Node{State: state, NodeUUID: node_uuid, RedfishManagerID: redfishManagerID, RedfishSystemID: redfishSystemID, RedfishVersion: redfishVersion})

End:
	if err != nil {
		logger.Log.Error("Failed to update Node to READYWAIT state", zap.String("Node Name", bmnode.Name))
		nodestatus = NodeStatus{NodeUUID: node_uuid, Status: false}
		state = FAILED
		node.Update(&node.Node{State: state, NodeUUID: node_uuid})
	}
	nodeStatusChan <- nodestatus
	wg.Done()

}

func SetupreadyHandler(bmnode BMNode, nodeStatusChan chan<- NodeStatus, wg *sync.WaitGroup) {
	logger.Log.Info("SetupreadyHandler()", zap.String("Node Name", bmnode.Name))
	var nodestatus NodeStatus
	var err error
	var state string
	isogenClient := &plugin.BMHNode{bmnode.Node}

	//Update the DB Now
	nodeuuidString := bmnode.Node.NodeUUID.String()
	err = node.Update(&node.Node{State: INTRANSITION, NodeUUID: bmnode.NodeUUID})
	if err != nil {
		goto End
	}

	fmt.Println(nodeuuidString)
	//check for firmware upgrade
	if bmnode.AllowFirmwareUpgrade {
		redfishClient := &plugin.BMHNode{bmnode.Node}
		_, err = redfishClient.DispenseClientRequest("updatefirmware")
		if err != nil {
			logger.Log.Error("Failed to upgrade Firmware", zap.String("Node Name", bmnode.Name))
			goto End
		}
	}


	//Create iSO
	_, err = isogenClient.DispenseClientRequest("createiso")

	if err != nil {
		logger.Log.Error("Failed to create ISO file ", zap.String("Node Name", bmnode.Name))
		goto End
	}


	nodestatus = NodeStatus{NodeUUID: bmnode.NodeUUID, Status: true}
	state = SETUPREADYWAIT

	err = node.Update(&node.Node{State: state, NodeUUID: bmnode.NodeUUID})
End:
	if err != nil {
		logger.Log.Error("Failed to update to SETUPREADYWAIT state", zap.String("Node Name", bmnode.Name))
		nodestatus = NodeStatus{NodeUUID: bmnode.NodeUUID, Status: false}
		state = FAILED
		node.Update(&node.Node{State: state, NodeUUID: bmnode.NodeUUID})
	}
	nodeStatusChan <- nodestatus
	wg.Done()
}
func DeployedHandler(bmnode BMNode, nodeStatusChan chan<- NodeStatus, wg *sync.WaitGroup) {
	logger.Log.Info("DeployedHandler()", zap.String("Node Name", bmnode.Name))
	var nodestatus NodeStatus
	var state string
	redfishClient := &plugin.BMHNode{bmnode.Node}
	//Update the DB Now
	err := node.Update(&node.Node{State: INTRANSITION, NodeUUID: bmnode.NodeUUID})

	if err != nil {
		logger.Log.Error("Failed to update to TRANSITION state", zap.String("Node Name", bmnode.Name))
		goto End
	}

	fmt.Println(bmnode.NodeUUID)

	_, err = redfishClient.DispenseClientRequest("deployiso")

	if err != nil {
		logger.Log.Error("Failed to Deply ISO", zap.String("Node Name", bmnode.Name))
		goto End
	}


	nodestatus = NodeStatus{NodeUUID: bmnode.NodeUUID, Status: true}
	state = DEPLOYING
	err = node.Update(&node.Node{State: state, NodeUUID: bmnode.NodeUUID})
End:
	if err != nil {
		logger.Log.Error("Failed to update to DEPLOYING state", zap.String("Node Name", bmnode.Name))
		nodestatus = NodeStatus{NodeUUID: bmnode.NodeUUID, Status: false}
		state = FAILED
		node.Update(&node.Node{State: state, NodeUUID: bmnode.NodeUUID})
	}
	nodeStatusChan <- nodestatus
	wg.Done()
}
