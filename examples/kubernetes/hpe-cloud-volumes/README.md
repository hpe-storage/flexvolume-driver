# HPE Cloud Volumes StorageClass parameters
A `StorageClass` is used to provision or clone an HPE Cloud Volumes-backed persistent volume. It can also be used to import an existing HPE Cloud Volumes volume or clone of a snapshot into the Kubernetes cluster. The parameters are grouped below by those same workflows.

A sample [StorageClass](sc-cv.yaml) is provided.

**Note:** These are optional parameters.

## Common parameters for Provisioning and Cloning
These parameters are mutable betweeen a parent volume and creating a clone from a snapshot.

| Parameter | String | Description |
| --------- | ------ | ----------- |
| nameSuffix| Text   | Suffix to append to Cloud Volumes. |
| destroyOnRm | Boolean | Indicates the backing Cloud volume (including snapshots) should be destroyed when the PVC is deleted. |
| limitIOPS | Integer | The IOPS limit of the volume. The IOPS limit should be in the range 300 to 50000. |
| perfPolicy | Text | The name of the performance policy to assign to the volume. Default example performance policies include "Other, Exchange, Oracle, SharePoint, SQL, Windows File Server". |
| protectionTemplate | Text | The name of the protection template to assign to the volume. Default examples of protection templates include "daily:3, daily:7, daily:14, hourly:6, hourly:12, hourly:24, twicedaily:4, twicedaily:8, twicedaily:14, weekly:2, weekly:4, weekly:8, monthly:3, monthly:6, monthly:12 or none". |
| volumeType | Text | Cloud Volume type. Supported types are PF and GPF. |

## Provisioning parameters
These parameters are immutable for clones once a volume has been created.

| Parameter | String | Description |
| --------- | ------ | ----------- |
| fsOwner | userId:groupId | The user id and group id that should own the root directory of the filesystem. |
| fsMode | Octal digits | 1 to 4 octal digits that represent the file mode to be applied to the root directory of the filesystem. |
| encryption | Boolean | Indicates that the volume should be encrypted. |

## Cloning parameters
Cloning supports two modes of cloning. Either use `cloneOf` and reference a PVC in the current namespace or use `importVolAsClone` and reference a Cloud volume name to clone and import to Kubernetes.

| Parameter | String | Description |
| --------- | ------ | ----------- |
| cloneOf | Text | The name of the PV to be cloned. `cloneOf` and `importVolAsClone` are mutually exclusive. |
| importVolAsClone | Text | The name of the Cloud Volume volume to clone and import. `importVolAsClone` and `cloneOf` are mutually exclusive. |
| snapshot | Text | The name of the snapshot to base the clone on. This is optional. If not specified, a new snapshot is created. |
| createSnapshot | Boolean | Indicates that a new snapshot of the volume should be taken matching the name provided in the `snapshot` parameter. If the `snapshot` parameter is not specified, a default name will be created. |
| snapshotPrefix | Text | A prefix to add to the beginning of the snapshot name. |
| replStore | Text | Replication store name. Should be used with importVolAsClone parameter to clone a replica volume |

## Import parameters
Importing volumes to Kubernetes requires the source Cloud volume to be not attached to any nodes. All previous Access Control Records will be stripped from the volume when put under control of the HPE Volume Driver for Kubernetes FlexVolume Plugin.

| Parameter | String | Description |
| --------- | ------ | ----------- |
| importVol | Text | The name of the Cloud volume to import. |
| forceImport | Boolean | Forces the import of a volume that is provisioned by another K8s cluster but not attached to any nodes. |
