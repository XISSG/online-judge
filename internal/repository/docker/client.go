package docker

import (
	"github.com/docker/docker/client"
)

type DockerClient struct {
	client *client.Client
}

func NewDockerClient() *DockerClient {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	return &DockerClient{
		client: cli,
	}
}
