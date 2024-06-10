package docker

import (
	"bytes"
	"context"
	"fmt"
	"github.com/docker/docker/api/types/image"
	"io"
	"log"
)

func (docker *DockerClient) ImageList() []string {
	ctx := context.Background()
	options := image.ListOptions{
		All: true,
	}
	res, err := docker.client.ImageList(ctx, options)
	if err != nil {
		log.Printf("Error getting")
	}

	var result []string
	for _, img := range res {
		result = append(result, img.RepoTags...)
	}
	return result
}
func (docker *DockerClient) ImagePull(imageTag string) {
	resp, err := docker.client.ImagePull(context.Background(), imageTag, image.PullOptions{})
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Close()
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, resp)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(buf.String())

}

func (docker *DockerClient) ImageRemove(imageId string) {
	_, err := docker.client.ImageRemove(context.Background(), imageId, image.RemoveOptions{})
	if err != nil {
		log.Fatal(err)
	}
}
