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

package image

import (
	"sync"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"k8s.io/kubernetes/pkg/kubelet/apis/cri/v1alpha1/runtime"
)

// LxdImageService holds the implementation of the images repo
type LxdImageService struct {
	sync.Mutex

	FakeImageSize uint64
	Called        []string
}

// NewLxdImageService creates an image storage that allows CRUD operations for images.
func NewLxdImageService() runtime.ImageServiceServer {
	return &LxdImageService{}
}

// RemoveImage removes the image from the image store.
func (s *LxdImageService) RemoveImage(ctx context.Context, req *runtime.RemoveImageRequest) (*runtime.RemoveImageResponse, error) {
	glog.Infof("+++++++ RemoveImage ")
	return &runtime.RemoveImageResponse{}, nil
}

// ImageStatus returns the status of the image.
func (s *LxdImageService) ImageStatus(ctx context.Context, req *runtime.ImageStatusRequest) (*runtime.ImageStatusResponse, error) {
	glog.Infof("+++++++ ImageStatus ")
	// api expected response for "Image does not exist"
	return &runtime.ImageStatusResponse{}, nil
}

// ListImages lists images in the store
func (s *LxdImageService) ListImages(ctx context.Context, req *runtime.ListImagesRequest) (*runtime.ListImagesResponse, error) {
	glog.Infof("+++++++ ListImages ")
	return &runtime.ListImagesResponse{Images: nil}, nil
}

// PullImage pulls an image into the store
func (s *LxdImageService) PullImage(ctx context.Context, req *runtime.PullImageRequest) (*runtime.PullImageResponse, error) {
	glog.Infof("+++++++ PullImage ")
	return &runtime.PullImageResponse{
		ImageRef: "",
	}, nil
}

// ImageFsInfo gets the info os an image
func (s *LxdImageService) ImageFsInfo(context.Context, *runtime.ImageFsInfoRequest) (*runtime.ImageFsInfoResponse, error) {
	glog.Infof("+++++++ ImageFsInfo ")
	return &runtime.ImageFsInfoResponse{ImageFilesystems: nil}, nil
}
