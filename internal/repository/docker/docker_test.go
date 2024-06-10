package docker

import (
	"fmt"
	"testing"
)

func TestImagePull(t *testing.T) {
	client := NewDockerClient()
	//client.ImagePull("nginx:latest")
	images := client.ImageList()
	fmt.Println(images)
}
