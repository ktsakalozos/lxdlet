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

package runtime

import (
	"bytes"
	"io/ioutil"
	"math/rand"
	"os"
	"syscall"
	"time"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"

	runtimeApi "k8s.io/kubernetes/pkg/kubelet/apis/cri/v1alpha1/runtime"
	"k8s.io/kubernetes/pkg/kubelet/server/streaming"

	"github.com/ktsakalozos/lxdlet/lxdlet/util"
)

type LxdRuntime struct {
	imageStore  runtimeApi.ImageServiceServer
	lxdDataPath string
}

const internalAppPrefix = "lxdletinternal-"
const sandboxes_path = "/var/tmp/lxdlet/sandboxes"

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// New creates a new RuntimeServiceServer backed by lxd
func NewLxdRuntimeService() runtimeApi.RuntimeServiceServer {
	rand.Seed(time.Now().UnixNano())
	_ = os.MkdirAll(sandboxes_path, os.ModePerm)
	runtime := &LxdRuntime{
		imageStore:  nil,
		lxdDataPath: sandboxes_path,
	}

	streamConfig := streaming.DefaultConfig
	go func() {
		// TODO, runtime.streamServer.Stop() for SIGTERM or any other clean
		// shutdown of rktlet
		glog.Infof("listening for execs on: %v", streamConfig.Addr)
	}()

	return runtime
}

func translateState(lxdState string) runtimeApi.ContainerState {
	if lxdState == "RUNNING" {
		return runtimeApi.ContainerState_CONTAINER_RUNNING
	}

	return runtimeApi.ContainerState_CONTAINER_UNKNOWN
}

func (r *LxdRuntime) Version(ctx context.Context, req *runtimeApi.VersionRequest) (*runtimeApi.VersionResponse, error) {
	glog.Infof("*********** Version ")
	name := "lxd"
	version := "0.1.0"
	lxdClient, err := util.NewLxdClient("/var/snap/lxd/common/lxd")
	if err != nil {
		return nil, err
	}
	info, err := lxdClient.GetInfo()
	if err != nil {
		return nil, err
	}

	return &runtimeApi.VersionResponse{
		Version:           version, // kubelet/remote version, must be 0.1.0
		RuntimeName:       name,
		RuntimeVersion:    info.Environment.ServerVersion,
		RuntimeApiVersion: info.APIVersion,
	}, nil
}

func (r *LxdRuntime) ListContainers(ctx context.Context, req *runtimeApi.ListContainersRequest) (*runtimeApi.ListContainersResponse, error) {
	// We assume the containers in data dir are all managed by kubelet.
	glog.V(6).Infof("*********** ListContainers ")

	lxdClient, err := util.NewLxdClient("/var/snap/lxd/common/lxd")
	if err != nil {
		return nil, err
	}

	// TODO(kjackal): be smarter here and ask only the container requested
	allLxcContainers, err := lxdClient.GetContainers()
	if err != nil {
		return nil, err
	}

	var containers []*runtimeApi.Container
	for _, lxcContainer := range allLxcContainers {

		imgSpec := &runtimeApi.ImageSpec{
			Image: lxcContainer.Config["image.serial"],
		}
		container := &runtimeApi.Container{
			//			Annotations:  resp.Status.Annotations,
			CreatedAt: lxcContainer.CreatedAt.UnixNano(),    // resp.Status.CreatedAt,
			Id:        lxcContainer.Name,                    //resp.Status.Id,
			Image:     imgSpec,                              //resp.Status.Image,
			ImageRef:  lxcContainer.Config["image.release"], //"resp.Status.ImageRef",
			//			Labels:       resp.Status.Labels,
			//			Metadata:     resp.Status.Metadata,
			//			PodSandboxId: p.UUID,
			State: translateState(lxcContainer.Config["volatile.last_state.power"]),
		}

		//		if passFilter(container, req.Filter) {
		containers = append(containers, container)
		//		}
	}

	return &runtimeApi.ListContainersResponse{Containers: containers}, nil
}

