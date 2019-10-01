/*
(c) Copyright 2017 Hewlett Packard Enterprise Development LP

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package flexvol

import (
	"encoding/json"
	"fmt"
	"github.com/hpe-storage/common-host-libs/docker/dockervol"
	"github.com/hpe-storage/common-host-libs/linux"
	log "github.com/hpe-storage/common-host-libs/logger"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	// InitCommand  - Initializes the driver.
	InitCommand = "init"
	// AttachCommand - Attach the volume specified by the given spec.
	AttachCommand = "attach"
	//DetachCommand - Detach the volume from the kubelet.
	DetachCommand = "detach"
	//MountCommand - Mount device mounts the device to a global path which individual pods can then bind mount.
	MountCommand = "mount"
	//UnmountCommand - Unmounts the filesystem for the device.
	UnmountCommand = "unmount"
	//GetVolumeNameCommand - Get the name of the volume.
	GetVolumeNameCommand = "getvolumename"
	//SuccessStatus indicates success
	SuccessStatus = "Success"
	//FailureStatus indicates failure
	FailureStatus = "Failure"
	//NotSupportedStatus indicates not supported
	NotSupportedStatus = "Not supported"
	//FailureJSON is a pre-marshalled response used in the case of a marshalling error
	FailureJSON = "{\"status\":\"Failure\",\"message\":\"Unknown error.\"}"
	//mountPathRegex describes the uuid and flexvolume name in the path
	//examples:
	// /var/lib/origin/openshift.local.volumes/pods/88917cdb-514d-11e7-93fb-5254005e615a/volumes/hpe~nimble/test2
	// /var/lib/kubelet/pods/fb36bec9-51f7-11e7-8eb8-005056968cbc/volumes/hpe~nimble/test
	mountPathRegex = "/var/lib/.*/pods/(?P<uuid>[\\w\\d-]*)/volumes/"

	//docker volume status key
	devicePathKey  = "devicePath"
	maxTries       = 3
	notMounted     = "not mounted"
	noFileOrDirErr = "no such file or directory"
)

var (
	//createVolumes indicate whether the driver should create missing volumes. We should disable this by default as the volumes need to have filesystem on it else mount will fail.
	createVolumes = false

	execPath string

	dvp *dockervol.DockerVolumePlugin
)

// Response containers the required information for each invocation
type Response struct {
	//"status": "<Success/Failure/Not Supported>",
	Status string `json:"status"`
	//"message": "<Reason for success/failure>",
	Message string `json:"message,omitempty"`
	//"device": "<Path to the device attached. This field is valid only for attach calls>"
	Device string `json:"device,omitempty"`
	//"volumeName:" "undocumented"
	VolumeName string `json:"volumeName,omitempty"`
	//"attached": <True/False (Return true if volume is attached on the node. Valid only for isattached call-out)>
	Attached bool `json:"attached,omitempty"`
	//Capabilities reported on Driver init
	DriverCapabilities map[string]bool `json:"capabilities,omitempty"`
}

//AttachRequest is used to create a volume if one with this name doesn't exist
type AttachRequest struct {
	Name           string
	PvOrVolumeName string `json:"kubernetes.io/pvOrVolumeName,omitempty"`
	FsType         string `json:"kubernetes.io/fsType,omitempty"`
	ReadWrite      string `json:"kubernetes.io/readwrite,omitempty"`
}

func (ar *AttachRequest) getBestName() string {
	if ar.Name != "" {
		return ar.Name
	}
	return ar.PvOrVolumeName
}

// Config controls the docker behavior
func Config(ePath string, options *dockervol.Options) (err error) {
	dvp, err = dockervol.NewDockerVolumePlugin(options)
	createVolumes = options.CreateVolumes
	execPath = ePath
	return err
}

// BuildJSONResponse marshals a message into the FlexVolume JSON Response.
// If error is not nil, the default Failure message is returned.
func BuildJSONResponse(response *Response) string {
	if len(response.Status) < 1 {
		response.Status = NotSupportedStatus
	}

	jmess, err := json.Marshal(response)
	if err != nil {
		return FailureJSON
	}
	return string(jmess)
}

