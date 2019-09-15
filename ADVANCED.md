# Advanced Configuration
This document some of the advanced configuration steps available to tweak behavior of the HPE Volume Driver for Kubernetes FlexVolume Plugin. 

## Set defaults at the compute node level
During normal operations, defaults are set in either the `ConfigMap` or in a `StorageClass` itself. The picking order is:

- StorageClass
- ConfigMap
- driver.json

Please see [DIAGNOSTICS.md](DIAGNOSTICS.md) to locate the driver for your particular environment. Add this object to the configuration file, `nimble.json`, for example:
```json
{
    "defaultOptions": [{"option1": "value1"}, {"option2": "value2"}]
}
```
Where `option1` and `option2` are valid backend volume plugin create options.

**Note:** It's highly recommended to control defaults with `StorageClass` API objects or the `ConfigMap`.

## Global options
Each driver supports setting certain "global" options in the `ConfigMap`. Some options are common, some are driver specific.

### Common

| Parameter | String | Description |
| --------- | ------ | ----------- |
| volumeDir | Text   | Root directory on the host to mount the volumes. This parameter needs correlation with the `podsmountdir` path in the `volumeMounts` stanzas of the deployment. | 
| logDebug  | Boolean | Turn on debug logging, set to false by default. |

### HPE Nimble Storage

| Parameter | String | Description |
| --------- | ------ | ----------- |
| TBD       | TBD    | TBD         |

### HPE Cloud Volumes

| Parameter | String | Description |
| --------- | ------ | ----------- |
| TBD       | TBD    | TBD         |
