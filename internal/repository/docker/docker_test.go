package docker

import "testing"

func TestImagePull(t *testing.T) {
	client := NewDockerClient()
	//client.ImagePull("nginx:latest")
	client.ImageList()
}
