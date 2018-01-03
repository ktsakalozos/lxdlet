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

package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang/glog"
	"github.com/ktsakalozos/lxdlet/lxdlet/cmd/server/options"
	lxdImage "github.com/ktsakalozos/lxdlet/lxdlet/image"
	lxdRuntime "github.com/ktsakalozos/lxdlet/lxdlet/runtime"
	"github.com/spf13/pflag"
	"google.golang.org/grpc"
	"k8s.io/kubernetes/pkg/kubectl/util/logs"
	"k8s.io/kubernetes/pkg/kubelet/apis/cri/v1alpha1/runtime"
	"flag"
)

type RemoteRuntime struct {
	server *grpc.Server
}

func NewLxdRemoteRuntime() *RemoteRuntime {

	f := &RemoteRuntime{
		server: grpc.NewServer(),
	}
	runtime.RegisterRuntimeServiceServer(f.server, lxdRuntime.NewLxdRuntimeService())
	runtime.RegisterImageServiceServer(f.server, lxdImage.NewLxdImageService())

	return f
}

const defaultUnixSock = "/var/tmp/lxdlet.sock"

func main() {
	//This is needed to fix glog suffix parse error-> ERROR: logging before flag.Parse:
	flag.CommandLine.Parse([]string{})

	s := options.NewLxdletServer()
	s.AddFlags(pflag.CommandLine)
	logs.InitLogs()
	defer logs.FlushLogs()

	exitCh := make(chan os.Signal, 1)
	signal.Notify(exitCh, syscall.SIGINT, syscall.SIGTERM)

	glog.Info("This lxd CRI server implementation is under heavy development.")

	socketPath := defaultUnixSock
	defer os.Remove(socketPath)

	sock, err := net.Listen("unix", socketPath)
	if err != nil {
		glog.Fatalf("Error listening on sock %q: %v ", socketPath, err)
	}
	defer sock.Close()

	r := NewLxdRemoteRuntime()

	glog.Infof("Starting to serve on %q", socketPath)
	go r.server.Serve(sock)

	<-exitCh

	glog.Infof("lxdlet service exiting...")
}
