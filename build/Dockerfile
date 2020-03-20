FROM registry.svc.ci.openshift.org/openshift/release:golang-1.10 AS builder
WORKDIR /go/src/github.com/metamorph
COPY . .
RUN go build -o build/_output/bin/metamorph cmd/manager/main.go

FROM ubuntu:18.04

COPY --from=builder /go/src/github.com/metamorph/build/_output/bin/metamorph /

LABEL io.k8s.display-name="MetaMorph BareMetal Provisoner" \
      io.k8s.description="This is the image for the MetaMorph BareMetal Provisioner."

