# flexvolume-driver

FlexVolume Driver for Kubernetes leverages HPE Nimble Storage or HPE Cloud Volumes to provide scalable and persistent storage for stateful applications.

## Deploying to Kubernetes

### Step 1: Create a secret

#### For on-prem deployments, create a secret with your array details:

Replace the password string (`YWRtaW4=`) with a base64 encoded version of your password and replace the ip with your array IP address and save it as `secret.yaml`:

```
apiVersion: v1
kind: Secret
metadata:
  name: nimble-secret
  namespace: kube-system
stringData:
  ip: 10.0.0.1
  username: admin
  protocol: "iscsi"
data:
  # echo -n "admin" | base64
  password: YWRtaW4=
```

#### For HPE cloud volumes deployments, create a secret with your access details:

Replace the username and password strings (`YWRtaW4=`) with a base64 encoded version of your HPE Cloud Volumes access_key and access_secret. Also, replace the ip with HPE Cloud Volumes portal URL address and save it as `secret.yaml`:

```
apiVersion: v1
kind: Secret
metadata:
  name: nimble-secret
  namespace: kube-system
stringData:
  ip: cloudvolumes.hpe.com
  protocol: "iscsi"
data:
  # echo -n "$cv_access_key" | base64
  username: YWRtaW4=
  # echo -n "$cv_access_secret" | base64
  password: YWRtaW4=
```

#### Then create the secret using kubectl:

```
$ kubectl create -f secret.yaml
secret "nimble-secret" created
```

You should now see the nimble secret in the `kube-system` namespace along with other secrets

```
$ kubectl get secret -n kube-system
NAME                  TYPE                                  DATA      AGE
default-token-jxx2s   kubernetes.io/service-account-token   3         2d1h
nimble-secret         Opaque                                5         149m
```

### Step 2. Create a ConfigMap

#### For on-prem deployments, perform the following steps:

Edit the below default parameters as required for flexvolume driver and save it as `nimble-config.yaml`

```
---
kind: ConfigMap
apiVersion: v1
metadata:
  name: nimble-config
  namespace: kube-system
data:
  volume-driver.json: |-
    {
      "global":   {},
      "defaults": {
                 "sizeInGiB":"10",
                 "limitIOPS":"-1",
                 "limitMBPS":"-1",
                 "perfPolicy": "Other",
                 "mountConflictDelay": 120
                },
      "overrides":{}
    }
```

#### For HPE cloud volumes deployments, perform the following steps:

Edit the below parameters as required with your public cloud info and save it as `nimble-config.yaml`. Initiator IP addresses of all nodes in the cluster needs to be specified to successfully mount volume on any node. This will be handled automatically during release.

```
---
kind: ConfigMap
apiVersion: v1
metadata:
  name: nimble-config
  namespace: kube-system
data:
  volume-driver.json: |-
    {
    "global": {
                "nameSuffix": ".docker",
                "snapPrefix": "BaseFor",
                "initiators": ["10.1.0.1", "10.1.0.8", "10.1.0.7"],
                "automatedConnection": true,
                "existingCloudSubnet": "10.1.0.0/24",
                "region": "us-east-1",
                "privateCloud": "vpc-data",
                "cloudComputeProvider": "Google Cloud Platform"
    },
    "defaults": {
                "sizeInGiB": 10,
                "limitIOPS": 1000,
                "fsOwner": "0:0",
                "fsMode": "600",
                "description": "persistent volume for container",
                "perfPolicy": "Other",
                "protectionTemplate": "twicedaily:4",
                "encryption": true,
                "volumeType": "PF",
                "destroyOnRm": true,
                "mountConflictDelay": 120
    },
    "overrides": {
    }
}
```

#### Then create the ConfigMap using kubectl:

```
$ kubectl create -f nimble-config.yaml
configmap/nimble-config created
```

### Step 3. Deploy the HPE FlexVolume driver and HPE Dynamic Provisioner (doryd)

Deploy HPE FlexVolume driver as Daemonset and HPE Dynamic Provisioner(doryd) as Deployment using below command:

```
$ kubectl create -f hpe-flexvolume-driver.yaml
```

or for HPE Cloud Volumes

```
$ kubectl create -f cv-flexvolume-driver.yaml
```

Check to see all hpe-flexvolume and kube-storage-controller-doryd pods are running using kubectl:

