FROM debian:stretch
RUN apt-get update && \
    apt-get install -y \
    multipath-tools \
    e2fsprogs \
    btrfs-progs \
    xfsprogs \
    iproute2 \
    util-linux \
    systemd \
    libisns0 \
    dmidecode

# create install paths
RUN mkdir -p /opt/hpe-storage/etc/ && mkdir /opt/hpe-storage/log/
RUN mkdir -p /opt/hpe-storage-mounts/

# add config files for each plugin type
ADD [ "config/*", "/opt/hpe-storage/etc/" ]

# add binaries
ADD [ "dockervolumed", "/opt/hpe-storage/dockervolumed" ]
ADD [ "flexvolume-driver", "/opt/hpe-storage/flexvolume" ]

# add host conformance checks and configuration
ADD [ "conform/*", "/opt/hpe-storage/lib/" ]

# add diag script
ADD [ "diag/*", "/opt/hpe-storage/bin/" ]

# add config files to tune multipath settings
ADD [ "tune/*", "/opt/hpe-storage/nimbletune/"]

# entrypoint
ADD [ "conform/entrypoint.sh", "/entrypoint.sh" ]
RUN [ "chmod", "+x", "/entrypoint.sh"]

ENTRYPOINT [ "/entrypoint.sh" ]