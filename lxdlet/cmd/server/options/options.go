/*
Copyright 2015 The Kubernetes Authors.

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

// Package options contains all of primary arguments for a rktlet
package options

import (
	//"github.com/ktsakalozos/lxdlet/lxdlet"
	"github.com/spf13/pflag"
)

type LxdletServer struct {
	//*lxdlet.Config
}

func NewLxdletServer() *LxdletServer {
	//config := lxdlet.DefaultConfig
	return &LxdletServer{
	//	Config: config,
	}
}

func (s *LxdletServer) AddFlags(fs *pflag.FlagSet) {
/*
	fs.StringVar(&s.LxdPath, "lxd-path", s.LxdPath, "Path of lxd binary. Leave empty to use the first lxd in $PATH.")
	fs.StringVar(&s.LxdDatadir, "lxd-data-dir", s.LxdDatadir, "Path to lxd's data directory. Defaults to '/var/lib/lxdlet/data'.")
	fs.StringVar(&s.StreamServerAddress, "stream-server-address", s.StreamServerAddress, "Address to listen on for api-server streaming requests. MUST BE SECURED BY SOME EXTERNAL MECHANISM.")
	fs.StringVar(&s.LxdStage1Name, "lxd-stage1-name", s.LxdStage1Name, "Name of an image to use as stage1. This needs to be specified as 'image:version'. If the image is present in the local store, the version can be ommitted.")
*/
}