func (r *LxdRuntime) ContainerStatus(ctx context.Context, req *runtimeApi.ContainerStatusRequest) (*runtimeApi.ContainerStatusResponse, error) {
	// Container ID is in the form of "uuid:appName".
	glog.V(6).Infof("*********** ContainerStatus : ", req.ContainerId)
	/*	uuid, appName, err := parseContainerID(req.ContainerId)
		if err != nil {
			return nil, err
		}

		// TODO(kjackal): decipher this
		resp, err := r.RunCommand("app", "status", uuid, "--app="+appName, "--format=json")
		if err != nil {
			return nil, err
		}

		if len(resp) != 1 {
			return nil, fmt.Errorf("unexpected result %q", resp)
		}
	*/
	/*
		var app rkt.App
		if err := json.Unmarshal([]byte(resp[0]), &app); err != nil {
			return nil, fmt.Errorf("failed to unmarshal container: %v", err)
		}

		status, err := toContainerStatus(uuid, &app)
		if err != nil {
			return nil, fmt.Errorf("failed to convert to container status: %v", err)
		}
	*/
	//return &runtimeApi.ContainerStatusResponse{Status: status}, nil
	return &runtimeApi.ContainerStatusResponse{Status: nil}, nil
}

func (r *LxdRuntime) CreateContainer(ctx context.Context, req *runtimeApi.CreateContainerRequest) (*runtimeApi.CreateContainerResponse, error) {
	imageID := req.GetConfig().GetImage().Image

	glog.V(6).Infof("*********** CreateContainer called with image: ", imageID)

	// TODO(kjackal): decipher this
	/*	command, err := generateAppAddCommand(req, imageID)
		if err != nil {
			return nil, err
		}
		if _, err := r.RunCommand(command[0], command[1:]...); err != nil {
			return nil, err
		}

		appName := buildAppName(req.Config.Metadata.Attempt, req.Config.Metadata.Name)
		containerID := buildContainerID(req.PodSandboxId, appName)

		return &runtimeApi.CreateContainerResponse{ContainerId: containerID}, nil
	*/
	return &runtimeApi.CreateContainerResponse{ContainerId: ""}, nil
}

func (r *LxdRuntime) StartContainer(ctx context.Context, req *runtimeApi.StartContainerRequest) (*runtimeApi.StartContainerResponse, error) {
	// Container ID is in the form of "uuid:appName".
	glog.V(6).Infof("*********** StartContainer contained id: ", req.ContainerId)
	/*	uuid, appName, err := parseContainerID(req.ContainerId)
		if err != nil {
			return nil, err
		}

		// TODO(kjackal): decipher this
		if _, err := r.RunCommand("app", "start", uuid, "--app="+appName); err != nil {
			return nil, err
		}
	*/
	return &runtimeApi.StartContainerResponse{}, nil
}

func (r *LxdRuntime) StopContainer(ctx context.Context, req *runtimeApi.StopContainerRequest) (*runtimeApi.StopContainerResponse, error) {
	// Container ID is in the form of "uuid:appName".
	glog.V(6).Infof("*********** StopContainer contained id: ", req.ContainerId)
	/*	uuid, appName, err := parseContainerID(req.ContainerId)
		if err != nil {
			return nil, err
		}

		// TODO(kjackal): decipher this
		// TODO(yifan): Support timeout.
		if _, err := r.RunCommand("app", "stop", uuid, "--app="+appName); err != nil {
			return nil, err
		}
	*/
	return &runtimeApi.StopContainerResponse{}, nil
}

func (r *LxdRuntime) RemoveContainer(ctx context.Context, req *runtimeApi.RemoveContainerRequest) (*runtimeApi.RemoveContainerResponse, error) {
	// Container ID is in the form of "uuid:appName".
	glog.V(6).Infof("*********** RemoveContainer contained id: ", req.ContainerId)
	/*
		uuid, appName, err := parseContainerID(req.ContainerId)
		if err != nil {
			return nil, err
		}
		// TODO(kjackal): decipher this

		// TODO(yifan): Support timeout.
		if _, err := r.RunCommand("app", "rm", uuid, "--app="+appName); err != nil {
			return nil, err
		}
	*/
	return &runtimeApi.RemoveContainerResponse{}, nil
}

