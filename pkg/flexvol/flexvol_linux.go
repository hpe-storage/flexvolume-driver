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
	"fmt"
	"github.com/hpe-storage/common-host-libs/docker/dockervol"
	"github.com/hpe-storage/common-host-libs/linux"
	log "github.com/hpe-storage/common-host-libs/logger"
	"io/ioutil"
	"os"
	"strings"
)

// TO-DO - We should use CHAPI module instead of using PowerShell command
// directly.
func doMount(flexvolPath, dockerPath, dockerName, mountID string) error {
	devPath, err := linux.GetDeviceFromMountPoint(dockerPath)
	if err != nil {
		return err
	}

	if devPath == "" {
		//we're probably running in a different namespace
		//so we need to pull the device path from the
		//docker volume driver
		log.Infof("doMount: devPath was empty for flexvolPath=% volume=%s", flexvolPath, dockerName)

		//get the volume info
		var volRes *dockervol.GetResponse
		volRes, err = getVolume(dockerName)
		if err != nil {
			return err
		}

		devPath, found := volRes.Volume.Status[devicePathKey].(string)
		if !found || devPath == "" {
			log.Errorf("Unable to get device for flexvolPath=%s from docker volume=%+v (path=%s)", flexvolPath, volRes, dockerPath)
			return fmt.Errorf("Unable to get device for flexvolPath=%s from docker volume=%s", flexvolPath, dockerPath)
		}
		log.Tracef("doMount: found devPath=%s for volume=%s", devPath, dockerName)

		//mount devicePath onto flexvolPath
		_, err = linux.MountDeviceWithFileSystem(devPath, flexvolPath, nil)
		if err != nil {
			return err
		}
		log.Tracef("doMount: mounted devPath=%s at flexvolPath=%s", devPath, flexvolPath)

		//create a hidden file in the flexvolume path that maps flexvolume mount to the docker volume (breadcrumb)
		var metadata string
		metadata, err = getMountMetadataPath(flexvolPath)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(metadata, []byte(dockerPath), 0600)
		if err != nil {
			return err
		}
		log.Tracef("doMount: stored dockerPath=%s at metadata=%s", dockerPath, metadata)

	} else {
		//bind mount the docker path to the flexvol path
		_ = linux.RetryBindMount(dockerPath, flexvolPath, false)
		log.Tracef("doMount: bind mounted dockerPath=%s at flexvolPath=%s", dockerPath, flexvolPath)
	}
	// Set selinux context if configured
	// References:
	//    https://github.com/kubernetes/kubernetes/issues/20813
	//    https://github.com/openshift/origin/issues/741
	//    https://github.com/projectatomic/atomic-site/blob/master/source/blog/2015-06-15-using-volumes-with-docker-can-cause-problems-with-selinux.html.md
	err = linux.Chcon("svirt_sandbox_file_t", flexvolPath)
	if err != nil {
		return err
	}

	return err
}

// Unbind the mount and delete mountpoint
//nolint :gocyclo
func doUnMount(flexvolPath, mountID string) error {

	devPath, err := linux.GetDeviceFromMountPoint(flexvolPath)
	if err != nil {
		return err
	}

	log.Tracef("Umount of \"%s\" from %s", flexvolPath, devPath)
	err = linux.BindUnmount(flexvolPath)
	if err != nil && !strings.Contains(err.Error(), notMounted) {
		return err
	}

	dockerPath, metadata, err := retryGetDockerPathAndMetadata(flexvolPath, devPath)
	if err != nil && !strings.Contains(strings.ToLower(err.Error()), noFileOrDirErr) {
		return err
	}

	dockerVolumeName, err := retryGetVolumeNameFromMountPath(flexvolPath, dockerPath)
	if err != nil {
		return err
	}
	log.Tracef("docker unmount of %s %s", dockerVolumeName, mountID)
	err = dvp.Unmount(dockerVolumeName, mountID)
	if err != nil && !strings.Contains(err.Error(), notMounted) {
		return err
	}

	if metadata != "" {
		dockerVolumeName, err = getVolumeNameFromMountPath(flexvolPath, dockerPath)
		if err != nil {
			// an error means that we didn't find the volume mounted
			// this means we can clean up the breadcrumbs
			log.Tracef("Unmount: removing metadata=%s", metadata)
			os.Remove(metadata)
		}
		log.Tracef("Unmount: dockerVolumeName=%s still has an active mount at %s.", dockerVolumeName, dockerPath)
	}
	return nil
}
