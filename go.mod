module github.com/bm-metamorph/MetaMorph

go 1.13

require (
	github.com/gin-contrib/zap v0.0.1
	github.com/gin-gonic/gin v1.6.3
	github.com/go-playground/validator/v10 v10.3.0 // indirect
	github.com/go-resty/resty/v2 v2.0.0
	github.com/golang/protobuf v1.4.2
	github.com/google/uuid v1.1.1
	github.com/hashicorp/go-version v1.2.1
	github.com/jinzhu/gorm v1.9.12
	github.com/manojkva/go-redfish-api-wrapper v1.0.6
	github.com/mitchellh/go-homedir v1.1.0
	github.com/natefinch/lumberjack v2.0.0+incompatible
	github.com/pkg/errors v0.8.1
	github.com/spf13/cobra v0.0.7
	github.com/spf13/viper v1.6.2
	github.com/stretchr/testify v1.4.0
	go.uber.org/zap v1.15.0
	golang.org/x/crypto v0.0.0-20200323165209-0ec3e9974c59 // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	golang.org/x/sys v0.0.0-20200519105757-fe76b779f299 // indirect
	google.golang.org/grpc v1.28.0
	google.golang.org/protobuf v1.23.0
	gopkg.in/yaml.v2 v2.3.0 // indirect
	k8s.io/api v0.18.1 // indirect
	opendev.org/airship/go-redfish/client v0.0.0-0
)

replace opendev.org/airship/go-redfish/client => /root/go/src/opendev.org/airship/go-redfish/client // Use opendev/org/airship/go-redfish refs/changes/77/737177/3

replace github.com/manojkva/go-redfish-api-wrapper => /root/go/src/github.com/manojkva/go-redfish-api-wrapper //Replace the above redfish PS in the local dir of api-wrapper too
