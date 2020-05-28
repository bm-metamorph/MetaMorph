FROM registry.svc.ci.openshift.org/openshift/release:golang-1.13 AS builder
WORKDIR /go/src/bitbucket.com/metamorph
COPY . .

#temp fix to resolve dependencies
COPY tmp/go-redfish-api-wrapper /Users/aneesh/go/src/github.com/manojkva/go-redfish-api-wrapper
COPY tmp/go-redfish /Users/aneesh/go/src/opendev.org/airship/go-redfish


ENV GOFLAGS=""
RUN go build -o build/_output/bin/metamorph main.go
RUN go build -o build/_output/bin/metamorph_agent pkg/agent/agent.go

FROM ubuntu:18.04 as exporter

COPY --from=builder /go/src/bitbucket.com/metamorph/build/_output/bin/metamorph /usr/local/bin/
COPY --from=builder /go/src/bitbucket.com/metamorph/build/_output/bin/metamorph_agent /opt/metamorph/assets/files/metamorph_agent
COPY --from=builder /go/src/bitbucket.com/metamorph/configs /opt/metamorph/configs

RUN apt update && apt  -y upgrade && apt-get -y install --no-install-recommends rsync mkisofs && \
    apt clean && rm -rf /var/lib/apt/lists/*

LABEL io.k8s.display-name="MetaMorph BareMetal Provisoner" \
      io.k8s.description="This is the image for the MetaMorph BareMetal Provisioner."
