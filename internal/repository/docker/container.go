package docker

import (
	"context"
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/pkg/archive"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"io"
	"time"
)

// 配置执行命令
func (docker *DockerClient) ContainerCreate(imageName string, containerName string, workingDir string, cmds []string, timeOut time.Duration) string {
	ctx := context.Background()
	stopTimeout := new(int)
	*stopTimeout = int(timeOut)
	config := &container.Config{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Image:        imageName,
		WorkingDir:   workingDir,
		Cmd:          cmds,
		StopTimeout:  stopTimeout,
	}
	hostConfig := &container.HostConfig{}
	networkingConfig := &network.NetworkingConfig{}
	platform := v1.Platform{}

	res, err := docker.client.ContainerCreate(ctx, config, hostConfig, networkingConfig, &platform, containerName)
	if err != nil {
		return ""
	}
	return res.ID
}

func (docker *DockerClient) CopyToContainer(containerId string, dstDir string, srcFile string) error {
	ctx := context.Background()
	content, err := archive.Tar(srcFile, archive.Uncompressed)
	if err != nil {
		return err
	}
	options := types.CopyToContainerOptions{}
	return docker.client.CopyToContainer(ctx, containerId, dstDir, content, options)
}

func (docker *DockerClient) ContainerWait(containerId string) (chanResponse <-chan container.WaitResponse, chanErr <-chan error) {
	ctx := context.Background()
	return docker.client.ContainerWait(ctx, containerId, container.WaitConditionNotRunning)
}

func (docker *DockerClient) ContainerStart(containerId string) error {
	ctx := context.Background()
	options := container.StartOptions{}
	return docker.client.ContainerStart(ctx, containerId, options)
}

// 获取执行时间，退出码等信息
func (docker *DockerClient) ContainerInspect(containerId string) (int, int64) {
	ctx := context.Background()
	res, err := docker.client.ContainerInspect(ctx, containerId)
	if err != nil {
		return 1, 0
	}

	exitCode := res.State.ExitCode

	start := res.State.StartedAt
	startTime, err := time.Parse(time.RFC3339, start)
	finish := res.State.FinishedAt
	endTime, err := time.Parse(time.RFC3339, finish)

	execTime := endTime.Sub(startTime).Milliseconds()
	return exitCode, execTime
}

// 获取执行结果
func (docker *DockerClient) ContainerLogs(containerId string) (string, error) {
	ctx := context.Background()
	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	}
	resp, err := docker.client.ContainerLogs(ctx, containerId, options)
	if err != nil {
		return "", err
	}
	defer resp.Close()
	if err != nil {
		panic(err)
	}

	// 读取容器日志
	logs, err := io.ReadAll(resp)
	if err != nil {
		return "", err
	}

	return string(logs), err
}

// 获取内存信息
func (docker *DockerClient) ContainerStats(containerId string) (uint64, error) {
	ctx := context.Background()
	res, err := docker.client.ContainerStats(ctx, containerId, true)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()
	var data types.StatsJSON
	err = json.NewDecoder(res.Body).Decode(&data)
	return data.MemoryStats.Usage, err
}

func (docker *DockerClient) ContainerStop(containerId string) error {
	ctx := context.Background()
	options := container.StopOptions{}
	return docker.client.ContainerStop(ctx, containerId, options)
}

func (docker *DockerClient) ContainerRemove(containerId string) error {
	ctx := context.Background()
	options := container.RemoveOptions{}
	return docker.client.ContainerRemove(ctx, containerId, options)
}

func (docker *DockerClient) IsContainerRunning(containerId string) (bool, error) {
	ctx := context.Background()
	res, err := docker.client.ContainerInspect(ctx, containerId)
	if err != nil {
		return false, err
	}
	return res.State.Running, nil
}