// ErrorResponse creates a Response with Status and Message set.
func ErrorResponse(err error) *Response {
	response := &Response{
		Status: FailureStatus,
	}
	response.Message = err.Error()
	return response
}

//Get a volume (create if necessary) This was added to k8s 1.6
func Get(jsonRequest string) (string, error) {
	log.Infof("get called with (%s)\n", jsonRequest)
	req := &AttachRequest{}
	err := json.Unmarshal([]byte(jsonRequest), req)
	if err != nil {
		return "", err
	}
	name, err := getOrCreate(req.getBestName(), jsonRequest)
	if err != nil {
		return "", err
	}
	response := &Response{
		Status:     SuccessStatus,
		VolumeName: name,
	}
	return BuildJSONResponse(response), nil
}

//Attach doesn't attach a volume.  It simply creates a volume if necessary.  It then returns "Not Supported".
//This worked well in 1.5 in that it broke the create and mount into 2 timeout windows, but
//this has changed in 1.6.
func Attach(jsonRequest string) (string, error) {
	log.Tracef("attach called with %s\n", jsonRequest)
	req := &AttachRequest{}
	err := json.Unmarshal([]byte(jsonRequest), req)
	if err != nil {
		return "", err
	}

	_, err = getOrCreate(req.getBestName(), jsonRequest)
	if err != nil {
		return "", err
	}

	return BuildJSONResponse(&Response{Status: NotSupportedStatus, Message: "Not supported."}), nil
}

func getOrCreate(name, jsonRequest string) (string, error) {
	log.Tracef("getOrCreate called with %s and %s\n", name, jsonRequest)
	volume, err := getVolume(name)
	if err != nil || volume.Volume.Name != name {
		if !createVolumes {
			return "", fmt.Errorf("configured to NOT create volumes")
		}

		log.Infof("volume %s was not found(err=%v), creating a new volume using %v", name, err, jsonRequest)
		var options map[string]interface{}
		err := json.Unmarshal([]byte(jsonRequest), &options)
		if err != nil {
			log.Errorf("unable to unmarshal options for %v - %s", jsonRequest, err.Error())
			return "", err
		}
		newName, err := dvp.Create(name, options)
		log.Tracef("getOrCreate returning %v for %s", newName, name)
		if err != nil {
			return "", err
		}
		return newName, nil
	}
	return volume.Volume.Name, nil
}

// wrapper for dvp.Get() with retries incorporated
func getVolume(name string) (volume *dockervol.GetResponse, err error) {
	log.Tracef("getVolume called with %s", name)
	try := 0
	for {
		log.Tracef("dvp.Get() called with %s try:%d", name, try+1)
		volume, err = dvp.Get(name)
		log.Tracef("volume returned from dvp.Get() is %#v", volume)
		if volume != nil {
			return volume, nil
		}
		if err != nil {
			if try < maxTries {
				try++
				time.Sleep(time.Duration(try) * time.Second)
				continue
			}
			return nil, err
		}
		return volume, nil
	}
}

//Mount a volume
func Mount(args []string) (string, error) {
	log.Tracef("mount called with %v\n", args)
	err := ensureArg("mount", args, 2)
	if err != nil {
		return "", err
	}

	req := &AttachRequest{}
	//json seems to be in the second or third argument
	jsonRequest, err := findJSON(args, req)
	if err != nil {
		return "", err
	}

	dockerVolName := req.getBestName()
	_, err = getOrCreate(dockerVolName, jsonRequest)
	if err != nil {
		return "", err
	}

	mountID, err := getMountID(args[0])
	if err != nil && !(strings.Contains(err.Error(), "object was not found")) {
		return "", err
	}

	path, err := dvp.Mount(dockerVolName, mountID)
	if err != nil {
		return "", err
	}

	//Mkdir
	err = os.MkdirAll(args[0], 0755)
	if err != nil {
		return "", err

	}
	log.Tracef("flexpath %s, mountpath %s", args[0], path)

	err = doMount(args[0], path, dockerVolName, mountID)
	if err != nil {
		return "", err
	}

	return BuildJSONResponse(&Response{Status: SuccessStatus}), nil
}

