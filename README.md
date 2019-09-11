# HPE Volume Driver for Kubernetes FlexVolume Plugin

HPE Volume Driver for Kubernetes FlexVolume Plugin leverages HPE Nimble Storage or HPE Cloud Volumes to provide scalable and persistent storage for stateful applications.

## Requirements

* OpenShift Container Platform 3.9, 3.10 and 3.11.
* Kubernetes 1.10 and above.
* Redhat / CentOS 7.5+
* Ubuntu LTS 16.04 / 18.04

## Support Matrix for HPE Volume Driver for Kubernetes FlexVolume Plugin on HPE Nimble Storage

| Release                 | HPE Nimble Storage Version    |
|-------------------------|----------|
| v3.0.0              | 5.1.3.x |

**Note:** Synchronous Replication (Peer persistence) is not supported by the HPE Volume Driver for Kubernetes FlexVolume Plugin

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

### Step 3. Deploy the HPE Volume Driver for Kubernetes FlexVolume Plugin and HPE Dynamic Provisioner for Kubernetes

Deploy HPE FlexVolume driver as Daemonset and HPE Dynamic Provisioner as Deployment using below command:

```
$ kubectl create -f https://raw.githubusercontent.com/hpe-storage/co-deployments/master/yaml/flexvolume-driver/hpe-flexvolume-driver-v3.0.0.yaml
```

or for HPE Cloud Volumes

```
$ kubectl create -f https://raw.githubusercontent.com/hpe-storage/co-deployments/master/yaml/flexvolume-driver/hpecv-flexvolume-driver.yaml
```

**Note:** The deployment yaml files for HPE Volume Driver for Kubernetes FlexVolume can be found [here](https://github.com/hpe-storage/co-deployments/tree/master/yaml/flexvolume-driver)

Check to see all hpe-flexvolume and hpe-dynamic-provisioner pods are running using kubectl:

```markdown

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
pod/hpe-dynamic-provisioner-59f9d495d4-hxh29    1/1     Running   2          12d
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

## Building the HPE Volume Driver for Kubernetes FlexVolume Plugin

Instructions on how to build the HPE Volume Driver for Kubernetes FlexVolume Plugin can be found in [BUILDING.md](BUILDING.md)

## Using the HPE Volume Driver for Kubernetes FlexVolume Plugin

Getting started with the HPE Volume Driver for Kubernetes FlexVolume Plugin, setting up `StorageClass`. See [USING.md](USING.md) for examples.

## FlexVolume-driver config options

Following options are supported by flexvolume-driver which can be provided under flexvolume exec path with convention {driver-name}.json

```json

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
```

## Logging and Diagnostic

Log files associated with the HPE Volume Driver for Kubernetes FlexVolume Plugin logs data to the standard output stream. If the logs need to be retained for long term, use a standard logging solution. Some of the logs on the host are persisted which follow standard logrotate policies.

### FlexVolume Plugin Logs

* Flexvolume plugin logs:
  `kubectl logs -f daemonset.apps/hpe-flexvolume-driver -n kube-system`
  The logs are persisted at `/var/log/hpe-docker-plugin.log` and `/var/log/dory.log`

* HPE Dynamic Provisioner logs:
  `kubectl logs -f  deployment.apps/hpe-dynamic-provisioner`
  The logs are persisted at `/var/log/hpe-dynamic-provisioner.log`

### Log Collector

Log collector script `hpe-logcollector.sh` can be used to collect diagnostic logs from the hosts.

```markdown

hpe-logcollector.sh -h
Diagnostic LogCollector Script to collect HPE Storage logs
```

## Support

Please file any issues, questions or feature requests [here](https://github.com/hpe-storage/flexvolume-driver/issues). You may also join our Slack community to chat with HPE folks close to this project. We hang out in `#NimbleStorage` and `#Kubernetes` at [slack.hpedev.io](https://slack.hpedev.io/).

## Contributing

We value all feedback and contributions. If you find any issues or want to contribute, please feel free to open an issue or file a PR. More details in [CONTRIBUTING.md](CONTRIBUTING.md)

## License

This is open source software licensed using the Apache License 2.0. Please see [LICENSE](LICENSE) for details.

## FAQ

Some of the troubleshooting and FAQs can be found at [FAQ](examples/kubernetes/hpe-nimble-storage/FAQ.md)
