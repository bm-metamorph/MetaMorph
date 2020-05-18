# MetaMorph

Lifecycle your BareMetal


## Setting up Development Environment

Setup ENV variables

export METAMORPH_CONFIGPATH=<path of config.yaml location>
export REDFISH_SLEEPTIME_SECS=10 //time duration in between subsequent redfish API calls

1. To Run the Controller
    `go run main.go controller`

2. To Run the API
    `go run main.go api`

	
