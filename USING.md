# Using the HPE Volume Driver for Kubernetes FlexVolume Plugin

These instructions are provided as an example on how to use the HPE Volume Driver for Kubernetes FlexVolume Plugin with a HPE Nimble Storage Array.

The below YAML declarations are meant to be created with `kubectl create`. Either copy the content to a file on the host where `kubectl` is being executed, or copy & paste into the terminal, like this:

```
kubectl create -f-
< paste the YAML >
^D (CTRL + D)
```

**Note:**  Some of the examples supported by the HPE Volume Driver for Kubernetes FlexVolume Plugin are available in the [examples/kubernetes/hpe-nimble-storage](examples/kubernetes/hpe-nimble-storage) directory or [examples/kubernetes/cloud-volumes](examples/kubernetes/cloud-volumes) and all the HPE Nimble Storage Array Flexvolume `StorageClass` parameters can be found in [examples/kubernetes/hpe-nimble-storage](examples/kubernetes/hpe-nimble-storage).

To get started, create a `StorageClass` API object referencing the `hpe-secret` and defining additional (optional) `StorageClass` parameters:

## Sample storage classes

Sample storage classes can be found for [HPE Nimble Storage](examples/kubernetes/hpe-nimble-storage/sc-nimble.yaml), [Cloud Volumes](examples/kubernetes/hpe-cloud-volumes/sc-cv.yaml) and SimpliVity (TBD).

## Test and verify volume provisioning

Create a `StorageClass` with volume parameters as required.

```
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: sc-nimble
provisioner: hpe.com/nimble
parameters:
  description: "Volume from HPE FlexVolume driver"
  perfPolicy: "Other"
  limitIOPS: "76800"
```

Create a PersistentVolumeClaim. This makes sure a volume is created and provisioned on your behalf:

```
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: pvc-nimble
spec:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
  storageClassName: sc-nimble
```

Check that a new `PersistentVolume` is created based on your claim:

```
$ kubectl get pv
NAME                                            CAPACITY     ACCESS MODES   RECLAIM POLICY   STATUS   CLAIM               STORAGECLASS   REASON   AGE
sc-nimble-13336da3-7ca3-11e9-826c-00505693581f  10Gi         RWO            Delete           Bound    default/pvc-nimble  sc-nimble               3s
```

The above output means that the FlexVolume driver successfully provisioned a new volume and bound to the requesting `PVC` to a new `PV`. The volume is not attached to any node yet. It will only be attached to a node if a workload is scheduled to a specific node. Now let us create a `Pod` that refers to the above volume. When the `Pod` is created, the volume will be attached, formatted and mounted to the specified container:

```
kind: Pod
apiVersion: v1
metadata:
  name: pod-nimble
spec:
  containers:
    - name: pod-nimble-con-1
      image: nginx
      command: ["bin/sh"]
      args: ["-c", "while true; do date >> /data/mydata.txt; sleep 1; done"]
      volumeMounts:
        - name: export1
          mountPath: /data
    - name: pod-nimble-cont-2
      image: debian
      command: ["bin/sh"]
      args: ["-c", "while true; do date >> /data/mydata.txt; sleep 1; done"]
      volumeMounts:
        - name: export1
          mountPath: /data
  volumes:
    - name: export1
      persistentVolumeClaim:
        claimName: pvc-nimble
```

Check if the pod is running successfully:

```
$ kubectl get pod pod-nimble
NAME         READY   STATUS    RESTARTS   AGE
pod-nimble   2/2     Running   0          2m29s
```

## Use case specific examples
This `StorageClass` examples help guide combinations of options when provisioning volumes.

## Data protection
This `StorageClass` creates thinly provisioned volumes with deduplication turned on. It will also apply the Performance Policy "SQL Server" along with a Protection Template. The Protection Template needs to be defined on the array.

```
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

## Clone and throttle for devs
This `StorageClass` will create clones of a "production" volume and throttle the performance of each clone to 1000 IOPS. When the PVC is deleted, it will be permanently deleted from the backend array.

```
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

## Clone a non-containerized volume
This `StorageClass` will clone a standard backend volume (without container metadata on it) from a particular pool on the backend.

```
kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  name: import-clone-legacy-prod
rovisioner: hpe.com/nimble
parameters:
  pool: "flash"
  importVolAsClone: "production-db-vol"
  destroyOnRm: "true"
```

## Import (cutover) a volume
This `StorageClass` will import an existing Nimble volume to Kubernetes. The source volume needs to be offline for the import to succeed.

```
kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  name: import-clone-legacy-prod
provisioner: hpe.com/nimble
parameters:
  pool: "flash"
  importVol: "production-db-vol"
```

## Using overrides
The HPE Dynamic Provisioner for Kubernetes understands a set of annotation keys a user can set on a `PVC`. If the corresponding keys exists in the list of the `allowOverrides` key in the `StorageClass`, the end-user can tweak certain aspects of the provisioning workflow. This opens up for very advanced data services.

### StorageClass
```
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: my-sc
provisioner: hpe.com/nimble
parameters:
  description: "Volume provisioned by StorageClass my-sc"
  dedupe: "false"
  destroyOnRm: "true"
  perfPolicy: "Windows File Server"
  folder: "myfolder"
  allowOverrides: snapshot,limitIOPS,perfPolicy
 ```

### PersistentVolumeClaim
```
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
 name: my-pvc
 annotations:
    hpe.com/description: "This is my custom description"
    hpe.com/limitIOPS: "8000"
    hpe.com/perfPolicy: "SQL Server"
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
  storageClassName: my-sc
```
This will create a `PV` of 8000 IOPS with the Performance Policy of "SQL Server" and a custom volume description.

## Creating clones of PVCs
Using a `StorageClass` to clone a `PV` is practical when there's needs to clone across namespaces (for example from prod to test or stage). If a user wants to clone any arbitrary volume, it becomes a bit tedious to create a `StorageClass` for each clone. The annotation `hpe.com/CloneOfPVC` allows a user to clone any `PVC` within a namespace.

```
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
 name: my-pvc-clone
 annotations:
    hpe.com/cloneOfPVC: my-pvc
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
  storageClassName: my-sc
```
