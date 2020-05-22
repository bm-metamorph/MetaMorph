#!/bin/bash
set -xe

# Set permissions o+r, beacause these files need to be readable
# for Armada in the container
METAMORPH_PERMISSIONS=$(stat --format '%a' metamorph-armada-manifest.yaml)
KUBE_CONFIG_PERMISSIONS=$(stat --format '%a' ~/.kube/config)

sudo chmod 0644 metamorph-armada-manifest.yaml
sudo chmod 0644 ~/.kube/config

# In the event that this docker command fails, we want to continue the script
# and reset the file permissions.
set +e

# Download latest Armada image and deploy Airship components
docker run --rm --net host -p 8000:8000 --name armada \
    -v ~/.kube/config:/armada/.kube/config \
    -v "$(pwd)"/metamorph:/metamorph \
    -v "$(pwd)"/metamorph-armada-manifest.yaml:/metamorph-armada-manifest.yaml \
    quay.io/airshipit/armada:latest-ubuntu_bionic \
    apply /metamorph-armada-manifest.yaml

# Set back permissions of the files
sudo chmod "${AIRSKIFF_PERMISSIONS}" metamorph-armada-manifest.yaml
sudo chmod "${KUBE_CONFIG_PERMISSIONS}" ~/.kube/config
