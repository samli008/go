package main

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {
	cli, err := client.NewClientWithOpts(client.WithVersion("1.41"))
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\t%s\t%s\t%s\t%s\n", "容器名称", "容器状态", "容器镜像", "容器网络", "容器IP")
	for _, container := range containers {
		fmt.Printf("%s\t\t%s\t%s\t%s\t\t", container.Names[0][1:], container.Status, container.Image, container.HostConfig.NetworkMode)
		for _, v := range container.NetworkSettings.Networks {
			fmt.Printf("%s\n", v.IPAddress)
		}
	}
}