// Unmount a volume
//nolint :gocyclo
func Unmount(args []string) (string, error) {
	log.Tracef("Unmount called with %v", args)
	mountID, err := getMountID(args[0])
	if err != nil {
		log.Errorf("Failed to get the mount id from flexvolpath %s ", args[0])
		return "", err
	}
	err = doUnMount(args[0], mountID)
	if err != nil {
		log.Errorf("Failed to unmount, flexvolpath %s ", args[0])
		return "", err
	}
	return BuildJSONResponse(&Response{Status: SuccessStatus}), nil
}

// retry getVolumeNameFromMountPath for maxTries
func retryGetVolumeNameFromMountPath(k8sPath, dockerPath string) (string, error) {
	log.Tracef("retryGetVolumeNameFromMountPath called with %s %s", k8sPath, dockerPath)
	try := 0
	for {
		log.Tracef("getVolumeNameFromMountPath called with %s %s try:%d", k8sPath, dockerPath, try+1)
		dockerVolumeName, err := getVolumeNameFromMountPath(k8sPath, dockerPath)
		if err != nil {
			if try < maxTries {
				try++
				time.Sleep(time.Duration(try) * time.Second)
				continue
			}
			return "", err
		}
		log.Tracef("dockerVolumeName %s found at k8sPath :%s", dockerVolumeName, k8sPath)
		return dockerVolumeName, nil
	}
}

func getMountID(path string) (string, error) {
	log.Tracef("getMountID called with %v\n", path)
	// Replace backward slash to forward slash.
	backSlash := "\\\\"
	cregx := regexp.MustCompile(backSlash)
	mpath := cregx.ReplaceAllString(path, "/")

	r := regexp.MustCompile(mountPathRegex)
	groups := r.FindStringSubmatch(mpath)
	if len(groups) < 2 {
		return "", fmt.Errorf("unable to extract uuid from path %s", mpath)
	}
	log.Tracef("getMountID returning \"%s\"", groups[1])
	return groups[1], nil

}

//nolint: gocyclo
func getVolumeNameFromMountPath(k8sPath, dockerPath string) (string, error) {
	log.Tracef("getVolumeNameFromMountPath called with %s and %s", k8sPath, dockerPath)
	// sometimes the dockerPath is empty in case of failover/failback scenarios for OSP 3.11 and greater make sure we return the volume if it exists mounted
	if dockerPath == "" && k8sPath != "" {
		//if docker path is empty but k8sPath exist, try to use that to the unmount try to use k8s path for volume name
		volNames := strings.Split(k8sPath, "/")
		if len(volNames) == 0 || volNames[len(volNames)-1] == "" {
			return "", fmt.Errorf("no volume found from k8s path %s", k8sPath)
		}
		return volNames[len(volNames)-1], nil
	}
	name := filepath.Base(dockerPath)
	dockerVolume, err := getVolume(name)
	log.Tracef("retrieved dockerVolume %#v with name %s", dockerVolume, name)
	if err != nil || dockerVolume.Volume.Name != name {
		// The docker plugin might not use the docker volume name in the path.
		// Therefore we need to look through the know volumes to find out who
		// is mounted at that path.
		volumes, err2 := dvp.List()
		if err2 != nil {
			log.Errorf("Unable to get list of volumes. - %s, get error was %s", err2, err)
			return "", err
		}
		for _, vol := range volumes.Volumes {
			if vol.Mountpoint == dockerPath {
				return vol.Name, nil
			}
		}
		return "", fmt.Errorf("unable to find docker volume for %s.  No docker volume claimed to be mounted at %s", k8sPath, dockerPath)
	}
	if dockerVolume.Volume.Mountpoint == "" {
		// it could be possible that the information on the container provider was removed as it was mounted by other host, so don't treat it as an error
		log.Tracef("found a docker volume but its MountPoint was \"\", checking from /proc/mounts")
		devPath, _ := linux.GetDeviceFromMountPoint(dockerPath)
		log.Tracef("devPath %s was found for volume %s since MountPoint was \"\" in dockerVolume status", devPath, dockerVolume.Volume.Name)
		if devPath != "" {
			log.Tracef("devPath %s for docker volume %s", devPath, dockerVolume.Volume.Name)
			return dockerVolume.Volume.Name, nil
		}

		return "", fmt.Errorf("found a docker volume but its MountPoint was \"\"")
	}
	return dockerVolume.Volume.Name, nil
}

