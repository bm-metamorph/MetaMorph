module github.com/manojkva/metamorph-isogen-plugin

go 1.13

replace opendev.org/airship/go-redfish/client => /root/go/src/opendev.org/airship/go-redfish/client // Use opendev/org/airship/go-redfish refs/changes/77/737177/3

replace github.com/bm-metamorph/MetaMorph => /go/src/github.com/bm-metamorph/MetaMorph

replace github.com/manojkva/metamorph-plugin => /root/go/src/github.com/manojkva/metamorph-plugin

require (
	github.com/bm-metamorph/MetaMorph v0.0.0-00010101000000-000000000000
	github.com/hashicorp/go-hclog v0.14.1
	github.com/hashicorp/go-plugin v1.3.0
	github.com/manojkva/metamorph-plugin v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.6.1
	go.uber.org/zap v1.15.0
)
