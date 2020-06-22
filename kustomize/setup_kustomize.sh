#!/bin/bash


# Copy kustomize binary to host
# Make sure to mount host's /user/local/bin to pod's /opt/bin/
cp /opt/metamorph/kustomize/bin/kustomize /opt/bin/kustomize

# Copy kustomize plugin to host
# make sure to mount host's /root/.config to pod's /opt/.config

mkdir -p /opt/.config/kustomize/plugin/metamorph.io/v1/userdata
cp -r /opt/metamorph/kustomize/plugin/metamorph.io/v1/userdata  /opt/.config/kustomize/plugin/metamorph.io/v1/userdata


echo "kustomize setup complete"

