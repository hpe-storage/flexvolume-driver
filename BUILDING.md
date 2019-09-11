
# Building the HPE Volume Driver for Kubernetes FlexVolume Plugin

- Clone the repo : git clone <https://github.com/hpe-storage/flexvolume-driver>
- cd to flexvolume-driver
- Turn on go modules support `export GO111MODULES=on`
- Set GOOS `export GOOS=linux`(optional)
- Set CONTAINER_REGISTRY env to point to your image registry, if other than docker.io(default)
- Set GOPATH, as go binaries are placed under $(GOPATH)/bin which is added to $(PATH)
- Run `make all`

Note1: Minimum go version of 1.12 required.
Note2: tests are only supported on Linux platform
