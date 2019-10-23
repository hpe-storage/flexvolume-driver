#!/bin/sh

# Tolerate zero errors and tell the world all about it
set -xe

# Check protocol and modify default for fc
if [ "${PROTOCOL}" = "fc" ]; then
    sed -i -e 's/Environment=PROTOCOL=iscsi/Environment=PROTOCOL=fc/g' \
    /opt/hpe-storage/lib/hpe-storage-node.service
fi

# Apply workaround for Rancher RKE(kubelet in container) related to
# https://github.com/kubernetes/kubernetes/issues/65825
if [ "${FLAVOR}" = "rancher" ]; then
    sed -i -e 's/Environment=FLAVOR=k8s/Environment=FLAVOR=rancher/g' \
        /opt/hpe-storage/lib/hpe-storage-node.service
    cp -f "/opt/hpe-storage/nimbletune/multipath.conf.upstream" \
        /usr_local/local/bin/multipath.conf
fi

# Copy HPE Storage Node Conformance checks and conf in place
cp -f "/opt/hpe-storage/lib/hpe-storage-node.service" \
      /lib/systemd/system/hpe-storage-node.service
cp -f "/opt/hpe-storage/lib/hpe-storage-node.sh" \
      /usr_local/local/bin/hpe-storage-node.sh
chmod +x /usr_local/local/bin/hpe-storage-node.sh

echo "running conformance checks..."
# Reload and run!
systemctl daemon-reload
systemctl restart hpe-storage-node

# Copy the hpe log collector script
echo "copy hpe log collector script..."
cp -f "/opt/hpe-storage/bin/hpe-logcollector.sh" \
      /usr_local/local/bin/hpe-logcollector.sh
chmod +x /usr_local/local/bin/hpe-logcollector.sh

echo "starting flexvolume plugin..."
# Serve! Serve!!!
exec /opt/hpe-storage/dockervolumed