package docker

import (
	"bytes"
	"context"
	"fmt"
	"github.com/docker/docker/api/types/image"
	"github.com/xissg/online-judge/internal/constant"
	"io"
)

func (docker *DockerClient) ImageList() ([]string, error) {
	ctx := context.Background()
	options := image.ListOptions{
		All: true,
	}
	res, err := docker.client.ImageList(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("repository layer: docker, image list: %w %+v", constant.ErrInternal, err)
	}

	var result []string
	for _, img := range res {
		result = append(result, img.RepoTags...)
	}
	return result, nil
}
func (docker *DockerClient) ImagePull(imageTag string) error {
	resp, err := docker.client.ImagePull(context.Background(), imageTag, image.PullOptions{})
	if err != nil {
		return fmt.Errorf("repository layer: docker, image pull: %w %+v", constant.ErrInternal, err)
	}
	defer resp.Close()
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, resp)
	if err != nil {
		return fmt.Errorf("repository layer: docker, image pull; io copy: %w %+v", constant.ErrInternal, err)
	}

	return nil
}

func (docker *DockerClient) ImageRemove(imageId string) error {
	_, err := docker.client.ImageRemove(context.Background(), imageId, image.RemoveOptions{})
	if err != nil {
		return fmt.Errorf("repository layer: docker, image remove: %w %+v", constant.ErrInternal, err)
	}
	return nil
}
