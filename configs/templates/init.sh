#!/bin/bash
echo "Depecrated. Will be removed from Future release."
echo "This is script is replaced by Boot Actions and Metamorph Agent"
sleep 30
curl -d '{ "UUID" : "{{.NodeUUID}}", "Name": "demo-provisioning","DeployStatus": "OSINSTALLED" }' -H "Content-Type: application/json" -X  POST http://{{.ProvisioningIP}}:{{.ProvisionerPort}}/nodes/{{.NodeUUID}}

INSTALL_LOCK=/var/install/metamorph-bootactions.lock
curl -s http://{{.ProvisioningIP}}:{{.HTTPPort}}/deploy_airskiff.sh | bash  > /var/log/airskiff_installation.log  2>&1

if [ $? == 0 ]; then
     touch $INSTALL_LOCK
     echo " AirSkiff Installation Success" >> /var/log/airskiff_installation.log
     curl -X POST --data-binary @/root/.kube/config -H "Content-type: text/x-yaml" http://{{.ProvisioningIP}}:{{.ProvisionerPort}}/update_kube_config/{{.NodeUUID}}
fi