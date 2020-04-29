FROM registry.svc.ci.openshift.org/openshift/release:golang-1.13 AS builder
WORKDIR /go/src/bitbucket.com/metamorph
COPY . .
ENV GOFLAGS=""
RUN go build -o build/_output/bin/metamorph main.go

FROM ubuntu:18.04

COPY --from=builder /go/src/bitbucket.com/metamorph/build/_output/bin/metamorph /
COPY configs/ src/bitbucket.com/metamorph/configs/
LABEL io.k8s.display-name="MetaMorph BareMetal Provisoner" \
      io.k8s.description="This is the image for the MetaMorph BareMetal Provisioner."
