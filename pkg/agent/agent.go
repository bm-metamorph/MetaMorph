package main

import(
	"fmt"
	//"time"
	"bitbucket.com/metamorph/proto"
	"bitbucket.com/metamorph/pkg/config"
	"bitbucket.com/metamorph/pkg/db/models/node"
	"google.golang.org/grpc"
	"context"
	"encoding/json"
	"sync"
	"sort"
	"net/http"
	"os"
	"io"
	"os/exec"
	"path/filepath"
	"strings"
)

var logdir string
var tempdir string


func grpcClient() ( proto.AgentServiceClient) {
	conn, err := grpc.Dial( config.Get("agent.cntrl_endpoint").(string), grpc.WithInsecure() )
	if err != nil { panic(err) }
	client := proto.NewAgentServiceClient(conn) 
	return client
}

func prioritizeTasks(bootactions []node.BootAction) map[uint][]node.BootAction {
	schedule  := make(map[uint][]node.BootAction)
	for _,ba := range bootactions {
	 schedule[ba.Priority] =  append(schedule[ba.Priority], ba)
  }
	 return schedule
}

func executeTask( task node.BootAction,wg *sync.WaitGroup ){
	defer wg.Done()
	task.Status = "runninng"
	UpdateTaskStatus(task)
	fmt.Println("executing", task.Name, "with Priority" , task.Priority)
	runCmd(task.Location, task.Control, task.Args)
	//interval := rand.Intn(20)
	//time.Sleep(time.Second * time.Duration(interval))
	task.Status = "done"
	UpdateTaskStatus(task)
	fmt.Println("Executed", task.Name, "with Priority" , task.Priority)
}

func ProcessTasks(tasks []node.BootAction){
	var wg sync.WaitGroup
	prioritizedTasks := prioritizeTasks(tasks)

	keys := []int{}
	for i,_ := range prioritizedTasks {
	  keys = append(keys, int(i))
	}
	sort.Ints(keys)

	for i := range keys {
        for _,ba := range prioritizedTasks[uint(keys[i])] {
          wg.Add(1)
          go executeTask(ba, &wg)
        }
        wg.Wait()
      }
}

func DownloadFile(filepath string, url string) error {

    // Get the data
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // Create the file
    out, err := os.Create(filepath)
    if err != nil {
        return err
    }
    defer out.Close()

    // Write the body to file
    _, err = io.Copy(out, resp.Body)
    return err
}


func runCmd(location string, cntrl string, args string) {

    // get `go` executable path
	goExecutable, _ := exec.LookPath( cntrl)

	_, file := filepath.Split(location)
	filePath := tempdir + "/" + file
	if err := DownloadFile(filePath, location); err != nil {
        panic(err)
	}
    f, _ := os.Create( logdir + "/" + file + ".log")
	defer f.Close()
	
	Args := []string{}
    if args != "" {
		Args = strings.Split(args, " ")
		Args = append([]string{goExecutable}, Args...)
    } else {
		Args = append(Args, goExecutable)
	}
	Args = append(Args, filePath)

    // construct `go version` command
    cmdGoVer := &exec.Cmd {
        Path: goExecutable,
        Args: Args,
        Stdout: f,
        Stderr: f,
	}
	
	fmt.Println( "Executing ", cmdGoVer.String() )

    // run `go version` command
    if err := cmdGoVer.Run(); err != nil {
        fmt.Println( "Error:", err );
    }

}

func UpdateTaskStatus(task node.BootAction){
	
	//Update the task status using GRPC
	ctx := context.Background()
	client := grpcClient()
	input, _ := json.Marshal(task)
	req := &proto.Request{Task: input }
	if response, err := client.UpdateTaskStatus(ctx, req); err == nil {
		
		fmt.Println(response.Result)

	} else {
		fmt.Println(err.Error())
	}


}


func UpdateNodeState(NodeID string, state string){
	ctx := context.Background()
	client := grpcClient()
	req := &proto.Request{NodeID: NodeID, NodeState : state }
	if response, err := client.UpdateNodeState(ctx, req); err == nil {
		fmt.Println(response.Result)

	} else {
		fmt.Println(err.Error())
	}

}


func main() {

	logdir = config.Get("agent.logdir").(string)
	tempdir = config.Get("agent.temp_dir").(string)

	os.Mkdir(logdir, 0744)
	os.Mkdir(tempdir, 0744)

	nodeID := config.Get("agent.node_id").(string)
	ctx := context.Background()
	client := grpcClient()
	req := &proto.Request{NodeID: nodeID }
	bootactions := []node.BootAction{}
	if response, err := client.GetTasks(ctx, req); err == nil {
		//Process the Tasks
		json.Unmarshal(response.Res, &bootactions)
		ProcessTasks(bootactions)
		UpdateNodeState(nodeID, "deployed")
	} else {
		fmt.Println(err.Error())
	}
}


