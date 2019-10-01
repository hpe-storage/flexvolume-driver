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
	"os"
	"path"
	"strings"

	log "github.com/hpe-storage/common-host-libs/logger"
	"github.com/hpe-storage/common-host-libs/windows"
)

func doMount(flexvolPath, dockerPath, dockerName, mountID string) error {
	// Get volume object for a disk mounted by plugin.
	err := windows.AddPartitionAccessPath(flexvolPath, dockerPath)
	if err != nil {
		log.Tracef("Could not find mounted volume %s on the host, error %v", dockerPath, err.Error())
		return err
	}

	return nil
}

// Unbind the mount and delete mountpoint
func doUnMount(flexvolPath, mountID string) error {
	log.Tracef("doUnmount called flexvolpath %s, mount id %s", flexvolPath, mountID)
	//Get the docker volume name form the partition.
	dockerVolumeName, err := getDockerVolName(flexvolPath)
	log.Tracef("docker volume name %s", dockerVolumeName)
	if err != nil {
		log.Errorf("Could not find the docker volume name for flexvolPath %s  error %v", flexvolPath, err.Error())
		return err
	}
	// Remove Partition Access Path
	err = windows.RemovePartitionAccessPath(flexvolPath)
	if err != nil {
		log.Errorf("Could not remove the partition flexpath %s  error %v", flexvolPath, err.Error())
		//return err
	}

	log.Tracef("dvp unmount of %s %s", dockerVolumeName, mountID)
	err = dvp.Unmount(dockerVolumeName, mountID)
	if err != nil && !strings.Contains(err.Error(), notMounted) {
		return err
	}
	// cleanup pod volume folder otherwise k8s master keep calling unmount workflow
	os.Remove(flexvolPath)
	return nil
}

func getDockerVolName(flexvolPath string) (string, error) {
	log.Tracef(" getDockerVolName,  flexvolPath %s", flexvolPath)
	dockerVolPath, err := windows.GetDockerVolAccessPath(flexvolPath)
	if err != nil {
		return "", err
	}
	log.Tracef("dockerVolPath is %s", dockerVolPath)
	// C:\var\lib\kubelet\pods\47abfc90-a161-11e9-bbe1-005056962f3b\volumes\dory~nimble\sql-pv\
	// Change strings from windows forward slash to backward slash
	dockerVolPath = strings.Replace(dockerVolPath, "\\", "/", -1)
	//dockerVolPath = strings.Replace(dockerVolPath, "\r\n", "", -1)
	dockerVolPath = strings.TrimSuffix(dockerVolPath, "\r\n")
	log.Tracef("Format unit style windows dockerVolPath %s", dockerVolPath)
	lastIndex := len(dockerVolPath) - 1
	log.Tracef("lastIndex %d, length %d, %c", lastIndex, len(dockerVolPath), dockerVolPath[lastIndex])
	// Chop last forward slash if exists
	if dockerVolPath[lastIndex] == '/' {
		log.Tracef("Chop forward slash from dockerVolPath %s", dockerVolPath)
		dockerVolPath = dockerVolPath[:lastIndex]
	}
	dockerFolder, dockerVolName := path.Split(dockerVolPath)
	log.Tracef("After split dockerFolder %s, docker volume %s", dockerFolder, dockerVolName)
	return dockerVolName, nil

}