```
$ kubectl get pods -n kube-system
NAME                                            READY   STATUS    RESTARTS   AGE
alertmanager-0                                  2/2     Running   0          3d
alertmanager-1                                  2/2     Running   10         14d
calico-node-2mhb2                               2/2     Running   3          14d
calico-node-dlwnq                               2/2     Running   3          14d
calico-node-zq6bw                               2/2     Running   4          14d
calico-typha-7dc5d97dbd-4v88x                   1/1     Running   0          3d
citadel-565dc89cb8-xd46t                        1/1     Running   0          3d
config-management-operator-d764468c9-dsjx7      1/1     Running   7          3d
grafana-0                                       2/2     Running   0          3d
hpe-flexvolume-driver-2rdt4                     1/1     Running   0          3d
hpe-flexvolume-driver-md562                     1/1     Running   0          3d
hpe-flexvolume-driver-x4k96                     1/1     Running   0          3d
k8s-bigip-ctlr-deployment-758fcbccdc-457sh      1/1     Running   0          3d
kube-dns-54f6d9b699-bwmrj                       3/3     Running   0          3d
kube-dns-54f6d9b699-z2lwf                       3/3     Running   13         14d
kube-dns-autoscaler-74b69ddf68-js2dr            1/1     Running   0          3d
kube-proxy-gfvw6                                1/1     Running   1          14d
kube-proxy-j6kwj                                1/1     Running   1          14d
kube-proxy-kttss                                1/1     Running   1          14d
kube-state-metrics-79c995cd8c-kq95s             2/2     Running   0          3d
kube-storage-controller-doryd-6db794764-l4rj8   1/1     Running   2          12d
load-balancer-f5-589b5b6d8d-7b62x               1/1     Running   0          3d
logging-operator-cbbc69c67-482km                1/1     Running   0          3d
node-exporter-9md2k                             2/2     Running   2          14d
node-exporter-f8klp                             2/2     Running   2          14d
node-exporter-mlcwh                             2/2     Running   2          14d
prometheus-0                                    2/2     Running   0          3d
prometheus-1                                    2/2     Running   2          14d
pushprox-client-78cc6d85cb-68jd4                2/2     Running   0          3d
stackdriver-operator-7b4bb8f667-ll56b           1/1     Running   0          3d
```

### Step 4. Test and verify volume provisioning

Create a StorageClass with volume parameters as required.  Note the example below will look for a protection template named "cloud-repl-template" that is configured to replicate as an HPE Cloud Volume.

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
  protectionTemplate: "cloud-repl-template"
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

The above output means that the HPE Nimble Storage FlexVolume driver successfully provisioned a new volume and automatically setup to replicate to HPE Cloud Volumes.  The volume is not attached to any node yet. It will only be attached to a node if a workload is scheduled to a specific node. Now let us create a Pod that refers to the above volume. When the Pod is created, the volume will be attached, formatted and mounted to the specified container:

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

### Step 5. Test and verify volume provisioning for HPE cloud volumes using import

Create a StorageClass with volume parameters as required and specify the replica volume name to import as clone along with its replication store name as below example. Note the provisioner name as `hpe.com/cv`.

```
---
kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
 name: sc-cv
provisioner: hpe.com/cv
parameters:
  description: "Nimble Storage Class"
  mountConflictDelay: "150"
  importVolAsClone: "base-replica-volume"
  replStore: "replication-store-name"
```

Create a PersistentVolumeClaim. This makes sure a volume is created and provisioned on your behalf:

```
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: pvc-cv
spec:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
  storageClassName: sc-cv
```

Check that a new `PersistentVolume` is created based on your claim:

```
$ kubectl get pv
NAME                                            CAPACITY     ACCESS MODES   RECLAIM POLICY   STATUS   CLAIM               STORAGECLASS   REASON   AGE
sc-cv-13336da3-7ca3-11e9-826c-00505693581f   10Gi         RWO            Delete           Bound    default/pvc-cv   sc-cv               3s
```

The above output means that the HPE Cloud Volumes FlexVolume driver successfully provisioned a new volume cloned from replicated on-prem volume. The volume is not attached to any node yet. It will only be attached to a node if a workload is scheduled to a specific node. Now follow the steps above to create a pod referencing the HPE Cloud Volume created above.

## Sample storage classes

Sample storage classes can be found for [Nimble Storage](examples/kubernetes/nimble-storage/README.md), [Cloud Volumes](examples/kubernetes/cloud-volumes/README.md), Simplivity, and 3Par.

## FlexVolume-driver config options

Following options are supported by flexvolume-driver which can be provided under flexvolume exec path with convention {driver-name}.json

{
    "logFilePath": "/var/log/hpe-flexvolume-driver.log",
    "logDebug": false,
    "stripK8sFromOptions": true,
    "dockerVolumePluginSocketPath": "/etc/hpe-storage/nimble.sock",
    "createVolumes": false,
    "enable1.6": false,
    "listOfStorageResourceOptions" :    ["size","sizeInGiB"],
    "factorForConversion": 1073741824,
    "defaultOptions": [{"manager": "k8s"}]
}

## FAQ

Some of the troubleshooting and FAQs can be found at [FAQ](examples/kubernetes/nimble-storage/FAQ.md)
