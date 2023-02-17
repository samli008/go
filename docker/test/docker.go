package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

var name = "node4"
var hostname = "node4"
var net = "net20"
var ip = "192.168.20.18"
var image = "c7:nodejs"

func main() {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	// result, err1 := create(cli, name, net, ip, image, hostname)
	// list(cli)
	// if err1 != nil {
	// 	fmt.Println(err1)
	// } else {
	// 	fmt.Println(result)
	// }

	// id, err1 := removeContainer("9483ac1a1156", cli)
	// if err1 != nil {
	// 	fmt.Println(err1)
	// } else {
	// 	fmt.Println(id + " deleted")
	// }

	//list(cli)
	// listImage(cli)
	//netList(cli)
	netRange(cli, "net20")
}

func listImage(cli *client.Client) {
	images, _ := cli.ImageList(context.Background(), types.ImageListOptions{})

	for _, image := range images {
		fmt.Printf("%s\t%s\t%s\n", image.RepoTags[0], image.ID[7:19], strconv.Itoa(int(image.Size))[:3]+"MB")
	}
}

func create(cli *client.Client, containName string, net string, ip string, image string, hostname string) (string, error) {
	config := &container.Config{
		Image:    image,
		Hostname: hostname,
	}

	hostConfig := &container.HostConfig{
		Privileged: true,
		RestartPolicy: container.RestartPolicy{
			Name: "always",
		},
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: "/opt",
				Target: "/opt",
			},
		},
	}

	netConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			"network": {
				NetworkID: net,
				IPAMConfig: &network.EndpointIPAMConfig{
					IPv4Address: ip,
				},
			},
		},
	}
	resp, err1 := cli.ContainerCreate(context.Background(), config, hostConfig, netConfig, nil, containName)
	if err1 != nil {
		return resp.ID, err1
	}
	cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{})
	return "created ok", err1
}

func removeContainer(containerID string, cli *client.Client) (string, error) {
	err := cli.ContainerRemove(context.Background(), containerID, types.ContainerRemoveOptions{Force: true})
	return containerID, err
}

func list(cli *client.Client) {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\t%s\t%s\t%s\t%s\t%s\n", "容器名称", "容器状态", "容器镜像", "容器ID", "容器网络", "容器IP")
	for _, container := range containers {
		fmt.Printf("%s\t\t%s\t%s\t%s\t\t", container.Names[0][1:], container.Status, container.Image, container.ID[:12])
		for i, v := range container.NetworkSettings.Networks {
			if i != "network" {
				fmt.Printf("%s\n", i)
				fmt.Printf("%s\n", v.IPAddress)
			}
		}
	}
}

func netList(cli *client.Client) {
	nets, err := cli.NetworkList(context.Background(), types.NetworkListOptions{})
	if err != nil {
		panic(err)
	}
	lists := make([]string, 0)
	fmt.Println()
	fmt.Printf("%s\t%s\n", "网络名称", "网络ID")
	for _, net := range nets {
		lists = append(lists, net.Name)
		if net.Driver == "macvlan" || net.Driver == "bridge" {
			fmt.Printf("%s\t\t%s\n", net.Name, net.ID[:12])
		}
	}
	fmt.Println()
}

func netRange(cli *client.Client, netId string) {
	nets, err := cli.NetworkInspect(context.Background(), netId, types.NetworkInspectOptions{})
	if err != nil {
		panic(err)
	}
	//lists := make([]string, 0)
	fmt.Println()
	fmt.Printf("%s\n", nets.Driver)
	fmt.Printf("%s\n", nets.Options["parent"])
	fmt.Printf("%s\t%s\n", nets.IPAM.Config[0].Subnet, nets.IPAM.Config[0].Gateway)
	for _, v := range nets.Containers {
		fmt.Printf("%s\t%s\n", v.Name, v.IPv4Address)
	}
	fmt.Println()
}
