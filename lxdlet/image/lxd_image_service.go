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
	"fmt"
	"sync"

	"github.com/golang/glog"
	"github.com/ktsakalozos/lxdlet/lxdlet/util"
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
	lxdClient, err := util.NewLxdClient("/var/snap/lxd/common/lxd")
	if err != nil {
		return nil, err
	}

	image := req.GetImage().GetImage()
	_, err = lxdClient.DeleteImage(image, true)
	if err != nil {
		return nil, err
	}
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
	lxdClient, err := util.NewLxdClient("/var/snap/lxd/common/lxd")
	if err != nil {
		return nil, err
	}

	lxcimages, err := lxdClient.ListImages()
	if err != nil {
		return nil, err
	}

	var images []*runtime.Image
	for _, lxcImage := range lxcimages {
		if len(lxcImage.Aliases) == 0 {
			continue
		}
		image := &runtime.Image{
			Id: lxcImage.Aliases[0].Name,
		}

		images = append(images, image)
	}

	return &runtime.ListImagesResponse{Images: nil}, nil
}

// PullImage pulls an image into the store
func (s *LxdImageService) PullImage(ctx context.Context, req *runtime.PullImageRequest) (*runtime.PullImageResponse, error) {
	image := req.GetImage().GetImage()
	glog.Infof("+++++++ PullImage %s", image)
	/*
		// Do not pull any images. Create container should pull the image.
		lxdClient, err := util.NewLxdClient("/var/snap/lxd/common/lxd")
		if err != nil {
			return nil, err
		}

		_, err = lxdClient.PullImage(image, true)
		if err != nil {
			return nil, err
		}
	*/
	return &runtime.PullImageResponse{
		ImageRef: image,
	}, nil
}

// ImageFsInfo gets the info os an image
func (s *LxdImageService) ImageFsInfo(context.Context, *runtime.ImageFsInfoRequest) (*runtime.ImageFsInfoResponse, error) {
	glog.Infof("+++++++ ImageFsInfo ")
	return nil, fmt.Errorf("not implemented")
}
