module github.com/manojkva/metamorph-redfish-plugin

go 1.13

replace opendev.org/airship/go-redfish/client => /root/go/src/opendev.org/airship/go-redfish/client // Use opendev/org/airship/go-redfish refs/changes/77/737177/3

replace github.com/manojkva/go-redfish-api-wrapper => /root/go/src/github.com/manojkva/go-redfish-api-wrapper //Replace the above redfish PS in the local dir of api-wrapper too

replace github.com/bm-metamorph/MetaMorph => /go/src/github.com/bm-metamorph/MetaMorph 

replace github.com/manojkva/metamorph-plugin => /root/go/src/github.com/manojkva/metamorph-plugin

require (
	github.com/bm-metamorph/MetaMorph v0.0.0
	github.com/go-resty/resty/v2 v2.3.0
	github.com/hashicorp/go-hclog v0.14.1
	github.com/hashicorp/go-plugin v1.3.0
	github.com/hashicorp/go-version v1.2.1
	github.com/manojkva/go-redfish-api-wrapper v1.0.6
	github.com/manojkva/metamorph-plugin v0.0.0-00010101000000-000000000000
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.6.1
	go.uber.org/zap v1.15.0
	golang.org/x/tools v0.0.0-20200708183856-df98bc6d456c // indirect
	opendev.org/airship/go-redfish/client v0.0.0 // indirect
)