func (r *LxdRuntime) UpdateRuntimeConfig(ctx context.Context, req *runtimeApi.UpdateRuntimeConfigRequest) (*runtimeApi.UpdateRuntimeConfigResponse, error) {
	// TODO, use the PodCIDR passed in once we have network plugins setup
	return &runtimeApi.UpdateRuntimeConfigResponse{}, nil
}

func (r *LxdRuntime) Status(ctx context.Context, req *runtimeApi.StatusRequest) (*runtimeApi.StatusResponse, error) {
	// TODO: implement

	glog.V(6).Infof("*********** Status")
	//Need to copy the consts to get pointers
	runtimeReady := runtimeApi.RuntimeReady
	networkReady := runtimeApi.NetworkReady
	tv := true

	conditions := []*runtimeApi.RuntimeCondition{
		&runtimeApi.RuntimeCondition{
			Type:   runtimeReady,
			Status: tv,
		},
		&runtimeApi.RuntimeCondition{
			Type:   networkReady,
			Status: tv,
		},
	}
	resp := runtimeApi.StatusResponse{
		Status: &runtimeApi.RuntimeStatus{
			Conditions: conditions,
		},
	}

	return &resp, nil
}

func (r *LxdRuntime) Attach(ctx context.Context, req *runtimeApi.AttachRequest) (*runtimeApi.AttachResponse, error) {
	return nil, nil
}

func (r *LxdRuntime) Exec(ctx context.Context, req *runtimeApi.ExecRequest) (*runtimeApi.ExecResponse, error) {
	return nil, nil
}

func (r *LxdRuntime) ExecSync(ctx context.Context, req *runtimeApi.ExecSyncRequest) (*runtimeApi.ExecSyncResponse, error) {
	return nil, nil
}

func (r *LxdRuntime) PortForward(ctx context.Context, req *runtimeApi.PortForwardRequest) (*runtimeApi.PortForwardResponse, error) {
	return nil, nil
}

// ContainerStats returns stats of the container. If the container does not
// exist, the call returns an error.
func (r *LxdRuntime) ContainerStats(ctx context.Context, req *runtimeApi.ContainerStatsRequest) (*runtimeApi.ContainerStatsResponse, error) {
	resp := runtimeApi.ContainerStatsResponse{}
	return &resp, nil
}

// ListContainerStats returns stats of all running containers.
func (r *LxdRuntime) ListContainerStats(context.Context, *runtimeApi.ListContainerStatsRequest) (*runtimeApi.ListContainerStatsResponse, error) {
	return nil, nil
}

///////////////////////
// Sandbox functions //
///////////////////////
func (r *LxdRuntime) getPodPath(podUUID string) string {
	var strBuffer bytes.Buffer
	strBuffer.WriteString(r.lxdDataPath)
	strBuffer.WriteString("/")
	strBuffer.WriteString(podUUID)
	return strBuffer.String()
}

func (r *LxdRuntime) getPodStatus(podID string) *runtimeApi.PodSandboxStatus {
	path := r.getPodPath(podID)
	status := runtimeApi.PodSandboxState_SANDBOX_NOTREADY
	var createdAt int64
	if stats, err := os.Stat(path); err == nil {
		status = runtimeApi.PodSandboxState_SANDBOX_READY
		var unixStat = stats.Sys().(*syscall.Stat_t)
		createdAt = unixStat.Ctim.Sec
	}

	return &runtimeApi.PodSandboxStatus{
		Id:        podID,
		State:     status,
		CreatedAt: createdAt,
	}
}

func podSandboxStatusMatchesFilter(sbx *runtimeApi.PodSandboxStatus, filter *runtimeApi.PodSandboxFilter) bool {
	if filter == nil {
		return true
	}
	if filter.Id != "" && filter.Id != sbx.Id {
		return false
	}

	if filter.State != nil && filter.GetState().State != sbx.State {
		return false
	}

	for key, val := range filter.LabelSelector {
		sbxLabel, exists := sbx.Labels[key]
		if !exists {
			return false
		}
		if sbxLabel != val {
			return false
		}
	}

	return true
}

