# FAQs

## 1. Troubleshooting Flexvolume Driver

### a. Config File

`nimble.json` config file can be used to override some of the default options.
By default the following options are implemented and can be overridden by the json file at

`/usr/libexec/kubernetes/kubelet-plugins/volume/exec/hpe.com~nimble/nimble.json`

```json
{
    "logFilePath": "/var/log/dory.log",
    "logDebug": false,
    "dockerVolumePluginSocketPath": "/etc/hpe-storage/nimble.sock",
    "defaultOptions": [{"option1": "value1"}, {"option2": "value2"}]
}
```

where `option1` and `option2` are valid flexvolume create options

### b. To ensure FlexVolume to Docker Volume Driver connectivity, manually test a mount using

```markdown
/usr/libexec/kubernetes/kubelet-plugins/volume/exec/hpe.com~nimble/nimble mount /tmp/1
'{"name":"testvol", "sizeInGiB":"20", "destroyOnRm": "true"}'
```

This should result in the following error:

```json
{"status":"Failure","message":"unable to split /tmp/1"}.
```

### c. Make sure dory can communicate with dockerplugin socket /run/docker/plugins/nimble.sock

```markdown
Info : 2019/03/14 18:21:52 dory.go:82: [14414] request: init []
Info : 2019/03/14 18:21:52 dory.go:100: [14414] reply  : init []: {"status":"Success","capabilities":{"attach":false}}
```
