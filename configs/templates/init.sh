#!/bin/bash
sleep 30
curl -d '{ "UUID" : "{{.NodeUuid}}", "Name": "demo-provisioning","DeployStatus": "OSINSTALLED" }' -H "Content-Type: application/json" -X  POST http://{{.ProvisioningIP}}:{{.ProvisionerPort}}/nodes/{{.NodeUuid}}

INSTALL_LOCK=/var/install/metamorph-bootactions.lock
curl -s http://{{.ProvisioningIP}}:{{.HttpPort}}/deploy_airskiff.sh | bash  > /var/log/airskiff_installation.log  2>&1

if [ $? == 0 ]; then
     touch $INSTALL_LOCK
     echo " AirSkiff Installation Success" >> /var/log/airskiff_installation.log
     curl -X POST --data-binary @/root/.kube/config -H "Content-type: text/x-yaml" http://{{.ProvisioningIP}}:{{.ProvisionerPort}}/update_kube_config/{{.NodeUuid}}
fi