func findJSON(args []string, req *AttachRequest) (string, error) {
	var err error
	for i := 1; i < len(args); i++ {
		log.Tracef("findJSON(%d) about to unmarshal %v", i, args[i])
		err = json.Unmarshal([]byte(args[i]), req)
		if err == nil {
			return args[i], nil
		}
	}
	return "", err
}

func retryGetDockerPathAndMetadata(flexvolPath, devPath string) (string, string, error) {
	log.Tracef("retryGetDockerPathAndMetadata called with flexvolPath(%s) devPath(%s)", flexvolPath, devPath)
	maxTries := 3
	try := 0
	for {
		dockerPath, metadata, err := getDockerPathAndMetadata(flexvolPath, devPath)
		if err != nil {
			log.Errorf("getDockerPathAndMetadata failed for flexvolPath %s, devPath %s : %s", flexvolPath, devPath, err.Error())
			if try < maxTries {
				try++
				log.Tracef("try=%d", try)
				time.Sleep(time.Duration(try) * time.Second)
				continue
			}
			return dockerPath, metadata, err
		}
		if err != nil {
			return dockerPath, metadata, err
		}
		return dockerPath, metadata, nil
	}
}

// return the dockerPath and a path to the metadata file if present
func getDockerPathAndMetadata(flexvolPath, devPath string) (string, string, error) {
	dockerPath, err := linux.GetMountPointFromDevice(devPath)
	if err != nil {
		return "", "", err
	}
	log.Tracef("getDockerPathAndMetadata: devPath=%s was mounted on dockerPath=%s", devPath, dockerPath)

	metadata := ""
	if dockerPath == "" {
		// if we didn't get a docker path its because we're running
		// in a different namespace (likely rkt)
		log.Infof("getDockerPathAndMetadata: didn't find a docker path for devPath=%s and flexvolPath=%s", devPath, flexvolPath)

		metadata, err = getMountMetadataPath(flexvolPath)
		if err != nil {
			log.Errorf("getDockerPathAndMetadata: unable to read metadata=%s for devPath=%s and flexvolPath=%s", metadata, devPath, flexvolPath)
			return "", "", err
		}

		var fileData []byte
		fileData, err = ioutil.ReadFile(metadata)
		if err != nil {
			log.Errorf("getDockerPathAndMetadata: unable to read file content from metadata=%s for devPath=%s and flexvolPath=%s", metadata, dockerPath, flexvolPath)
			return "", "", err
		}
		dockerPath = string(fileData)
		log.Tracef("getDockerPathAndMetadata: found dockerPath=%s for devPath=%s and flexvolPath=%s", dockerPath, devPath, flexvolPath)
	}
	log.Tracef("getDockerPathAndMetadata: devPath=%s was mounted on dockerPath=%s", devPath, dockerPath)
	return dockerPath, metadata, nil
}

func getMountMetadataPath(flexvolPath string) (string, error) {
	_, flexvolFilename := filepath.Split(flexvolPath)
	if flexvolFilename == "" {
		return "", fmt.Errorf("unable to get filename from %s", flexvolPath)
	}
	execPathDir, _ := filepath.Split(execPath)
	if flexvolFilename == "" {
		return "", fmt.Errorf("unable to get dir from %s", execPath)
	}
	return fmt.Sprintf("%s.%s", execPathDir, flexvolFilename), nil
}
