#!/bin/bash
#

set -x

# echo "******* Version call ******************"
crictl -r /var/tmp/lxdlet.sock version

# echo "**************** ps *******************"
crictl --debug -r /var/tmp/lxdlet.sock ps

#echo "******** create ***************"
crictl  -r /var/tmp/lxdlet.sock create integrationtestsandbox test-resources/container-config.json test-resources/sandbox-config.json
crictl -r /var/tmp/lxdlet.sock ps
lxc list

#echo "******** start ***************"
crictl  -r /var/tmp/lxdlet.sock start integrationtestsandbox
crictl -r /var/tmp/lxdlet.sock ps
lxc list

#echo "******** list sandboxes ***************"
crictl  -r /var/tmp/lxdlet.sock sandboxes


#echo "******** stop container ***************"
crictl  -r /var/tmp/lxdlet.sock stop integrationtestsandbox
crictl -r /var/tmp/lxdlet.sock ps
lxc list

#echo "******** stop sandbox ***************"
crictl  -r /var/tmp/lxdlet.sock stops integrationtestsandbox
crictl -r /var/tmp/lxdlet.sock ps
lxc list

#echo "******** rm container ***************"
crictl  -r /var/tmp/lxdlet.sock rm integrationtestsandbox
crictl -r /var/tmp/lxdlet.sock ps
lxc list

#echo "******** rm sandbox ***************"
crictl  -r /var/tmp/lxdlet.sock rms integrationtestsandbox
crictl -r /var/tmp/lxdlet.sock ps
lxc list
