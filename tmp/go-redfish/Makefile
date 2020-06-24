.PHONY: code-gen
code-gen:
	rm -rf client
	openapi-generator generate -i ./spec/openapi.yaml -g go --package-name "client" --git-repo-id go-redfish/client --git-user-id airship --git-host opendev.org -o client/ -p enumClassPrefix=true
	mkdir -p api
	go run api_generator.go | gofmt > api/service_interface.go
	go generate api/service_interface.go

.PHONY: code-gen-docker
code-gen-docker:
	rm -rf client
	docker run --rm -v ${PWD}:/local openapitools/openapi-generator-cli:v4.2.2 generate \
         -i  /local/spec/openapi.yaml \
         -g go \
	 --package-name "client" \
         --git-repo-id go-redfish/client \
         --git-user-id airship \
         --git-host opendev.org \
	 -p enumClassPrefix=true \
         -o /local/client
	mkdir -p api
	go run api_generator.go | gofmt > api/service_interface.go
	go generate api/service_interface.go


.PHONY: deps
deps:
	go get github.com/stretchr/testify/assert
	go get golang.org/x/oauth2
	go get golang.org/x/net/context
	go get github.com/antihax/optional
