## Using the HPE Volume Driver for Kubernetes FlexVolume Plugin

These instructions are provided as an example on how to use the HPE Volume Driver for Kubernetes FlexVolume Plugin with the HPE Nimble Storage Array.

### Test and verify volume provisioning
The below YAML declarations are meant to be created with `kubectl create`. Either copy the content to a file on the host where `kubectl` is being executed, or copy & paste into the terminal, like this:

```
kubectl create -f-
< paste the YAML >
^D (CTRL + D)
```

**Note:**  Some of the examples supported by the HPE Volume Driver for Kubernetes FlexVolume Plugin are available in the [examples/kubernetes/hpe-nimble-storage](examples/kubernetes/hpe-nimble-storage) directory or [examples/kubernetes/cloud-volumes](examples/kubernetes/cloud-volumes) and all the HPE Nimble Storage Array Flexvolume `StorageClass` parameters can be found in [examples/kubernetes/hpe-nimble-storage](examples/kubernetes/hpe-nimble-storage).

To get started, create a `StorageClass` API object referencing the `nimble-secret` and defining additional (optional) `StorageClass` parameters:

### Step 1. Sample storage classes

Sample storage classes can be found for [Nimble Storage](examples/kubernetes/hpe-nimble-storage/storage-class.yaml), [Cloud Volumes](examples/kubernetes/cloud-volumes/README.md), Simplivity, and 3Par.

### Step 2. Test and verify volume provisioning

Create a StorageClass with volume parameters as required.

```yaml
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

```yaml
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

The above output means that the HPE Nimble Storage FlexVolume driver successfully provisioned a new volume and automatically. The volume is not attached to any node yet. It will only be attached to a node if a workload is scheduled to a specific node. Now let us create a Pod that refers to the above volume. When the Pod is created, the volume will be attached, formatted and mounted to the specified container:

```yaml
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
