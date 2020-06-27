# MetaMorph

Lifecycle your BareMetal


## Setting up Development Environment

Setup ENV variables

export METAMORPH_CONFIGPATH=<path of config.yaml location>
export REDFISH_SLEEPTIME_SECS=10 //time duration in between subsequent redfish API calls. Default = 120 secs
export METAMORPH_POWERCHANGE_TIMEOUT=300 //To handle Nodes that are powered off at the start of RAID Creation.
export METAMORPH_LOG_LEVEL= 1 // default is DEBUG, 1 = INFO 2 = WARN 3 = ERROR 
export REDFISH_JOBCHECKTIMEOUT_MTS= 10 // default is 10 minutes. Job related to firmware updates takes more time. 

1. To Run the Controller
    `go run main.go controller`

2. To Run the API
    `go run main.go api`

	
