module github.com/bm-metamorph/MetaMorph

go 1.13

require (
	github.com/gin-contrib/zap v0.0.1
	github.com/gin-gonic/gin v1.6.3
	github.com/golang/protobuf v1.4.2
	github.com/google/uuid v1.1.1
	github.com/hashicorp/go-hclog v0.14.1
	github.com/hashicorp/go-plugin v1.3.0
	github.com/jinzhu/gorm v1.9.12
	github.com/manojkva/metamorph-plugin v0.0.0-00010101000000-000000000000
	github.com/mitchellh/go-homedir v1.1.0
	github.com/spf13/cobra v0.0.7
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.4.0
	go.uber.org/zap v1.15.0
	google.golang.org/grpc v1.30.0
	google.golang.org/protobuf v1.25.0
	opendev.org/airship/go-redfish/client v0.0.0-0
)

replace opendev.org/airship/go-redfish/client => /root/go/src/opendev.org/airship/go-redfish/client // Use opendev/org/airship/go-redfish refs/changes/77/737177/3


replace github.com/manojkva/metamorph-plugin => /root/go/src/github.com/manojkva/metamorph-plugin
