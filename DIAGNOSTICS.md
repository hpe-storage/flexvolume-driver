# Diagnostics
This document outlines a few troubleshooting steps for the HPE Volume Driver for Kubernetes Plugin. This product is supported by HPE, please consult with your support organization (Nimble, Cloud Volumes etc) prior attempting any configuration changes.

## Troubleshooting FlexVolume driver
The FlexVolume driver is a binary executed by the kubelet to perform mount/unmount/attach/detach operations as workloads request storage resources. The binary relies on communicating with a socket on the host where the volume plugin responsible for the MUAD operations perform control-plane or data-plane operations against the backend system hosting the actual volumes.

### Locations
The driver has a configuration file where certain defaults can be tweaked to accommodate a certain behavior. Under normal circumstances, this file does not need any tweaking.

The name and the location of the binary varies based on Kubernetes distribution (the default 'exec' path) and what backend driver is being used. In a typical scenario, using Nimble, this is expected:

Binary: `/usr/libexec/kubernetes/kubelet-plugins/volume/exec/hpe.com~nimble/nimble`
Config file: `/usr/libexec/kubernetes/kubelet-plugins/volume/exec/hpe.com~nimble/nimble.json`

### Override defaults
By default, it contains only the path to the socket file for the volume plugin:
```json
{
    "dockerVolumePluginSocketPath": "/etc/hpe-storage/nimble.sock"
}
```

Valid options for the FlexVolume driver can be inspected by executing the binary on the host with the `config` argument:
```
/usr/libexec/kubernetes/kubelet-plugins/volume/exec/hpe.com~nimble/nimble config
Error processing option 'logFilePath' - key:logFilePath not found
Error processing option 'logDebug' - key:logDebug not found
Error processing option 'supportsCapabilities' - key:supportsCapabilities not found
Error processing option 'stripK8sFromOptions' - key:stripK8sFromOptions not found
Error processing option 'createVolumes' - key:createVolumes not found
Error processing option 'listOfStorageResourceOptions' - key:listOfStorageResourceOptions not found
Error processing option 'factorForConversion' - key:factorForConversion not found
Error processing option 'enable1.6' - key:enable1.6 not found

Driver=nimble Version=v2.5.1-50fbff2aa14a693a9a18adafb834da33b9e7cc89
Current Config:
  dockerVolumePluginSocketPath = /etc/hpe-storage/nimble.sock
           stripK8sFromOptions = true
                   logFilePath = /var/log/dory.log
                      logDebug = false
                 createVolumes = false
                     enable1.6 = false
           factorForConversion = 1073741824
  listOfStorageResourceOptions = [size sizeInGiB]
          supportsCapabilities = true
```
An example tweak could be to enable debug logging and enable support for Kubernetes 1.6 (which we don't officially support). The config file would then end up like this:

```json
{
    "dockerVolumePluginSocketPath": "/etc/hpe-storage/nimble.sock",
    "logDebug": true,
    "enable1.6": true
}
```
Execute the binary again (`nimble config`) to ensure the parameters and config file gets parsed correctly. Since the config file is read on each FlexVolume operation, no restart of anything is needed.

See [ADVANCED.md](ADVANCED.md) for more parameters for the driver.json file.

### Connectivity
To verify the FlexVolume binary can actually communicate with the backend volume plugin, issue a faux mount request:

```
/usr/libexec/kubernetes/kubelet-plugins/volume/exec/hpe.com~nimble/nimble mount no/op '{"name":"myvol1"}'
```

If the FlexVolume driver can successfully communicate with the volume plugin socket:
```
{"status":"Failure","message":"configured to NOT create volumes"}
```

In the case of any other output, check if the backend volume plugin is alive with `curl`:
```
curl --unix-socket /etc/hpe-storage/nimble.sock -d '{}' http://localhost/VolumeDriver.Capabilities
```
It should output:
```json
{"capabilities":{"scope":"global"},"Err":""}
```

## FlexVolume and dynamic provisioner driver logs
Log files associated with the HPE Volume Driver for Kubernetes FlexVolume Plugin logs data to the standard output stream. If the logs need to be retained for long term, use a standard logging solution. Some of the logs on the host are persisted which follow standard logrotate policies.

* FlexVolume driver logs:
  `kubectl logs -f daemonset.apps/hpe-flexvolume-driver -n kube-system`
  The logs are persisted at `/var/log/hpe-docker-plugin.log` and `/var/log/dory.log`

* Dynamic Provisioner logs:
  `kubectl logs -f  deployment.apps/hpe-dynamic-provisioner`
  The logs are persisted at `/var/log/hpe-dynamic-provisioner.log`

## Log Collector
Log collector script `hpe-logcollector.sh` can be used to collect diagnostic logs from the hosts.

Depending on the severity of the problem, the script may need to executed directly on the compute node itself. Normal usage would entail running the collector script using `kubectl` in one of the driver `Pods` from the `DaemonSet` or the dynamic provisioner `Deployment`.

If the issue is related to mount, unmount, attach or detach:
```
kubectl get pods -n kube-system -l name=hpe-flexvolume-driver -o wide
NAME                          READY   STATUS    RESTARTS   AGE     IP              NODE              
hpe-flexvolume-driver-b4rdn   1/1     Running   0          2d22h   10.21.192.224   tme-lnx9-salsa-wrk
hpe-flexvolume-driver-ch9cr   1/1     Running   0          2d22h   10.21.192.225   tme-lnx7-salsa-wrk
hpe-flexvolume-driver-x282s   1/1     Running   0          2d22h   10.21.192.222   tme-lnx8-salsa-wrk
```
Pay attention the node name, the script should be executed on the node exhibiting the issue.

If the issue is related to provisioning or deleting Persistent Volumes:
```
kubectl get pods -n kube-system -l daemon=kube-storage-controller-daemon -o wide
NAME                                             READY   STATUS    RESTARTS   AGE     IP              NODE              
kube-storage-controller-doryd-7d87f75586-xpbl8   1/1     Running   0          2d22h   10.21.192.225   tme-lnx7-salsa-wrk
```
It's important the script is executed in a FlexVolume driver `Pod`, not the dynamic provisioner `Pod`. In the example above if there were issues with deleting volumes, the script should've been executed in the `Pod` `hpe-flexvolume-driver-ch9cr` which reside on the same node as the provisioner: "tme-lnx7-salsa-wrk".

Now, copy & paste the `NAME` into this `kubectl` command to make sure connectivity to the `Pod`:
```
kubectl exec -nkube-system -it hpe-flexvolume-driver-b4rdn -- hpe-logcollector.sh -h
Diagnostic LogCollector Script to collect HPE Storage logs

Usage: hpe-logcollector.sh [CASE_NUMBER]
       where CASE_NUMBER is an optional parameter <HPE Nimble Storage Support Case Number>
       needed to upload the logs to the HPE Nimble Storage FTP Server
```

Now execute, with "00000000" replaced with the actual "HPE Nimble Storage Support Case Number".
```
kubectl exec -nkube-system -it hpe-flexvolume-driver-b4rdn -- hpe-logcollector.sh 00000000
======================================
Collecting the Diagnostic Information
======================================
CaseNumber is 00000000
Transfer to the ftp server
FTP of /var/log/00000000-hpe-storage-logs-tme-lnx9-salsa-wrk-20190915_030454.tar.gz to upload completed
Complete
====================================
```
