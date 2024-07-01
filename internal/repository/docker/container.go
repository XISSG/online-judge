package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/pkg/archive"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/xissg/online-judge/internal/constant"
	"io"
	"time"
)

// 配置执行命令
func (docker *DockerClient) ContainerCreate(imageName string, containerName string, workingDir string, cmds []string, timeOut time.Duration) (string, error) {
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
		return "", fmt.Errorf("repository layer: docker, container create: %w %+v", constant.ErrInternal, err)
	}
	return res.ID, nil
}

func (docker *DockerClient) CopyToContainer(containerId string, dstDir string, srcFile string) error {
	ctx := context.Background()
	content, err := archive.Tar(srcFile, archive.Uncompressed)
	if err != nil {
		return fmt.Errorf("repository layer: docker, copy to container; archive.Tar: %w %+v", constant.ErrInternal, err)
	}
	options := types.CopyToContainerOptions{}
	err = docker.client.CopyToContainer(ctx, containerId, dstDir, content, options)
	if err != nil {
		return fmt.Errorf("repository layer: docker, copy to container: %w %+v", constant.ErrInternal, err)
	}
	return nil
}

func (docker *DockerClient) ContainerWait(containerId string) (chanResponse <-chan container.WaitResponse, chanErr <-chan error) {
	ctx := context.Background()
	return docker.client.ContainerWait(ctx, containerId, container.WaitConditionNotRunning)
}

func (docker *DockerClient) ContainerStart(containerId string) error {
	ctx := context.Background()
	options := container.StartOptions{}
	err := docker.client.ContainerStart(ctx, containerId, options)
	if err != nil {
		return fmt.Errorf("repository layer: docker, container start: %w %+v", constant.ErrInternal, err)
	}
	return nil
}

// 获取执行时间，退出码等信息
func (docker *DockerClient) ContainerInspect(containerId string) (exitCode int, execTime int64) {
	ctx := context.Background()
	res, err := docker.client.ContainerInspect(ctx, containerId)
	if err != nil {
		return 1, 0
	}

	exitCode = res.State.ExitCode

	start := res.State.StartedAt
	startTime, err := time.Parse(time.RFC3339, start)
	finish := res.State.FinishedAt
	endTime, err := time.Parse(time.RFC3339, finish)

	execTime = endTime.Sub(startTime).Milliseconds()
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
		return "", fmt.Errorf("repository layer: docker, container logs: %w %+v", constant.ErrInternal, err)
	}

	// 读取容器日志
	logs, err := io.ReadAll(resp)
	if err != nil {
		return "", fmt.Errorf("repository layer: docker, container logs;io.ReadAll: %w %+v", constant.ErrInternal, err)
	}

	return string(logs), err
}

// 获取内存信息
func (docker *DockerClient) ContainerStats(containerId string) (uint64, error) {
	ctx := context.Background()
	res, err := docker.client.ContainerStats(ctx, containerId, true)
	if err != nil {
		return 0, fmt.Errorf("repository layer: docker, container stats: %w %+v", constant.ErrInternal, err)
	}
	defer res.Body.Close()
	var data types.StatsJSON
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return 0, fmt.Errorf("repository layer: docker, container stats;json.NewDecoder.Decode: %w %+v", constant.ErrInternal, err)
	}
	return data.MemoryStats.Usage, nil
}

func (docker *DockerClient) ContainerStop(containerId string) error {
	ctx := context.Background()
	options := container.StopOptions{}
	err := docker.client.ContainerStop(ctx, containerId, options)
	if err != nil {
		return fmt.Errorf("repository layer: docker, container stop: %w %+v", constant.ErrInternal, err)
	}
	return nil
}

func (docker *DockerClient) ContainerRemove(containerId string) error {
	ctx := context.Background()
	options := container.RemoveOptions{}
	err := docker.client.ContainerRemove(ctx, containerId, options)
	if err != nil {
		return fmt.Errorf("repository layer: docker, container remove: %w %+v", constant.ErrInternal, err)
	}
	return nil
}

func (docker *DockerClient) IsContainerRunning(containerId string) (bool, error) {
	ctx := context.Background()
	res, err := docker.client.ContainerInspect(ctx, containerId)
	if err != nil {
		return false, fmt.Errorf("repository layer: docker, is docker running;container inspect: %w %+v", constant.ErrInternal, err)
	}
	return res.State.Running, nil
}
