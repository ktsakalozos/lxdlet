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
crictl  -r /var/tmp/lxdlet.sock runs test-resources/sandbox-config.json

#echo "******** list sandboxes ***************"
crictl  -r /var/tmp/lxdlet.sock sandboxes

#echo "******** create contaner ***************"
crictl  -r /var/tmp/lxdlet.sock create integrationtestcontaner test-resources/container-config.json test-resources/sandbox-config.json
crictl -r /var/tmp/lxdlet.sock ps
lxc list

#echo "******** start ***************"
crictl  -r /var/tmp/lxdlet.sock start integrationtestcontaner
crictl -r /var/tmp/lxdlet.sock ps
lxc list


#echo "******** stop container ***************"
crictl  -r /var/tmp/lxdlet.sock stop integrationtestcontaner
crictl -r /var/tmp/lxdlet.sock ps
lxc list

#echo "******** rm container ***************"
crictl  -r /var/tmp/lxdlet.sock rm integrationtestcontaner
crictl -r /var/tmp/lxdlet.sock ps
lxc list

#echo "******** stop sandbox ***************"
crictl  -r /var/tmp/lxdlet.sock stops hdishd83djaidwnduwk28bcsb
crictl -r /var/tmp/lxdlet.sock ps
lxc list

#echo "******** rm sandbox ***************"
crictl  -r /var/tmp/lxdlet.sock rms hdishd83djaidwnduwk28bcsb
crictl -r /var/tmp/lxdlet.sock ps
lxc list
