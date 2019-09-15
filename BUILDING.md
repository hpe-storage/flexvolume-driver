# Building the HPE Volume Driver for Kubernetes FlexVolume Plugin
The driver is written in Go and require a recent (Go 1.12 minimum) functional Go build environment.

Clone the repo:
```
git clone https://github.com/hpe-storage/flexvolume-driver
cd flexvolume-driver
```
Turn on go modules support:
```
export GO111MODULES=on
export GOOS=linux
```
**Note:** Setting `GOOS` is optional as the binary runs on any platform, but tests are only supported on Linux.

Set `CONTAINER_REGISTRY` environment variable to point to your image registry, if other than docker.io (default).
Make sure `GOPATH` is set, as go binaries are placed under `$GOPATH/bin` which is added to `$PATH`.

```
make all
```
