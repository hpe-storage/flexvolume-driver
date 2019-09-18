# HPE Nimble Storage StorageClass parameters
A `StorageClass` is used to provision or clone an HPE Nimble Storage-backed persistent volume. It can also be used to import an existing HPE Nimble Storage volume or clone of a snapshot into the Kubernetes cluster. The parameters are grouped below by those same workflows.

A sample [StorageClass](sc-nimble.yaml) is provided.

**Note:** These are optional parameters.

## Common parameters for Provisioning and Cloning
These parameters are mutable betweeen a parent volume and creating a clone from a snapshot.

| Parameter | String | Description |
| --------- | ------ | ----------- |
| nameSuffix| Text   | Suffix to append to Nimble volumes. Defaults to .docker |
| destroyOnRm | Boolean | Indicates the backing Nimble volume (including snapshots) should be destroyed when the PVC is deleted. |
| limitIOPS | Integer | The IOPS limit of the volume. The IOPS limit should be in the range 256 to 4294967294, or -1 for unlimited (default). |
| limitMBPS | Integer | The MB/s throughput limit for the volume. |
| description | Text | Text to be added to the volume's description on the Nimble array. |
| perfPolicy | Text | The name of the performance policy to assign to the volume. Default example performance policies include "Backup Repository", "Exchange 2003 data store", "Exchange 2007 data store", "Exchange 2010 data store", "Exchange log", "Oracle OLTP", "Other Workloads", "SharePoint", "SQL Server", "SQL Server 2012", "SQL Server Logs". |
| protectionTemplate | Text | The name of the protection template to assign to the volume. Default examples of protection templates include "Retain-30Daily", "Retain-48Hourly-30aily-52Weekly", and "Retain-90Daily". |
| folder | Text | The name of the Nimble folder in which to place the volume. |
| thick | Boolean | Indicates that the volume should be thick provisioned. |
| dedupeEnabled | Boolean | Indicates that the volume should enable deduplication. |
| syncOnUnmount | Boolean | Indicates that a snapshot of the volume should be synced to the replication partner each time it is detached from a node. |

**Note**: Performance Policies, Folders and Protection Templates are Nimble specific constructs that can be created on the Nimble array itself to address particular requirements or workloads. Please consult with the storage admin or read the admin guide found on [HPE InfoSight](https://infosight.hpe.com).

## Provisioning parameters
These parameters are immutable for clones once a volume has been created.

| Parameter | String | Description |
| --------- | ------ | ----------- |
| fsOwner | userId:groupId | The user id and group id that should own the root directory of the filesystem. |
| fsMode | Octal digits | 1 to 4 octal digits that represent the file mode to be applied to the root directory of the filesystem. |
| encryption | Boolean | Indicates that the volume should be encrypted. |
| pool | Text | The name of the pool in which to place the volume. |

## Cloning parameters
Cloning supports two modes of cloning. Either use `cloneOf` and reference a PVC in the current namespace or use `importVolAsClone` and reference a Nimble volume name to clone and import to Kubernetes.

| Parameter | String | Description |
| --------- | ------ | ----------- |
| cloneOf | Text | The name of the PV to be cloned. `cloneOf` and `importVolAsClone` are mutually exclusive. |
| importVolAsClone | Text | The name of the Nimble volume to clone and import. `importVolAsClone` and `cloneOf` are mutually exclusive. |
| snapshot | Text | The name of the snapshot to base the clone on. This is optional. If not specified, a new snapshot is created. |
| createSnapshot | Boolean | Indicates that a new snapshot of the volume should be taken matching the name provided in the `snapshot` parameter. If the `snapshot` parameter is not specified, a default name will be created. |
| snapshotPrefix | Text | A prefix to add to the beginning of the snapshot name. |
| destroyOnDetach | Boolean | Indicates that the Nimble volume (including snapshots) backing this volume should be destroyed when this volume is unmounted or detached. |

## Import parameters
Importing volumes to Kubernetes requires the source Nimble volume to be offline. All previous Access Control Records and Initiator Groups will be stripped from the volume when put under control of the HPE Volume Driver for Kubernetes FlexVolume Plugin.

| Parameter | String | Description |
| --------- | ------ | ----------- |
| importVol | Text | The name of the Nimble volume to import. |
| snapshot | Text | The name of the Nimble snapshot to restore the imported volume to after takeover. If not specified, the volume will not be restored. |
| restore  | Boolean | Restores the volume to the last snapshot taken on the volume. |
| takeover | Boolean | Indicates the current group will takeover ownership of the Nimble volume and volume collection. This should be performed against a downstream replica. |
| reverseRepl | Boolean | Reverses the replication direction so that writes to the Nimble volume are replicated back to the group where it was replicated from. |
| forceImport | Boolean | Forces the import of a volume that is not owned by the group and is not part of a volume collection. If the volume is part of a volume collection, use takeover instead.

**Note:** HPE Nimble Docker Volume workflows works with 1-1 mapping between volume and volume collection.