// RunPodSandbox creates and starts a Pod
func (r *LxdRuntime) RunPodSandbox(ctx context.Context, req *runtimeApi.RunPodSandboxRequest) (*runtimeApi.RunPodSandboxResponse, error) {
	glog.Infof("======= RunPodSandbox ")
	podUUID := randString(64)
	serialisedRequest, err := proto.Marshal(req)
	if err != nil {
		glog.Error("Failed to masharl snadbox creation request.")
		return nil, err
	}

	err = ioutil.WriteFile(r.getPodPath(podUUID), serialisedRequest, 0644)
	if err != nil {
		glog.Error("Failed to store snadbox creation request.")
		return nil, err
	}

	return &runtimeApi.RunPodSandboxResponse{PodSandboxId: podUUID}, nil
}

// StopPodSandbox stops a pod
func (r *LxdRuntime) StopPodSandbox(ctx context.Context, req *runtimeApi.StopPodSandboxRequest) (*runtimeApi.StopPodSandboxResponse, error) {
	glog.V(4).Infof("======= StopPodSandbox %s", req.PodSandboxId)
	//err := r.stopPodSandbox(ctx, req.PodSandboxId, false)
	// TODO(kjackal): Stop the container if running on this sandbox
	return &runtimeApi.StopPodSandboxResponse{}, nil
}

// RemovePodSandbox removes a pod
func (r *LxdRuntime) RemovePodSandbox(ctx context.Context, req *runtimeApi.RemovePodSandboxRequest) (*runtimeApi.RemovePodSandboxResponse, error) {
	glog.V(4).Infof("======= RemovePodSandbox %s", req.PodSandboxId)
	// Force stop first, per api contract "if there are any running containers in
	// the sandbox, they must be forcibly terminated
	//r.stopPodSandbox(ctx, req.PodSandboxId, true)
	// TODO(kjackal): Stop the container if running on this sandbox
	os.Remove(r.getPodPath(req.PodSandboxId))
	return &runtimeApi.RemovePodSandboxResponse{}, nil
}

// PodSandboxStatus gets the status of a pod
func (r *LxdRuntime) PodSandboxStatus(ctx context.Context, req *runtimeApi.PodSandboxStatusRequest) (*runtimeApi.PodSandboxStatusResponse, error) {
	glog.V(4).Infof("======= PodSandboxStatus %s", req.PodSandboxId)
	podStatus := r.getPodStatus(req.PodSandboxId)
	return &runtimeApi.PodSandboxStatusResponse{Status: podStatus}, nil
}

// ListPodSandbox lists all pods
func (r *LxdRuntime) ListPodSandbox(ctx context.Context, req *runtimeApi.ListPodSandboxRequest) (*runtimeApi.ListPodSandboxResponse, error) {
	glog.V(4).Infof("======= ListPodSandbox")
	files, err := ioutil.ReadDir(r.lxdDataPath)
	if err != nil {
		glog.Error("Failed to list pods.")
		return nil, err
	}

	sandboxes := make([]*runtimeApi.PodSandbox, 0, len(files))
	for _, file := range files {
		podID := file.Name()
		sandboxStatus := r.getPodStatus(podID)

		if !podSandboxStatusMatchesFilter(sandboxStatus, req.GetFilter()) {
			continue
		}
		sandboxes = append(sandboxes, &runtimeApi.PodSandbox{
			Id:        sandboxStatus.Id,
			Labels:    sandboxStatus.Labels,
			Metadata:  sandboxStatus.Metadata,
			State:     sandboxStatus.State,
			CreatedAt: sandboxStatus.CreatedAt,
		})
	}
	return &runtimeApi.ListPodSandboxResponse{Items: sandboxes}, nil
}

// UpdateContainerResources updates ContainerConfig of the container.
func (r *LxdRuntime) UpdateContainerResources(context.Context, *runtimeApi.UpdateContainerResourcesRequest) (*runtimeApi.UpdateContainerResourcesResponse, error) {
	return nil, nil
}