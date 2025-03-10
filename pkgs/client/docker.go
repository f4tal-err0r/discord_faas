package client

import (
	"bufio"
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"

	docker "github.com/docker/docker/client"
)

type DockerImage struct {
	Name        string
	Runtime     string
	ContainerID string
	c           *docker.Client
}

func NewDockerClient() (*DockerImage, error) {
	opts := []docker.Opt{
		docker.FromEnv,
		docker.WithAPIVersionNegotiation(),
	}

	client, err := docker.NewClientWithOpts(opts...)
	if err != nil {
		return nil, err
	}

	return &DockerImage{c: client}, nil
}

func (d *DockerImage) BuildImage(dockerfile *os.File, tag string) (string, error) {
	resp, err := d.c.ImageBuild(context.Background(), dockerfile, types.ImageBuildOptions{
		Tags: []string{tag},
	})
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func (d *DockerImage) StopContainer() error {
	err := d.c.ContainerStop(context.Background(), d.ContainerID, container.StopOptions{})
	if err != nil {
		return err
	}
	err = d.c.ContainerRemove(context.Background(), d.ContainerID, container.RemoveOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (d *DockerImage) StartContainer() error {
	err := d.c.ContainerStart(context.Background(), d.ContainerID, container.StartOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (d *DockerImage) TailContainer() (*bufio.Scanner, error) {
	reader, err := d.c.ContainerLogs(context.Background(), d.ContainerID, container.LogsOptions{})
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	scanner := bufio.NewScanner(reader)

	return scanner, nil
}
