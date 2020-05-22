FROM registry.svc.ci.openshift.org/openshift/release:golang-1.13 AS builder
WORKDIR /go/src/bitbucket.com/metamorph
COPY . .
ENV GOFLAGS=""
RUN go build -o build/_output/bin/metamorph main.go
RUN go build -o build/_output/bin/metamorph_agent pkg/agent/agent.go

FROM ubuntu:18.04 as exporter

COPY --from=builder /go/src/bitbucket.com/metamorph/build/_output/bin/metamorph /usr/local/bin/
COPY --from=builder /go/src/bitbucket.com/metamorph/build/_output/bin/metamorph_agent /opt/metamorph/assets/files/metamorph_agent
COPY --from=builder /go/src/bitbucket.com/metamorph/configs /opt/metamorph/configs



LABEL io.k8s.display-name="MetaMorph BareMetal Provisoner" \
      io.k8s.description="This is the image for the MetaMorph BareMetal Provisioner."
