# HPE Volume Driver for Kubernetes FlexVolume Plugin

HPE Volume Driver for Kubernetes FlexVolume Plugin leverages HPE Nimble Storage or HPE Cloud Volumes to provide scalable and persistent storage for stateful applications.

## Host Platform Requirements

* OpenShift Container Platform 3.9, 3.10 and 3.11.
* Kubernetes 1.10 and above.
* Redhat/CentOS 7.5+
* Ubuntu 16.04/18.04 LTS

## Storage Platform Requirements

### HPE Nimble Storage

| Driver      | HPE Nimble Storage Version | Release Notes    |
|-------------|----------------------------|------------------|
| v3.0.0      | 5.0.8.x and 5.1.3.x onwards          | [v3.0.0](release-notes/v3.0.0.md)|

**Note:** Synchronous replication (Peer Persistence) is not supported by the HPE Volume Driver for Kubernetes FlexVolume Plugin.

### HPE Cloud Volumes

* Amazon EKS 1.12/1.13
* Microsoft Azure AKS 1.12/1.13
* US regions only

## Deploying to Kubernetes
The recommend way to deploy and manage the HPE Volume Driver for Kubernetes FlexVolume Plugin is to use Helm. Please see the [co-deployments](https://github.com/hpe-storage/co-deployments) repository for further information.

Use the following steps for a manual installation.

### Step 1: Create a secret

#### HPE Nimble Storage
Replace the `password` string (`YWRtaW4=`) below with a base64 encoded version of your password and replace the `backend` with your array IP address and save it as `hpe-secret.yaml`.

```
apiVersion: v1
kind: Secret
metadata:
  name: hpe-secret
  namespace: kube-system
stringData:
  backend: 192.168.1.1
  username: admin
  protocol: "iscsi"
data:
  # echo -n "admin" | base64
  password: YWRtaW4=
```

#### HPE Cloud Volumes
Replace the `username` and `password` strings (`YWRtaW4=`) with a base64 encoded version of your HPE Cloud Volumes "access_key" and "access_secret". Also, replace the `backend` with HPE Cloud Volumes portal fully qualified domain name (FQDN) and save it as `hpe-secret.yaml`.

```
apiVersion: v1
kind: Secret
metadata:
  name: hpe-secret
  namespace: kube-system
stringData:
  backend: cloudvolumes.hpe.com
  protocol: "iscsi"
  serviceName: cv-cp-svc
  servicePort: "8080"
data:
  # echo -n "<my very confidential access key>" | base64
  username: YWRtaW4=
  # echo -n "<my very confidential secret key>" | base64
  password: YWRtaW4=
```

#### Create the secret

```
$ kubectl create -f hpe-secret.yaml
secret "hpe-secret" created
```

You should now see the HPE secret in the `kube-system` namespace.

```
$ kubectl get secret/hpe-secret -n kube-system
NAME                  TYPE                                  DATA      AGE
hpe-secret            Opaque                                5         3s
```

### Step 2. Create a ConfigMap
The `ConfigMap` is used to set and tweak defaults for both the FlexVolume driver and Dynamic Provisioner.

#### HPE Nimble Storage
Edit the below default parameters as required for FlexVolume driver and save it as `hpe-config.yaml`.

```
kind: ConfigMap
apiVersion: v1
metadata:
  name: hpe-config
  namespace: kube-system
data:
  volume-driver.json: |-
    {
      "global":   {},
      "defaults": {
                 "limitIOPS":"-1",
                 "limitMBPS":"-1",
                 "perfPolicy": "Other"
                },
      "overrides":{}
    }
```

#### HPE Cloud Volumes
Edit the below parameters as required with your public cloud info and save it as `hpe-config.yaml`.

**Note**: Initiator IP addresses of all nodes in the cluster needs to be specified to successfully mount volume on any node. This will be handled automatically for the v3.1.0 release.

```
kind: ConfigMap
apiVersion: v1
metadata:
  name: hpe-config
  namespace: kube-system
data:
  volume-driver.json: |-
    {
      "global": {
                "snapPrefix": "BaseFor",
                "initiators": ["eth0"],
                "automatedConnection": true,
                "existingCloudSubnet": "10.1.0.0/24",
                "region": "us-east-1",
                "privateCloud": "vpc-data",
                "cloudComputeProvider": "Amazon AWS"
      },
      "defaults": {
                "limitIOPS": 1000,
                "fsOwner": "0:0",
                "fsMode": "600",
                "description": "Volume provisioned by the HPE Volume Driver for Kubernetes FlexVolume Plugin",
                "perfPolicy": "Other",
                "protectionTemplate": "twicedaily:4",
                "encryption": true,
                "volumeType": "PF",
                "destroyOnRm": true
      },
      "overrides": {
      }
    }
```

#### Create the ConfigMap

```
$ kubectl create -f hpe-config.yaml
configmap/hpe-config created
```

Please see [ADVANCED.md](ADVANCED.md) for more `volume-driver.json` configuration options.

### Step 3. Deploy the HPE Volume Driver for Kubernetes FlexVolume Plugin and HPE Dynamic Provisioner for Kubernetes
Deploy the driver as a `DaemonSet` and the dynamic provisioner as a `Deployment`.

#### HPE Nimble Storage

```
$ kubectl create -f https://raw.githubusercontent.com/hpe-storage/co-deployments/master/yaml/flexvolume-driver/hpe-flexvolume-driver-v3.0.0.yaml
```

#### HPE Cloud Volumes

```
$ kubectl create -f https://raw.githubusercontent.com/hpe-storage/co-deployments/master/yaml/flexvolume-driver/hpe-cloud-volumes/hpecv-cp.yaml

Amazon EKS:
$ kubectl create -f https://raw.githubusercontent.com/hpe-storage/co-deployments/master/yaml/flexvolume-driver/hpe-cloud-volumes/hpecv-aws-flexvolume-driver.yaml

Microsoft Azure AKS:
$ kubectl create -f https://raw.githubusercontent.com/hpe-storage/co-deployments/master/yaml/flexvolume-driver/hpe-cloud-volumes/hpecv-azure-flexvolume-driver.yaml

Generic:
$ kubectl create -f https://raw.githubusercontent.com/hpe-storage/co-deployments/master/yaml/flexvolume-driver/hpe-cloud-volumes/hpecv-flexvolume-driver.yaml
```

**Note:** The declarations for HPE Volume Driver for Kubernetes FlexVolume Plugin can be found [here](https://github.com/hpe-storage/co-deployments/tree/master/yaml/flexvolume-driver)

Check to see all `hpe-flexvolume-driver` `Pods` (one per compute node) and the `hpe-dynamic-provisioner` Pod are running:
```
$ kubectl get pods -n kube-system
NAME                                            READY   STATUS    RESTARTS   AGE
hpe-flexvolume-driver-2rdt4                     1/1     Running   0          45s
hpe-flexvolume-driver-md562                     1/1     Running   0          44s
hpe-flexvolume-driver-x4k96                     1/1     Running   0          44s
hpe-dynamic-provisioner-59f9d495d4-hxh29        1/1     Running   0          24s
```

## Using
Get started using the FlexVolume driver by setting up `StorageClass`, `PVC` API objects. See [USING.md](USING.md) for examples.

## Building
Instructions on how to build the FlexVolume driver from sources can be found in [BUILDING.md](BUILDING.md)

## Diagnostics
Logging and other troubleshooting steps can be in [DIAGNOSTICS.md](DIAGNOSTICS.md)

## Support
The HPE Volume Driver for Kubernetes FlexVolume Plugin is supported software by Hewlett Packard Enterprise. Reach out to your HPE representation to be connected with the support organization with any general issue you need help resolving.

We also encourage open collaboration, file issues, questions or feature requests [here](https://github.com/hpe-storage/flexvolume-driver/issues). You may also join our Slack community to chat with HPE folks close to this project. We hang out in `#NimbleStorage` and `#Kubernetes` at [slack.hpedev.io](https://slack.hpedev.io/).

## Contributing
We value all feedback and contributions. If you find any issues or want to contribute, please feel free to open an issue or file a PR. More details in [CONTRIBUTING.md](CONTRIBUTING.md)

## License
This is open source software licensed using the Apache License 2.0. Please see [LICENSE](LICENSE) for details.
