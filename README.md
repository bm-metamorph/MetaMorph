# MetaMorph

MetaMorph is a tool introduced to provision baremetal nodes in the kubernetes native way. It is fully compliant with Kubernetes Cluster API, its a Baremetal Provisioner.   MetaMorph uses native redfish APIs to provision the baremetal nodes thus eliminating complex traditional pre-requisties like DHCP, TFTP, PXE booting etc. ISO used to provision the OS will be mounted from an HTTP share using `VirtualMedia` feature of Redfish

## Features

1. **Minimum Pre-requisties/Dependencies**  The only Pre-requisties Metamorph has is the Redfish Protocal support on the node to be provisioned.
2. **Edge Node Deployment** Since Metamorph eliminates the complex pre-requisties like DHCP, TFTP, PXE booting etc, Its very easy and reliable to deploy edge nodes.
3. **Vendor Independent** Servers manufactored by any vendor can be deployed (provided it supports Redfish protocol)
4. **Boot Actions** Boot Actions are jobs that will be executed on the first boot of the deployed node. It can be used to deploy any kind of software on the target node. Boot Actions can be written in any languages.
5. **Plugin Support**  Plugins can do variety of things. It can extend default features/functionlity, Add support for old hardware that doesn't support Redfish protocol (eg: HP iLO4 RAID Config), Integrate othe tools/services to MetaMorph \**coming soon*.
	


## Setup Dev env 

1. `git clone https://github.com/bm-metamorph/MetaMorph.git`


Setup ENV variables

```
export METAMORPH_CONFIGPATH=<path of config.yaml location>
export REDFISH_SLEEPTIME_SECS=10 //time duration in between subsequent redfish API calls. Default = 120 secs
export METAMORPH_POWERCHANGE_TIMEOUT=300 //To handle Nodes that are powered off at the start of RAID Creation.
export METAMORPH_LOG_LEVEL= 1 // default is DEBUG, 1 = INFO 2 = WARN 3 = ERROR 
export REDFISH_JOBCHECKTIMEOUT_MTS= 10 // default is 10 minutes. Job related to firmware updates takes more time. 
cd MetaMorph
```

2. To Run the Controller
    `go run main.go controller`

3. To Run the API
    `go run main.go api`

## Resources

* [Usage Documentation](https://metamorph.readthedocs.io/en/latest/usageGuide)
* [API References](https://metamorph.readthedocs.io/en/latest/references)



	
