package docker

import (
	"context"
	"fmt"
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
	res, err := cli.Ping(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println(res.APIVersion)
	return &DockerClient{
		client: cli,
	}
}
