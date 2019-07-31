# HPE Nimble Storage Flexvolume Driver Worklows

## Storage Class Properties

A storage class is used to create or clone an HPE Nimble Storage-backed persistent volume.  It can also be used to import an existing HPE Nimble Storage volume or clone of a snapshot into the Kubernetes cluster.  The parameters are grouped below by workflow.

### 1. Valid Parameters for Create

```markdown
- limitIOPS:"X"
    - Where X is the IOPS limit of the volume. The IOPS limit should be in the range [256, 4294967294], or -1 for unlimited.

- limitMBPS:"X"
    - Where X is the MB/s throughput limit for this volume. If both limitIOPS and limitMBPS are specified, limitIOPS must be specified first.
    - destroyOnRm:"true"
        - Indicates that the Nimble volume (including snapshots) backing this volume should be destroyed when this volume is deleted.

- sizeInGiB:"X"
    - Where X is the size of volume specified in GiB.

- size:"X"
    - Where X is the size of volume specified in GiB (short form of sizeInGiB).

- fsOwner:"X"
    - Where X is the user id and group id that should own the root directory of the filesystem, in the form of [userId:groupId].

- fsMode:"X"
    - Where X is 1 to 4 octal digits that represent the file mode to be applied to the root directory of the filesystem.

- description:"X"
    - Where X is the text to be added to volume description (optional).

- perfPolicy:"X"
    - Where X is the name of the performance policy (optional).
    - Examples of Performance Policies include: Exchange 2003 data store, Exchange 2007 data store, Exchange log, SQL Server, SharePoint, Exchange 2010 data store, SQL Server Logs, SQL Server 2012, Oracle OLTP, Windows File Server, Other Workloads, Backup Repository.

- pool:"X"
    - Where X is the name of pool in which to place the volume (optional).

- folder:"X"
    - Where X is the name of folder in which to place the volume (optional).

- encryption:"true"
    - Indicates that the volume should be encrypted (optional, dedupe and encryption are mutually exclusive).

- thick:"true"
    - Indicates that the volume should be thick provisioned (optional, dedupe and thick are mutually exclusive).

- dedupe:"true"
    - Indicates that the volume should be deduplicated.

- protectionTemplate:"X"
    - Where X is the name of the protection template (optional).
    - Examples of Protection Templates include: Retain-30Daily, Retain-90Daily, Retain-48Hourly-30Daily-52Weekly.
```

#### Example

```yaml

---
kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
    name: oltp-prod
provisioner: hpe.com/nimble
parameters:
    thick: "false"
    dedupe: "true"
    perfPolicy: "SQL Server"
    protectionTemplate: "Retain-48Hourly-30Daily-52Weekly"
```

### 2. Valid Parameters for Clone

```markdown
- cloneOf:"X"
    - Where X is the name of the Docker Volume to be cloned.

- snapshot:"X"
    - Where X is the name of the snapshot to base the clone on. This is optional. If not specified, a new snapshot is created.

- createSnapshot:"true"
    - Indicates that a new snapshot of the volume should be taken and used for the clone (optional).

- destroyOnRm:"true"
    - Indicates that the Nimble volume (including snapshots) backing this volume should be destroyed when this volume is deleted.
```

#### Example

```yaml
---
  kind: StorageClass
  apiVersion: storage.k8s.io/v1
  metadata:
    name: oltp-dev-clone-of-prod
  provisioner: hpe.com/nimble
  parameters:
    limitIOPS: "1000"
    cloneOf: "oltp-prod-1adee106-110b-11e8-ac84-00505696c45f"
    destroyOnRm: "true"
```

### 3. Valid Parameters for Import Clone of Snapshot

```markdown
- importVolAsClone:"X"
    - Where X is the name of the Nimble Volume and Nimble Snapshot to clone and import.

- snapshot:"X"
    - Where X is the name of the Nimble Snapshot to clone and import (optional, if missing, will use the most recent snapshot).

- createSnapshot:"true"
    - Indicates that a new snapshot of the volume should be taken and used for the clone (optional).

- pool:"X"
    - Where X is the name of the pool in which the volume to be imported resides (optional).

- folder:"X"
    - Where X is the name of the folder in which the volume to be imported resides (optional).

- destroyOnRm:"true"
    - Indicates that the Nimble volume (including snapshots) backing this volume should be destroyed when this volume is deleted.

- destroyOnDetach
    - Indicates that the Nimble volume (including snapshots) backing this volume should be destroyed when this volume is unmounted or detached.
```

#### Example

```yaml
---
  kind: StorageClass
  apiVersion: storage.k8s.io/v1
  metadata:
    name: import-clone-legacy-prod
 provisioner: hpe.com/nimble
  parameters:
    pool: "flash"
    importVolAsClone: "production-db-vol"
    destroyOnRm: "true"
```

### 4. Valid Parameters for Import Volume

```markdown
- importVol:"X"
    - Where X is the name of the Nimble volume to import.

- pool:"X"
    - Where X is the name of the pool in which the volume to be imported resides (optional).

- folder:"X"
    - Where X is the name of the folder in which the volume to be imported resides (optional).

- forceImport:"true"
    - Force the import of the volume. Note that overwrites application metadata (optional).

- restore
    - Restores the volume to the last snapshot taken on the volume (optional).

- snapshot:"X"
    - Where X is the name of the snapshot which the volume will be restored to. Only used with -o restore (optional).

- takover
    - Indicates the current group will takeover the ownership of the Nimble volume and volume collection (optional).

- reverseRepl
    - Reverses the replication direction so that writes to the Nimble volume are replicated back to the group where it was replicated from (optional).
```

#### Example

```yaml
---
  kind: StorageClass
  apiVersion: storage.k8s.io/v1
  metadata:
    name: import-clone-legacy-prod
  provisioner: hpe.com/nimble
  parameters:
    pool: "flash"
    importVol: "production-db-vol"
```

### 5. Valid Parameter for allowOverrides

#### Example

```yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
 name: mysc900
provisioner: hpe.com/nimble
parameters:
 description: "Volume from doryd"
 size: "900"
 dedupe: "false"
 destroyOnRm: "true"
 perfPolicy: "Windows File Server"
 folder: "mysc900"
 allowOverrides: snapshot,description,limitIOPS,size,perfPolicy,thick,folder
 ```

#### Persistent Volume Claim

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
 name: mypvc901
 annotations:
    hpe.com/description: "This is my custom description"
    hpe.com/limitIOPS: "8000"
    hpe.com/cloneOfPVC: myPVC
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
  storageClassName: mysc900```
