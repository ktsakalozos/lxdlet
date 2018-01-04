#!/bin/bash
#

set -x

# echo "******* Version call ******************"
crictl -r /var/tmp/lxdlet.sock version

# echo "**************** ps *******************"
crictl --debug -r /var/tmp/lxdlet.sock ps

#echo "******** list sandboxes ***************"
crictl  -r /var/tmp/lxdlet.sock sandboxes

#echo "******** add sandbox ***************"
SANDBOX=$(crictl  -r /var/tmp/lxdlet.sock runs test-resources/sandbox-config.json)
echo "Sandbox created: " $SANDBOX

#echo "******** list sandboxes ***************"
crictl  -r /var/tmp/lxdlet.sock sandboxes

#echo "******** create contaner ***************"
CONTAINER=$(crictl  -r /var/tmp/lxdlet.sock create $SANDBOX test-resources/container-config.json test-resources/sandbox-config.json)
echo "Container created: " $CONTAINER
crictl -r /var/tmp/lxdlet.sock ps
lxc list

#echo "******** start ***************"
crictl  -r /var/tmp/lxdlet.sock start $CONTAINER
crictl -r /var/tmp/lxdlet.sock ps
lxc list


#echo "******** stop container ***************"
crictl  -r /var/tmp/lxdlet.sock stop $CONTAINER
crictl -r /var/tmp/lxdlet.sock ps
lxc list

#echo "******** rm container ***************"
crictl  -r /var/tmp/lxdlet.sock rm $CONTAINER
crictl -r /var/tmp/lxdlet.sock ps
lxc list

#echo "******** stop sandbox ***************"
crictl  -r /var/tmp/lxdlet.sock stops $SANDBOX
crictl -r /var/tmp/lxdlet.sock ps
lxc list

#echo "******** rm sandbox ***************"
crictl  -r /var/tmp/lxdlet.sock rms $SANDBOX
crictl -r /var/tmp/lxdlet.sock ps
lxc list
