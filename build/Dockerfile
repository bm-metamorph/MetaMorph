FROM registry.svc.ci.openshift.org/openshift/release:golang-1.13 AS builder
WORKDIR /go/src/github.com/bm-metamorph/MetaMorph
COPY . .

#temp fix to resolve dependencies
COPY tmp/go-redfish-api-wrapper /root/go/src/github.com/manojkva/go-redfish-api-wrapper
COPY tmp/go-redfish /root/go/src/opendev.org/airship/go-redfish


ENV GOFLAGS=""
RUN go build -o build/_output/bin/metamorph main.go
RUN go build -o build/_output/bin/metamorph_agent pkg/agent/agent.go


FROM registry.svc.ci.openshift.org/openshift/release:golang-1.13 AS kustomize-builder

#Build Kustomize with Plugin support
ENV GOFLAGS=""
ENV GO111MODULE=on

COPY kustomize/plugin/metamorph.io/v1/userdata  /opt/metamorph/kustomize/plugin/metamorph.io/v1/userdata
RUN go get sigs.k8s.io/kustomize/kustomize/v3@v3.5.4
RUN pushd /opt/metamorph/kustomize/plugin/metamorph.io/v1/userdata && \
    go build -buildmode plugin -o UserData.so UserData.go

FROM ubuntu:18.04 as exporter

COPY --from=builder /go/src/github.com/bm-metamorph/MetaMorph/build/_output/bin/metamorph /usr/local/bin/
COPY --from=builder /go/src/github.com/bm-metamorph/MetaMorph/build/_output/bin/metamorph_agent /opt/metamorph/assets/files/metamorph_agent
COPY --from=builder /go/src/github.com/bm-metamorph/MetaMorph/configs /opt/metamorph/configs
COPY --from=kustomize-builder /go/bin/kustomize /opt/metamorph/kustomize/bin/kustomize
COPY --from=kustomize-builder /opt/metamorph/kustomize/plugin/metamorph.io/v1/userdata  /opt/metamorph/kustomize/plugin/metamorph.io/v1/userdata
COPY kustomize/setup_kustomize.sh /usr/local/bin/setup_kustomize.sh

RUN apt update && apt  -y upgrade && apt-get -y install --no-install-recommends rsync mkisofs vim curl&& \
    apt clean && rm -rf /var/lib/apt/lists/*

LABEL io.k8s.display-name="MetaMorph BareMetal Provisoner" \
      io.k8s.description="This is the image for the MetaMorph BareMetal Provisioner."
