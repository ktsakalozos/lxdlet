/*
Copyright 2016 The Kubernetes Authors.

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

package util

import (
	"github.com/lxc/lxd/client"
	"fmt"
	"github.com/lxc/lxd/shared/api"
)

type lxdDaemon struct {
	s    lxd.ContainerServer
	path string
	images       []api.Image
	networks     []api.Network
	storagePools []api.StoragePool
}

func NewLxdClient(path string) (*lxdDaemon, error) {
	// Connect to the LXD daemon
	s, err := lxd.ConnectLXDUnix(fmt.Sprintf("%s/unix.socket", path), nil)
	if err != nil {
		return nil, err
	}

	// Setup our internal struct
	d := &lxdDaemon{s: s, path: path}
	return d, nil
}

func (d *lxdDaemon) GetInfo() (*api.Server, error) {
	info, _, err := d.s.GetServer()
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (d *lxdDaemon) GetContainers() ([]api.Container, error) {
	// Containers
	containers, err := d.s.GetContainers()
	if err != nil {
		return nil, err
	}

	return containers, nil
}

func (d *lxdDaemon) GetContainer(name string) (*api.Container, error) {
	// Containers
	container, _, err := d.s.GetContainer(name)
	if err != nil {
		return nil, err
	}

	return container, nil
}

