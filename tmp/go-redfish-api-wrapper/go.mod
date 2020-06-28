module github.com/manojkva/go-redfish-api-wrapper

go 1.13

replace opendev.org/airship/go-redfish/client => /home/ekuamaj/go/src/opendev.org/airship/go-redfish/client

require (
	github.com/antihax/optional v1.0.0
	github.com/stretchr/testify v1.4.0
	opendev.org/airship/go-redfish/client v0.0.0-0
	sigs.k8s.io/controller-runtime v0.5.1
)
