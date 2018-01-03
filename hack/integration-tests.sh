#!/bin/bash
#

set -x

echo "***************************************"
echo "******* Version call ******************"
echo "***************************************"
crictl -r /var/tmp/lxdlet.sock version

echo "***************************************"
echo "**************** ps *******************"
echo "***************************************"
crictl --debug -r /var/tmp/lxdlet.sock ps

echo "***************************************"
echo "******** create sandbox ***************"
echo "***************************************"
crictl  -r /var/tmp/lxdlet.sock create integrationtestsandbox test-resources/container-config.json test-resources/sandbox-config.json
crictl -r /var/tmp/lxdlet.sock ps
