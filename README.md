# HPE Volume Driver for Kubernetes FlexVolume Plugin

HPE Volume Driver for Kubernetes FlexVolume Plugin leverages HPE Nimble Storage or HPE Cloud Volumes to provide scalable and persistent storage for stateful applications.

**Note:** The majority of the documentation of the HPE Volume Driver for Kubernetes FlexVolume Plugin has [moved to SCOD](https://scod.hpedev.io/flexvolume_driver/container_provider/). Sections below link directly to SCOD.

## Platform requirements
Supported platforms can be found in [platform requirements](https://scod.hpedev.io/flexvolume_driver/container_provider/#platform_requirements).

## Deploying to Kubernetes
The FlexVolume is preferably deployed to Kubernetes with the [Helm chart](https://hub.helm.sh/charts/hpe-storage/hpe-flexvolume-driver). Other methods are available at [Deploying to Kubernetes](https://scod.hpedev.io/flexvolume_driver/container_provider/#deploying_to_kubernetes).

## Using
Get started using the FlexVolume driver by setting up `StorageClass`, `PVC` API objects. See [Using](https://scod.hpedev.io/flexvolume_driver/container_provider/#using) for examples.

## Building
Instructions on how to build the FlexVolume driver from sources can be found in [BUILDING.md](BUILDING.md)

## Diagnostics
Logging and other troubleshooting steps can be in [Diagnostics](https://scod.hpedev.io/flexvolume_driver/container_provider/#diagnostics)

## Support
The HPE Volume Driver for Kubernetes FlexVolume Plugin is supported software by Hewlett Packard Enterprise. Reach out to your HPE representation to be connected with the support organization with any general issue you need help resolving.

We also encourage open collaboration, file issues, questions or feature requests [here](https://github.com/hpe-storage/flexvolume-driver/issues). You may also join our Slack community to chat with HPE folks close to this project. We hang out in `#NimbleStorage` and `#Kubernetes` at [slack.hpedev.io](https://slack.hpedev.io/).

## Contributing
We value all feedback and contributions. If you find any issues or want to contribute, please feel free to open an issue or file a PR. More details in [CONTRIBUTING.md](CONTRIBUTING.md)

## License
This is open source software licensed using the Apache License 2.0. Please see [LICENSE](LICENSE) for details.
